package google

import (
	"context"
	"github.com/flyteorg/flytestdlib/atomic"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/golang/groupcache/singleflight"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	v1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"sync"
	"time"
)

const (
	defaultGracePeriod             = 300
	gcpServiceAccountAnnotationKey = "iam.gke.io/gcp-service-account"
	workflowIdentityDocUrl         = "https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity"
)

type tokensKey struct {
	k8sNamespace      string
	k8sServiceAccount string
}

type GkeTokenSource struct {
	kubeClient        kubernetes.Interface
	tokens            sync.Map
	singleflight      singleflight.Group
	identityNamespace string
	scope             []string
	deletion          atomic.Bool
}

func (m *GkeTokenSource) getCachedToken(ctx context.Context, k8sNamespace string, k8sServiceAccount string) (*oauth2.Token, bool) {
	key := tokensKey{
		k8sNamespace:      k8sNamespace,
		k8sServiceAccount: k8sServiceAccount,
	}
	v, ok := m.tokens.Load(key)

	if !ok {
		return nil, false
	}

	token := v.(*oauth2.Token)

	if !hasExpired(token) {
		return token, true
	} else {
		// poor man's TTL cache, cache size shouldn't be a concern because it is limited by the number of SAs
		// if token has expired, clean other expired tokens, otherwise, nobody will clean them
		if m.deletion.CompareAndSwap(false, true) {
			m.tokens.Range(func(rawKey, value interface{}) bool {
				key := rawKey.(tokensKey)
				token := value.(*oauth2.Token)

				if hasExpired(token) {
					logger.Infof(ctx, "Removed expired token for [%s/%s]", key.k8sNamespace, key.k8sServiceAccount)
					m.tokens.Delete(key)
				}

				return true
			})

			m.deletion.Store(false)
		}
	}

	return nil, false
}

func hasExpired(token *oauth2.Token) bool {
	return token.Expiry.Round(0).Add(-defaultGracePeriod).Before(time.Now())
}

func (m *GkeTokenSource) getK8sServiceAccountToken(ctx context.Context, k8sNamespace string, k8sServiceAccount string) (string, error) {
	serviceAccounts := m.kubeClient.CoreV1().ServiceAccounts(k8sNamespace)
	createTokenResponse, err := serviceAccounts.CreateToken(ctx, k8sServiceAccount, &v1.TokenRequest{
		Spec: v1.TokenRequestSpec{
			Audiences: []string{
				m.identityNamespace,
			},
		},
	}, metav1.CreateOptions{})

	if err != nil {
		return "", errors.Wrap(err, "failed to create k8s service account token")
	}

	return createTokenResponse.Status.Token, nil
}

func (m *GkeTokenSource) getGcpServiceAccount(ctx context.Context, k8sNamespace string, k8sServiceAccount string) (string, error) {
	serviceAccounts := m.kubeClient.CoreV1().ServiceAccounts(k8sNamespace)

	serviceAccountResponse, err := serviceAccounts.Get(ctx, k8sServiceAccount, metav1.GetOptions{})

	if err != nil {
		return "", errors.Wrap(err, "failed to create k8s service account token")
	}

	for key, value := range serviceAccountResponse.Annotations {
		if key == gcpServiceAccountAnnotationKey {
			return value, nil
		}
	}

	return "", errors.Errorf(
		"[%v] annotation doesn't exist on k8s service account [%v/%v], read more at %v",
		gcpServiceAccountAnnotationKey,
		k8sNamespace,
		k8sServiceAccount,
		workflowIdentityDocUrl)
}

func (m *GkeTokenSource) GetTokenSource(ctx context.Context, k8sNamespace string, k8sServiceAccount string) (oauth2.TokenSource, error) {
	if k8sServiceAccount == "" {
		k8sServiceAccount = "default"
	}

	if cachedToken, ok := m.getCachedToken(ctx, k8sNamespace, k8sServiceAccount); ok {
		return oauth2.StaticTokenSource(cachedToken), nil
	}

	// when tokens expire, of we hit a miss, singleflight will do a most one request per SA
	value, err := m.singleflight.Do(k8sNamespace+"/"+k8sServiceAccount, func() (interface{}, error) {
		k8sServiceAccountToken, err := m.getK8sServiceAccountToken(ctx, k8sNamespace, k8sServiceAccount)

		if err != nil {
			return nil, err
		}

		gcpServiceAccount, err := m.getGcpServiceAccount(ctx, k8sNamespace, k8sServiceAccount)

		token, err := ExchangeToken(ctx, StsRequest{
			SubjectToken:      k8sServiceAccountToken,
			Scope:             m.scope,
			ServiceAccount:    gcpServiceAccount,
			IdentityNamespace: m.identityNamespace,
		})

		if err != nil {
			return nil, err
		}

		m.tokens.Store(tokensKey{
			k8sNamespace:      k8sNamespace,
			k8sServiceAccount: k8sServiceAccount,
		}, token)

		return token, nil
	})

	if err != nil {
		return nil, err
	}

	return oauth2.StaticTokenSource(value.(*oauth2.Token)), nil
}

func getKubeClient(kubeConfigPath string, kubeConfig KubeClientConfig) (*kubernetes.Clientset, error) {
	var kubecfg *rest.Config
	var err error
	if kubeConfigPath != "" {
		kubeConfigPath := os.ExpandEnv(kubeConfigPath)
		kubecfg, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, errors.Wrapf(err, "Error building kubeconfig")
		}
	} else {
		kubecfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, errors.Wrapf(err, "Cannot get InCluster kubeconfig")
		}
	}

	kubecfg.QPS = kubeConfig.QPS
	kubecfg.Burst = kubeConfig.Burst
	kubecfg.Timeout = kubeConfig.Timeout.Duration

	kubeClient, err := kubernetes.NewForConfig(kubecfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error building kubernetes clientset")
	}
	return kubeClient, err
}

func NewGkeTokenSource(config GKETokenSourceConfig) (GkeTokenSource, error) {
	kubeClient, err := getKubeClient(config.KubeConfigPath, config.KubeConfig)

	if err != nil {
		return GkeTokenSource{}, err
	}

	return GkeTokenSource{
		kubeClient:        kubeClient,
		tokens:            sync.Map{},
		identityNamespace: config.IdentityNamespace,
		scope:             config.Scope,
		singleflight:      singleflight.Group{},
		deletion:          atomic.NewBool(false),
	}, nil
}
