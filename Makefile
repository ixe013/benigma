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

all: fmt build start

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault/plugins/benigma -gcflags="all=-N -l" cmd/benigma/main.go

start:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=$$(pwd -P)/vault/plugins

enable:
	vault secrets enable benigma

clean:
	rm -f ./vault/plugins/benigma

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt start enable