package google

import (
	"context"
	"os"
	"sync"
	"time"

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
)

const (
	defaultGracePeriod             = 300
	gcpServiceAccountAnnotationKey = "iam.gke.io/gcp-service-account"
	workflowIdentityDocURL         = "https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity"
)

type gkeTokenSource struct {
	kubeClient             kubernetes.Interface
	tokens                 *sync.Map
	singleflight           *singleflight.Group
	identityNamespace      string
	scope                  []string
	deletion               *atomic.Bool
	gkeClusterURL          string
	federatedTokenEndpoint string // only for testing
	iamCredentialsEndpoint string // only for testing
}

func (m *gkeTokenSource) getCachedToken(ctx context.Context, identity Identity) (*oauth2.Token, bool) {
	v, ok := m.tokens.Load(identity)

	if !ok {
		return nil, false
	}

	token := v.(*oauth2.Token)

	if !hasExpired(token) {
		return token, true
	}

	// poor man's TTL cache, cache size shouldn't be a concern because it is limited by the number of SAs
	// if token has expired, clean other expired tokens, otherwise, nobody will clean them
	if m.deletion.CompareAndSwap(false, true) {
		m.tokens.Range(func(rawKey, value interface{}) bool {
			identity := rawKey.(Identity)
			token := value.(*oauth2.Token)

			if hasExpired(token) {
				logger.Infof(ctx, "Removed expired token for [%s/%s]", identity.K8sNamespace, identity.K8sServiceAccount)
				m.tokens.Delete(identity)
			}

			return true
		})

		m.deletion.Store(false)
	}

	return nil, false
}

func hasExpired(token *oauth2.Token) bool {
	return token.Expiry.Round(0).Add(-defaultGracePeriod).Before(time.Now())
}

func (m *gkeTokenSource) getK8sServiceAccountToken(ctx context.Context, identity Identity) (string, error) {
	serviceAccounts := m.kubeClient.CoreV1().ServiceAccounts(identity.K8sNamespace)
	createTokenResponse, err := serviceAccounts.CreateToken(ctx, identity.K8sServiceAccount, &v1.TokenRequest{
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

func (m *gkeTokenSource) getGcpServiceAccount(ctx context.Context, identity Identity) (string, error) {
	serviceAccounts := m.kubeClient.CoreV1().ServiceAccounts(identity.K8sNamespace)

	serviceAccountResponse, err := serviceAccounts.Get(ctx, identity.K8sServiceAccount, metav1.GetOptions{})

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
		identity.K8sNamespace,
		identity.K8sServiceAccount,
		workflowIdentityDocURL)
}

func (m gkeTokenSource) GetTokenSource(ctx context.Context, identity Identity) (oauth2.TokenSource, error) {
	if identity.K8sServiceAccount == "" {
		identity.K8sServiceAccount = "default"
	}

	if cachedToken, ok := m.getCachedToken(ctx, identity); ok {
		return oauth2.StaticTokenSource(cachedToken), nil
	}

	// when tokens expire, of we hit a miss, singleflight will do a most one request per SA
	value, err := m.singleflight.Do(identity.K8sNamespace+"/"+identity.K8sServiceAccount, func() (interface{}, error) {
		k8sServiceAccountToken, err := m.getK8sServiceAccountToken(ctx, identity)

		if err != nil {
			return nil, err
		}

		gcpServiceAccount, err := m.getGcpServiceAccount(ctx, identity)

		if err != nil {
			return nil, err
		}

		token, err := ExchangeToken(ctx, StsRequest{
			SubjectToken:           k8sServiceAccountToken,
			Scope:                  m.scope,
			ServiceAccount:         gcpServiceAccount,
			GKEClusterURL:          m.gkeClusterURL,
			IdentityNamespace:      m.identityNamespace,
			iamCredentialsEndpoint: m.iamCredentialsEndpoint,
			federatedTokenEndpoint: m.federatedTokenEndpoint,
		})

		if err != nil {
			return nil, err
		}

		m.tokens.Store(identity, token)

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

func newGKETokenSource(kubeClient kubernetes.Interface, config TokenSourceFactoryConfig) gkeTokenSource {
	deletion := atomic.NewBool(false)

	return gkeTokenSource{
		kubeClient:        kubeClient,
		tokens:            &sync.Map{},
		identityNamespace: config.IdentityNamespace,
		gkeClusterURL:     config.GKEClusterURL,
		scope:             config.Scope,
		singleflight:      &singleflight.Group{},
		deletion:          &deletion,
	}
}

func NewGKETokenSource(config TokenSourceFactoryConfig) (TokenSourceFactory, error) {
	kubeClient, err := getKubeClient(config.KubeConfigPath, config.KubeClientConfig)

	if err != nil {
		return gkeTokenSource{}, err
	}

	if config.IdentityNamespace == "" {
		return gkeTokenSource{}, errors.New("bigquery.googleTokenSource.identityNamespace is required when bigquery.googleTokenSource.type is 'gke'")
	}

	return newGKETokenSource(kubeClient, config), nil
}
