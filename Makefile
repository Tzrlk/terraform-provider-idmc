default: build

.build/:
	mkdir -p ${@}

CMD_OAPI_CODEGEN ?= go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
CMD_TFPLUGINDOCS ?= go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

################################################################################
#: Generate OpenAPI clients
codegen: \
		internal/idmc/admin/v2/idmc-admin-v2.gen.go \
		internal/idmc/admin/v3/idmc-admin-v3.gen.go
.PHONY: gen

%.gen.go: VPATH = internal/idmc
%.gen.go: %.yml
	${CMD_OAPI_CODEGEN} \
			-generate types,client,spec \
			-package $(lastword $(subst /, ,$(dir ${@}))) \
			-o ${@} \
			${<}

################################################################################
#: Format module sources.
format:
	go fmt
	terraform fmt -recursive ./examples/

################################################################################
#: Lint project
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

################################################################################
#: Compile module.
build: \
		codegen
	go mod tidy
	go install
.PHONY: build

################################################################################
#: Run unit tests
test:
	go test

################################################################################
#: Generate documentation.
docs: docs/*
.PHONY: docs

docs/*: \
	docs/data-sources/* \
	docs/resources/* \
	docs/functions/*

docs/data-sources/* docs/resources/* docs/functions/* &: \
		*.go \
		$(shell find examples/ -name *.tf) \
		*.md \
		| .build/
	${CMD_TFPLUGINDOCS} \
		generate \
		--provider-name idmc \
		--website-temp-dir .build/tfplugindocs
#		--providers-schema ??? (need to generate first)

################################################################################
#: Run acceptance tests.
verify:
	TF_ACC=1 go test ./... -v ${TESTARGS} -timeout 120m
.PHONY: verify
