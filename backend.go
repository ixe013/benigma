package benigma

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Factory configures and returns Enigma backends
func Factory(contex context.Context, configuration *logical.BackendConfig) (logical.Backend, error) {
	backend, err := newEnigmaBackend(contex, configuration)

	if err != nil {
		return nil, err
	}

	if err := backend.Setup(contex, configuration); err != nil {
		return nil, err
	}

	return backend, nil
}

func newEnigmaBackend(ctx context.Context, conf *logical.BackendConfig) (*enigmaBackend, error) {
	var b enigmaBackend
	b.Backend = &framework.Backend{
		Help: `An Enigma machine implemented as a Vault Secret Engine.`,
		Paths: []*framework.Path{

			b.getPathForModel(),
			b.getPathForSpecificModelOperations(),

			b.getPathForInstances(),
			b.getPathForInstanceCreation(),
			b.getPathForInstanceOperations(),
		},

		BackendType: logical.TypeLogical,
	}

	return &b, nil
}

type enigmaBackend struct {
	*framework.Backend
}
