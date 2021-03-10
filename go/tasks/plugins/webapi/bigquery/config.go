package bigquery

import (
	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/google"
	"time"

	pluginsConfig "github.com/flyteorg/flyteplugins/go/tasks/config"
	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/webapi"
	"github.com/flyteorg/flytestdlib/config"
)

//go:generate pflags Config --default-var=defaultConfig

var (
	defaultConfig = Config{
		WebAPI: webapi.PluginConfig{
			ResourceQuotas: map[core.ResourceNamespace]int{
				"default": 1000,
			},
			ReadRateLimiter: webapi.RateLimiterConfig{
				Burst: 100,
				QPS:   10,
			},
			WriteRateLimiter: webapi.RateLimiterConfig{
				Burst: 100,
				QPS:   10,
			},
			Caching: webapi.CachingConfig{
				Size:              500000,
				ResyncInterval:    config.Duration{Duration: 30 * time.Second},
				Workers:           10,
				MaxSystemFailures: 5,
			},
			ResourceMeta: nil,
		},
		ResourceConstraints: core.ResourceConstraintsSpec{
			ProjectScopeResourceConstraint: &core.ResourceConstraint{
				Value: 100,
			},
			NamespaceScopeResourceConstraint: &core.ResourceConstraint{
				Value: 50,
			},
		},
		TokenSource: "default",
		GKETokenSource: google.GKETokenSourceConfig{
			KubeConfig: google.KubeClientConfig{
				QPS:     5,
				Burst:   10,
				Timeout: config.Duration{Duration: 0},
			},
		},
	}

	configSection = pluginsConfig.MustRegisterSubSection("bigquery", &defaultConfig)
)

type Config struct {
	WebAPI              webapi.PluginConfig          `json:"webApi" pflag:",Defines config for the base WebAPI plugin."`
	ResourceConstraints core.ResourceConstraintsSpec `json:"resourceConstraints" pflag:"-,Defines resource constraints on how many executions to be created per project/overall at any given time."`
	TokenSource         string                       `json:"tokenSource" pflag:",Defines token source: default or GKE"`
	GKETokenSource      google.GKETokenSourceConfig  `json:"gkeTokenSource" pflag:",Defines GKE token source"`
}

func GetConfig() *Config {
	return configSection.GetConfig().(*Config)
}
