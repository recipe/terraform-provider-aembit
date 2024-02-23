default: testacc

# Run the GitHub CI Linters locally
lint: 
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

install:
	go mod tidy
	go install github.com/goreleaser/goreleaser@latest
	go get aembit.io/aembit
	go install -a -ldflags "-X main.version=1.0.0" .

# Run acceptance tests
.PHONY: testacc
testacc: install
	cd internal/provider
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 10m -coverprofile coverage.out
	go tool cover -html coverage.out -o coverage.html

# Locally create a build for local/qa testing using GoReleaser
#	Reference: https://developer.hashicorp.com/terraform/registry/providers/publishing#using-goreleaser-locally
build: testacc
	goreleaser build --snapshot --clean

release: testacc
	goreleaser release