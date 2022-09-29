package main

import (
    "flag"
	"fmt"
    "io"
	"os"

    "crypto/sha256"
	"encoding/json"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
	benigma "github.com/ixe013/benigma"
)

func computeSha256OfFile(exe string) (string, error) {
	file, err := os.Open(exe)
	if err == nil {
        defer file.Close()

        hash := sha256.New()
        _, err := io.Copy(hash, file)

        if err == nil {
            return fmt.Sprintf("%x", hash.Sum(nil)), nil
        }
    }

    return "unable to compute hash", err
}

func main() {
    flag.Parse()

    for _, fl := range flag.Args() {
        exe, _ := os.Executable()
        hash, _ := computeSha256OfFile(exe)

        if (fl == "version") {
            version := map[string]string{
                "version": benigma.Version,
                "commit":  benigma.Commit,
                "sha256":  hash,
            }
            bytes, err := json.Marshal(version)

            if err == nil {
                fmt.Println(string(bytes))
            }

	        os.Exit(0)

        } else if (fl == "hash") {
            fmt.Println(hash)
	        os.Exit(0)
        }
    }

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

	os.Exit(0)
}
