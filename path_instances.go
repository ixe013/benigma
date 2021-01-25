package benigma

import (
	"bytes"
	"context"
	"encoding/gob"
	"strings"

	"github.com/emedvedev/enigma"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *enigmaBackend) getPathForInstances() *framework.Path {
	return &framework.Path{
		Pattern: "instances/?$",

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.listInstances,
				Summary:  "List the available Enigma machine instances",
			},
		},

		HelpSynopsis: `List individual instances of Enigma machines.`,
	}
}

func (b *enigmaBackend) getPathForInstanceOperations() *framework.Path {
	return &framework.Path{
		Pattern: "instances/" + framework.GenericNameRegex("id") + "$",

		Fields: map[string]*framework.FieldSchema{
			"id": {
				Type:        framework.TypeString,
				Description: "The id of the Enigma model instance you are creating. Must be unique. Think of it as a machine's serial number.",
			},
			"keyboard": {
				Type:        framework.TypeString,
				Description: "The text that is typed on the Enigma machine's keyboard.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.readInstance,
				Summary:  "Returns the state of an Enigma machine instances",
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.encodeText,
				Summary:  "Encodes (encrypts or decrypts) text with this Enigma machine instances",
			},
		},

		HelpSynopsis: `Operations on an individual instance of an Enigma machine.`,
		HelpDescription: `This path is used to actually use an Enigma machines to encrypt (or decrypt) text. The state
of the machine is saved across invocations, as if you are communicating with someone on the other end.`,
	}
}

func (b *enigmaBackend) getPathForInstanceCreation() *framework.Path {
	return &framework.Path{
		Pattern: "models/" + framework.GenericNameRegex("model") + "/instance/?$",
		Fields: map[string]*framework.FieldSchema{
			"model": {
				Type:        framework.TypeString,
				Description: "The name of the Enigma model you want to produce an instance of.",
			},
			"id": {
				Type:        framework.TypeString,
				Description: "The id of the Enigma model instance you are creating. Must be unique. Think of it as a machine's serial number.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.createInstance,
				Summary:  "Process some text through this Engima instance. Will either encrypt or decrypt, depending on what you send it",
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.createInstance,
				Summary:  "Process some text through this Engima instance. Will either encrypt or decrypt, depending on what you send it",
			},
		},

		HelpSynopsis: `Factory of individual instances of a given Enigma machine model`,
		HelpDescription: `This path is used to create (and list) working copies of the Enigma machines, as if they were to be dispached to the
theater of operations. An instance will retain its state across Vault restart. Typically, the sender and 
receiver of an encrypted message will each have an instance of the same model.`,
	}
}

func (b *enigmaBackend) listInstances(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "instances/")

	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func serializeEnigmaInstance(machine *enigma.Enigma, logger hclog.Logger) ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)

	err := enc.Encode(machine)

	if err != nil {
		logger.Error("Unable to serialize enigma instnace", "len", buf.Len(), "error", err)
	}

	return buf.Bytes(), err
}

func deserializeEnigmaInstance(raw []byte) (*enigma.Enigma, error) {
	var machine enigma.Enigma
	b := bytes.Buffer{}
	b.Write(raw)

	d := gob.NewDecoder(&b)

	err := d.Decode(&machine)

	return &machine, err
}

func getStatefulInstanceFromStorage(id string, ctx context.Context, req *logical.Request, logger hclog.Logger) (*statefulInstanceEntry, error) {
	logger.Trace("Looking for instance", "id", id)

	jsonEntry, err := req.Storage.Get(ctx, "instances/"+id)

	if err != nil {
		logger.Error("Instance not found in cache", "id", id)
		return nil, err
	}

	var instance statefulInstanceEntry

	err = jsonEntry.DecodeJSON(&instance)

	if err != nil {
		logger.Error("Unable to deserialize instance found in cache")
		return nil, err
	}

	return &instance, err
}

func (b *enigmaBackend) createInstance(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	model := data.Get("model").(string)
	id := data.Get("id").(string)
	original, exists := builtinModels[model]

	if !exists {
		return logical.ErrorResponse("Cannot only create instances of builtin models (for now)"), logical.ErrInvalidRequest
	}

	state, err := serializeEnigmaInstance(original, b.Logger())

	if err != nil {
		//Should really be an error 500 here, it is not the client's fault
		return logical.ErrorResponse("Cannot create instances of builtin model"), logical.ErrInvalidRequest
	}

	instance := &statefulInstanceEntry{
		instanceEntry: instanceEntry{
			Model: model,
			Id:    id,
		},
		State: state,
		Steps: 0,
	}

	// Convert the model to JSON to store it
	jsonEntry, err := logical.StorageEntryJSON("instances/"+id, instance)
	if err != nil {
		return nil, err
	}

	/* The Value field of the jsonEntry blob looks like this (with id=z)
	   Offset(h) 00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F
	             -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	   00000000  7B 22 6D 6F 64 65 6C 22 3A 22 49 22 2C 22 69 64  {"model":"I","id
	   00000010  22 3A 22 7A 22 2C 22 73 74 61 74 65 22 3A 22 50  ":"z","state":"P
	   00000020  76 2B 42 41 77 45 42 42 6B 56 75 61 57 64 74 59  v+BAwEBBkVuaWdtY
	   00000030  51 48 2F 67 67 41 42 41 77 45 4A 55 6D 56 6D 62  QH/ggABAwEJUmVmb
	*/

	if err := req.Storage.Put(ctx, jsonEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *enigmaBackend) readInstance(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	id := data.Get("id").(string)

	instance, err := getStatefulInstanceFromStorage(id, ctx, req, b.Logger())

	if err != nil {
		return logical.ErrorResponse("Unable to find or deserialize instance " + id), logical.ErrInvalidRequest
	}

	b.Logger().Debug("Returning instance", "model", instance.Model, "id", instance.Id)

	return &logical.Response{
		Data: map[string]interface{}{
			"model": instance.Model,
			"id":    instance.Id,
			"steps": instance.Steps,
		},
	}, nil
}

func (b *enigmaBackend) encodeText(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	id := data.Get("id").(string)

	instance, err := getStatefulInstanceFromStorage(id, ctx, req, b.Logger())

	if err != nil {
		return logical.ErrorResponse("Unable to read state of cached instance " + id), logical.ErrInvalidRequest
	}

	machine, err := deserializeEnigmaInstance(instance.State)

	if err != nil {
		b.Logger().Error("Deserialization failed", id, err)
		return logical.ErrorResponse("Instance " + id + " has inconsistent state"), logical.ErrInvalidRequest
	}

	keyboard := data.Get("keyboard").(string)

	if keyboard == "" {
		b.Logger().Error("Nothing was typed on the keyboard")
		return logical.ErrorResponse("The field keyboard is mandatory for processing"), logical.ErrInvalidRequest
	}

	uppercase := strings.ToUpper(keyboard)

	sanitized := enigma.SanitizePlaintext(uppercase)

	var warning []string

	if uppercase != sanitized {
		warning = make([]string, 1)
		warning[0] = "Unsupported non-alphabetic characters removed from string"
		b.Logger().Warn(warning[0])
	}

	lights := machine.EncodeString(sanitized)

	b.Logger().Trace("Keyboard processed succesfully", "id", id)

	state, err := serializeEnigmaInstance(machine, b.Logger())

	instance.Steps = instance.Steps + len(keyboard)
	instance.State = state

	// Convert the model to JSON to store it
	jsonEntry, err := logical.StorageEntryJSON("instances/"+id, instance)
	if err != nil {
		b.Logger().Error("Can't convert instance to JSON")
		return logical.ErrorResponse("Internal error, can't convert instance to JSON"), err
	}

	if err := req.Storage.Put(ctx, jsonEntry); err != nil {
		b.Logger().Error("Can't persist JSON instance")
		return logical.ErrorResponse("Internal error, can't perstit JSON instance"), err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"lights": lights,
		},
		Warnings: warning,
	}, nil
}

type instanceEntry struct {
	Model string `json:"model"`
	Id    string `json:"id"`
}

type statefulInstanceEntry struct {
	instanceEntry
	State []byte `json:"state"`
	Steps int    `json:"steps"`
}
