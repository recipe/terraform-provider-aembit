default: testacc

install:
	go get aembit.io/aembit
	go install -a -ldflags "-X main.version=1.0.0" .

# Run acceptance tests
.PHONY: testacc
testacc: install
	cd internal/provider
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 10m
