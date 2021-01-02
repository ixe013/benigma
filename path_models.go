package benigma

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathListModels() *framework.Path {
	return &framework.Path{
		Pattern: "models/?$",

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathModelList,
				Summary:  "List the available Enigma machie models",
			},
			/*
						logical.UpdateOperation: &framework.PathOperation{
							Callback: b.pathRoleCreateUpdate,
						},
						logical.ReadOperation: &framework.PathOperation{
							Callback: b.pathRoleRead,
						},
						logical.DeleteOperation: &framework.PathOperation{
							Callback: b.pathRoleDelete,
						},
			//*/
		},

		HelpSynopsis:    pathModelHelpSyn,
		HelpDescription: pathModelHelpDesc,
	}
}

func (b *backend) pathModel() *framework.Path {
	return &framework.Path{
		Pattern: "models",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The name of the Enigma model you are creating. Cannot be one of the build-in model names.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathModelCreate,
				Summary:  "Creates an new Enigma machine model",
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathModelCreate,
				Summary:  "Creates an new Enigma machine model",
			},
			/*
						logical.UpdateOperation: &framework.PathOperation{
							Callback: b.pathRoleCreateUpdate,
						},
						logical.ReadOperation: &framework.PathOperation{
							Callback: b.pathRoleRead,
						},
						logical.DeleteOperation: &framework.PathOperation{
							Callback: b.pathRoleDelete,
						},
			//*/
		},

		HelpSynopsis:    pathModelHelpSyn,
		HelpDescription: pathModelHelpDesc,
	}
}

func (b *backend) pathModelList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "models/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(append(BuiltinModelNames(), entries...)), nil
}

func (b *backend) pathModelCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	name := data.Get("name").(string)

	model := &modelEntry{
		Name: name,
	}

	// Convert the model to JSON to store it
	jsonEntry, err := logical.StorageEntryJSON("models/"+name, model)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, jsonEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

type modelEntry struct {
	Name string `json:name`
}

const pathModelHelpSyn = `Operations on a given Enigma machine model`

const pathModelHelpDesc = `
This path is used to manage the Enigma machines and process plaintext or
ciphertext with them.
`
