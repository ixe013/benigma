package benigma

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *enigmaBackend) getPathForModel() *framework.Path {
	return &framework.Path{
		Pattern: "models/?$",

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The name of the Enigma model you are creating. Cannot be one of the build-in model names.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.listModels,
				Summary:  "List the available Enigma machine models",
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.createModel,
				Summary:  "Creates an new custom Enigma machine model",
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.createModel,
				Summary:  "Creates an new custom Enigma machine model",
			},
		},

		HelpSynopsis:    pathModelHelpSyn,
		HelpDescription: pathModelHelpDesc,
	}
}

func (b *enigmaBackend) getPathForSpecificModelOperations() *framework.Path {
	return &framework.Path{
		Pattern: "models/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The name of the Enigma model you are creating, or deleting.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.deleteModel,
				Summary:  "Deletes a custom Enigma machine model (builtin models can't be deleted).",
			},
		},

		HelpSynopsis:    pathModelHelpSyn,
		HelpDescription: pathModelHelpDesc,
	}
}

func (b *enigmaBackend) listModels(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "models/")
	if err != nil {
		b.Logger().Error("An error occured trying to list the models")
		return logical.ErrorResponse("Internal error, unable to list the current models"), err
	}

	return logical.ListResponse(append(builtinModelNames(), entries...)), nil
}

func (b *enigmaBackend) createModel(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	model := &modelEntry{
		Name: name,
	}

	// Convert the model to JSON to store it
	jsonEntry, err := logical.StorageEntryJSON("models/"+name, model)
	if err != nil {
		b.Logger().Error("An error occured trying to create the model")
		return logical.ErrorResponse("An error occured trying to create the model"), err
	}
	if err := req.Storage.Put(ctx, jsonEntry); err != nil {
		b.Logger().Error("An error occured trying to save the model")
		return logical.ErrorResponse("An error occured trying to save the model"), err
	}

	return nil, nil
}

func (b *enigmaBackend) deleteModel(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name := data.Get("name").(string)

	//Check to see if the instance name is builtin (and therefore can't be deleted)
	_, protected := builtinModels[name]

	if protected {
		return logical.ErrorResponse("Cannot delete builtin model " + name), logical.ErrInvalidRequest
	}

	if err := req.Storage.Delete(ctx, "models/"+name); err != nil {
		b.Logger().Error("An error occured trying to delete the model")
		return logical.ErrorResponse("An error occured trying to delete the model"), err
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
