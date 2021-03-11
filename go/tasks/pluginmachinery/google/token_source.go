package google

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	DefaultTokenSourceType = "default"
	GKETokenSourceType = "gke"
)

type Identity struct {
	K8sNamespace      string
	K8sServiceAccount string
}

type TokenSource interface {
	GetTokenSource(ctx context.Context, identity Identity) (oauth2.TokenSource, error)
}

func NewTokenSource(config TokenSourceConfig) (TokenSource, error) {
	if config.Type == DefaultTokenSourceType {
		return NewDefaultTokenSource()
	} else if config.Type == GKETokenSourceType {
		return NewGKETokenSource(config)
	} else {
		return nil, errors.New("unknown token source type [%v], possible values are: 'default', 'gke'")
	}
}