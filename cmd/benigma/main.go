package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
	benigma "github.com/vaups/benigma"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{})

	logger.Info("Enigma secret backend starting")

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: benigma.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})

	if err != nil {
		logger.Error("Enigma secret backent shutting down", "error", err)
		os.Exit(1)
	}

	os.Exit(0)
}
