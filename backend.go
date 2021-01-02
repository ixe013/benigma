package benigma

import (
	"context"
	//"encoding/json"
	//"fmt"
	//"strings"

	//"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	//"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// Factory configures and returns Enigma backends
func Factory(contex context.Context, configuration *logical.BackendConfig) (logical.Backend, error) {
	backend, err := Backend(contex, configuration)

	if err != nil {
		return nil, err
	}

	if err := backend.Setup(contex, configuration); err != nil {
		return nil, err
	}

	return backend, nil
}

func Backend(ctx context.Context, conf *logical.BackendConfig) (*backend, error) {
	var b backend
	b.Backend = &framework.Backend{
		Paths: []*framework.Path{
			// b.pathCreateEnigmaModelInstance(),
			// b.pathProcess(),
			// b.pathReset(),
			// b.pathKeys(),
			b.pathModel(),
			b.pathListModels(),
		},

		BackendType: logical.TypeLogical,
	}

	// determine cacheSize to use. Defaults to 0 which means unlimited
	return &b, nil
}

type backend struct {
	*framework.Backend
}
