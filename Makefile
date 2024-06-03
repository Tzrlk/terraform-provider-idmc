default: testacc

# Compile module.
build:
	go install
.PHONY: build

#: Generate documentation.
docs:
	go generate
.PHONY: docs

#: Run acceptance tests.
verify:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
.PHONY: verify
