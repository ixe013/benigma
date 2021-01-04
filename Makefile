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
PLUGINNAME = benigma

debug: GOBUILDFLAGS = -gcflags "all=-N -l"
debug: build

all: fmt debug dev

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o $(OUTPUTFOLDER)/$(PLUGINNAME) $(GOBUILDFLAGS) cmd/$(PLUGINNAME)/main.go
	sha256sum $(OUTPUTFOLDER)/$(PLUGINNAME)

dev:
	vault server --dev --dev-root-token-id root --log-level trace --dev-plugin-dir=$$(pwd -P)/$(OUTPUTFOLDER)

register:
	vault secrets disable $(PLUGINNAME)
	vault plugin deregister $(PLUGINNAME)
	vault plugin register --sha256=$$(sha256sum $(OUTPUTFOLDER)/$(PLUGINNAME)|cut -f1 -d " ") $(PLUGINNAME)
	vault secrets enable $(PLUGINNAME)
	vault path-help $(PLUGINNAME)

enable:
	vault secrets enable $(PLUGINNAME)

clean:
	rm -vf ./vault/plugins/$(PLUGINNAME)

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt dev enable
