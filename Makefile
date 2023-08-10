# Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin

# Tool Versions
CONTROLLER_TOOLS_VERSION ?= v0.12.0

# Tool Binaries
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen

generate: $(CONTROLLER_GEN)
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

tidy:
	go mod tidy
	go fmt ./...

lint:
	golangci-lint run ./...

test:
	go test -coverprofile=coverage.out -v ./...

clean:
	-rm -rf bin
	go clean -testcache

$(LOCALBIN):
	mkdir -p $(LOCALBIN)

$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen && $(LOCALBIN)/controller-gen --version | grep -q $(CONTROLLER_TOOLS_VERSION) || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: generate tidy lint test clean