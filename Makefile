#! /usr/bin/make

# Location for the built binaries
DEST_DIR=bin

PACKAGE := $(shell git remote get-url origin | sed -e 's|http://||' -e 's|^.*@||' -e 's|.git||' -e 's|:|/|')
VERSION_PACKAGE=$(PACKAGE)/pkg/cmd/$@
VFLAG=-X '$(VERSION_PACKAGE).name=$@' \
      -X '$(VERSION_PACKAGE).version=$(SYMBOLIC_REF)' \
      -X '$(VERSION_PACKAGE).buildDate=$(DATE)' \
      -X '$(VERSION_PACKAGE).buildSha=$(COMMIT_ID)'
	  
	  
default: all
.PHONY: default

$(NAME): $(shell find . -name \*.go)
    CGO_ENABLED=0 go build $(TAGS) -ldflags "$(VFLAG)" -o build/$@ ./cmd/$@


$(DEST_DIR):
	mkdir -p $(DEST_DIR)
.PHONY: $(DEST_DIR)

#$(DEST_DIR)/test-otel-app: $(DEST_DIR)
#	CGO_ENABLED=0 go build $(TAGS) -ldflags "$(VFLAG)" -o $@ cmd/test-otel-app/main.go

#$(DEST_DIR)/otaas: $(DEST_DIR)
#	CGO_ENABLED=0 go build $(TAGS) -ldflags "$(VFLAG)" -o $@ cmd/otaas/main.go

$(DEST_DIR)/opamp: $(DEST_DIR)
	CGO_ENABLED=0 go build $(TAGS) -ldflags "$(VFLAG)" -o $@ cmd/opamp/main.go

vendor: go.mod go.sum
	go mod vendor
.PHONY: vendor

# binaries: $(DEST_DIR)/test-otel-app  $(DEST_DIR)/otaas
binaries: $(DEST_DIR)/opamp
.PHONY: binaries

build: vendor binaries 

clean:
	rm -rf build .vendor/pkg

fmt:
	go fmt ./...

docker:
#	docker build --no-cache --build-arg validate=0 --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} --build-arg HTTPS_PROXY=${HTTPS_PROXY} --build-arg HTTP_PROXY=${HTTP_PROXY} -f Dockerfile.test-app --tag harishrajagopal/test-otel-app:0.1.0 .
#	docker build --no-cache --build-arg validate=0 --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} --build-arg HTTPS_PROXY=${HTTPS_PROXY} --build-arg HTTP_PROXY=${HTTP_PROXY} -f Dockerfile.otaas --tag harishrajagopal/otaas:0.1.0 .
	docker build --no-cache --build-arg validate=0 --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} --build-arg HTTPS_PROXY=${HTTPS_PROXY} --build-arg HTTP_PROXY=${HTTP_PROXY} -f Dockerfile.opamp --tag harishrajagopal/opamp:0.1.0 .

all: build
.PHONY: all

