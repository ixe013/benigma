GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all

OUTPUTFOLDER = ./vault/plugins
PROJECTNAME = benigma
PLUGINNAME = enigma
COMMIT=$(shell git rev-parse --short HEAD)
OUTPUTNAME = $(PLUGINNAME).$(COMMIT)

debug: GOBUILDFLAGS = -gcflags "all=-N -l"
debug: build

all: debug register test

upgrade: register reload

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o "$(OUTPUTFOLDER)/$(OUTPUTNAME)" $(GOBUILDFLAGS) -ldflags="-X 'github.com/vaups/benigma.Version=$$(git describe --abbrev=0 --dirty=+)' -X 'github.com/vaups/benigma.Commit=$(COMMIT)'" cmd/$(PROJECTNAME)/main.go
	sha256sum $(OUTPUTFOLDER)/$(OUTPUTNAME)

dev:
	vault server --dev --dev-root-token-id root --log-level trace --dev-plugin-dir=$$(pwd -P)/$(OUTPUTFOLDER)

test:
	find . -type f -name "*.shunit2" -exec shunit2 {} \;

unregister:
	vault secrets disable $(PLUGINNAME)
	vault plugin deregister $(PLUGINNAME)

register:
	curl -v --request PUT $$VAULT_ADDR/v1/sys/plugins/catalog/secret/$(PLUGINNAME) --header "X-Vault-Token: $$(vault print token)" --data "{ \"type\":\"secret\", \"command\":\"$(OUTPUTNAME)\", \"sha256\":\"$$(sha256sum $(OUTPUTFOLDER)/$(OUTPUTNAME)|cut -f1 -d ' ')\" }"

reload:
	vault write sys/plugins/reload/backend plugin=$(PLUGINNAME) scope=global mounts=$(PLUGINNAME)

clean:
	vault secrets disable $(PLUGINNAME) || true
	vault plugin deregister $(PLUGINNAME) || true
	rm -vf ./vault/plugins/*

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt dev enable
