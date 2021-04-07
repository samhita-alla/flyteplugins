package google

import "github.com/flyteorg/flytestdlib/config"

type TokenSourceFactoryType = string

const (
	TokenSourceTypeDefault = "default"
	TokenSourceTypeGKE = "gke"
)

type TokenSourceFactoryConfig struct {
	// Type is type of TokenSourceFactory, possible values are 'default' or 'gke'.
	// - 'default' uses default credentials, see https://cloud.google.com/iam/docs/service-accounts#default
	// - 'gke' uses GKE workload identity, see https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity
	Type TokenSourceFactoryType `json:"type" pflag:",Defines type of TokenSourceFactory, possible values are 'default' or 'gke'"`

	// IdentityNamespace is workload identity namespace, e.g. [project_id].svc.id.goog
	IdentityNamespace string `json:"identityNamespace" pflag:",Defines workload identity namespace, e.g. [project_id].svc.id.goog"`

	// GKEClusterURL is URL to GKE cluster, e.g. https://container.googleapis.com/v1/projects/<project>/locations/<location>/clusters/<cluster>
	GKEClusterURL string `json:"gkeClusterURL" pflag:",Defines URL to GKE cluster, e.g. https://container.googleapis.com/v1/projects/<project>/locations/<location>/clusters/<cluster>"`

	// Scope is OAuth 2.0 scopes to include on the resulting access token
	Scope []string `json:"scope" pflag:",Defines OAuth 2.0 scopes to include on the resulting access token"`

	// KubeConfigPath is path to the Kubernetes client config file
	KubeConfigPath string `json:"kubeConfig" pflag:",Path to Kubernetes client config file."`

	// KubeClientConfig is configuration to the Kubernetes client
	KubeClientConfig KubeClientConfig `json:"kubeClientConfig" pflag:",Configuration to control the Kubernetes client"`
}

type KubeClientConfig struct {
	// QPS indicates the maximum QPS to the master from this client.
	// If it's zero, the created RESTClient will use DefaultQPS: 5
	QPS float32 `json:"qps" pflag:"-,Max QPS to the master for requests to KubeAPI. 0 defaults to 5."`

	// Maximum burst for throttle.
	// If it's zero, the created RESTClient will use DefaultBurst: 10.
	Burst int `json:"burst" pflag:"-,Max burst rate for throttle. 0 defaults to 10"`

	// The maximum length of time to wait before giving up on a server request. A value of zero means no timeout.
	Timeout config.Duration `json:"timeout" pflag:",Max duration allowed for every request to KubeAPI before giving up. 0 implies no timeout."`
}

func GetDefaultConfig() TokenSourceFactoryConfig {
	return TokenSourceFactoryConfig{
		Type: "default",
		KubeClientConfig: KubeClientConfig{
			QPS:     5,
			Burst:   10,
			Timeout: config.Duration{Duration: 30},
		},
	}
}
