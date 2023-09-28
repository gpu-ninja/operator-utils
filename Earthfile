VERSION 0.7
FROM golang:1.21-bookworm
WORKDIR /app

generate:
  FROM +tools
  COPY . .
  RUN controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
  SAVE ARTIFACT ./reference/zz_generated.deepcopy.go AS LOCAL reference/zz_generated.deepcopy.go

tidy:
  LOCALLY
  RUN go mod tidy
  RUN go fmt ./...

lint:
  FROM golangci/golangci-lint:v1.54.2
  WORKDIR /app
  COPY . ./
  RUN golangci-lint run --timeout 5m ./...

test:
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  RUN go test -coverprofile=coverage.out -v ./...
  SAVE ARTIFACT ./coverage.out AS LOCAL coverage.out

tools:
  ARG CONTROLLER_TOOLS_VERSION=v0.12.0
  RUN go install sigs.k8s.io/controller-tools/cmd/controller-gen@${CONTROLLER_TOOLS_VERSION}