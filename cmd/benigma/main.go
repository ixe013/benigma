package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
	benigma "github.com/vaups/benigma"
)

func main() {
	if (len(os.Args) >= 2) && (os.Args[1] == "version") {
		fmt.Println(benigma.Version)
	} else {
		logger := hclog.New(&hclog.LoggerOptions{})

		logger.Info("Enigma secret engine starting", "version", benigma.Version, "commit", benigma.Commit)

		apiClientMeta := &api.PluginAPIClientMeta{}
		flags := apiClientMeta.FlagSet()
		flags.Parse(os.Args[1:])

		tlsConfig := apiClientMeta.GetTLSConfig()
		tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

		err := plugin.Serve(&plugin.ServeOpts{
			BackendFactoryFunc: benigma.Factory,
			TLSProviderFunc:    tlsProviderFunc,
		})

		logger.Error("Enigma secret engine shutting down", "error", err)

		if err != nil {
			os.Exit(1)
		}
	}

	os.Exit(0)
}
