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


debug: GOBUILDFLAGS = -gcflags "all=-N -l"
debug: build

all: fmt debug dev

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault/plugins/benigma $(GOBUILDFLAGS) cmd/benigma/main.go

dev:
	vault server --dev --dev-root-token-id root --log-level trace --dev-plugin-dir=$$(pwd -P)/vault/plugins

enable:
	vault secrets enable benigma

clean:
	rm -vf ./vault/plugins/benigma

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt dev enable
