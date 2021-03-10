package google

import "github.com/flyteorg/flytestdlib/config"

//go:generate pflags GKETokenSourceConfig --default-var=defaultConfig

type GKETokenSourceConfig struct {
	// IdentityNamespace is workload identity namespace, e.g. [project_id].svc.id.goog
	IdentityNamespace string `json:"identityNamespace" pflag:",Defines workload identity namespace, e.g. [project_id].svc.id.goog"`

	// Scope is OAuth 2.0 scopes to include on the resulting access token
	Scope []string `json:"scope" pflag:",Defines OAuth 2.0 scopes to include on the resulting access token"`

	// KubeConfigPath is path to the Kubernetes client config file
	KubeConfigPath string `json:"kubeConfig" pflag:",Path to Kubernetes client config file."`

	// KubeConfig is configuration to the Kubernetes client
	KubeConfig KubeClientConfig `json:"kubeClientConfig" pflag:",Configuration to control the Kubernetes client"`
}

type KubeClientConfig struct {
	// QPS indicates the maximum QPS to the master from this client.
	// If it's zero, the created RESTClient will use DefaultQPS: 5
	QPS float32 `json:"qps" pflag:"-,Max QPS to the master for requests to KubeAPI. 0 defaults to 5."`

	// Maximum burst for throttle.
	// If it's zero, the created RESTClient will use DefaultBurst: 10.
	Burst int `json:"burst" pflag:"-,Max burst rate for throttle. 0 defaults to 10"`

	// The maximum length of time to wait before giving up on a server request. A value of zero means no timeout.
	Timeout config.Duration `json:"timeout" pflag:"-,Max duration allowed for every request to KubeAPI before giving up. 0 implies no timeout."`
}
