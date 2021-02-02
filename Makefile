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
VERSION=$(shell git describe --abbrev=0 --dirty=d)
COMMIT=$(shell git rev-parse --short HEAD)
OUTPUTNAME = $(PLUGINNAME).$(VERSION)

debug: GOBUILDFLAGS = -gcflags "all=-N -l"
debug: build

all: debug register test

upgrade: register reload

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o "$(OUTPUTFOLDER)/$(OUTPUTNAME)" $(GOBUILDFLAGS) -ldflags="-X 'github.com/vaups/benigma.Version=$(VERSION)' -X 'github.com/vaups/benigma.Commit=$(COMMIT)'" cmd/$(PROJECTNAME)/main.go
	sha256sum $(OUTPUTFOLDER)/$(OUTPUTNAME)

dev:
	vault server --dev --dev-root-token-id root --log-level trace --dev-plugin-dir=$$(pwd -P)/$(OUTPUTFOLDER)

test: build
	find . -type f -name "*.shunit2" -exec {} \;

unregister:
	vault secrets disable $(PLUGINNAME)
	vault plugin deregister $(PLUGINNAME)

register: 
	curl -i --request PUT $$VAULT_ADDR/v1/sys/plugins/catalog/secret/$(PLUGINNAME) --header "X-Vault-Token: $$(vault print token)" --data "{ \"type\":\"secret\", \"command\":\"$(OUTPUTNAME)\", \"sha256\":\"$$($(OUTPUTFOLDER)/$(OUTPUTNAME) hash)\" }"

reload:
	vault write sys/plugins/reload/backend plugin=$(PLUGINNAME) scope=global

clean:
	vault secrets disable $(PLUGINNAME) || true
	vault plugin deregister $(PLUGINNAME) || true
	rm -vf ./vault/plugins/*

release: test
	tar czfv enigma.tar.gz --directory=$(OUTPUTFOLDER) $$(ls -1 $(OUTPUTFOLDER) | sort -r | head -1)

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt dev enable
