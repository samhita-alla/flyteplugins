package google

import (
	"context"
	"golang.org/x/oauth2"
)

type Identity struct {
	K8sNamespace      string
	K8sServiceAccount string
}

type TokenSource interface {
	GetTokenSource(ctx context.Context, identity Identity) (oauth2.TokenSource, error)
}
