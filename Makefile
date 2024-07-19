
# Settings
.DEFAULT: help
.DELETE_ON_ERROR:
.EXPORT_ALL_VARIABLES:
.ONESHELL:
.ALWAYS:

## CLEAR BUILT-INS #############################################################
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables
.SUFFIXES:

## INCLUDE LOCAL FILES #########################################################
MAKE_LOCALS := $(wildcard *local.mk)
$(if ${MAKE_LOCALS},$(foreach local,${MAKE_LOCALS},$(eval include ${local})),)

## HELP LISTS AVAILABLE TARGETS ################################################
help:
	@grep -E "^\s*include \S+$$" Makefile \
		| cut -f2 -d" " \
		| cat <(printf "%s\n" "Makefile") - \
		| xargs grep --no-filename -B1 -E "^[a-zA-Z0-9_-]+\:([^\=]|$$)" \
		| grep -v -- -- \
		| sed "N;s|\n|###|" \
		| sed -n "s|^#: \(.*\)###\(.*\):.*|\2###\1|p" \
		| column -t  -s '###'

.build/:
	mkdir -p ${@}

#: Wipes out any and all untracked and ignored files.
clobber:
	git clean -dx \
			--force

# Binary commands in case they need to be overridden.
CMD_OAPI_CODEGEN  ?= go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
CMD_TFPLUGINDOCS  ?= go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
CMD_GOLANGCI_LINT ?= go run github.com/golangci/golangci-lint/cmd/golangci-lint

################################################################################

API_SRC_FILES := $(wildcard internal/idmc/*/openapi.yml)
API_SRC_DIRS  := $(sort $(dir ${API_SRC_FILES}))

GO_SRC_FILES := $(shell find . -type f -name *.go -not -name *_test.go -not -name *.gen.go)
GO_SRC_DIRS  := $(sort $(dir ${GO_SRC_FILES}))
vpath %.go ${GO_SRC_DIRS}

GO_TEST_FILES := $(shell find . -type f -name *_test.go)
GO_TEST_DIRS  := $(sort $(dir ${GO_TEST_FILES}))
vpath %_test.go ${GO_TEST_DIRS}

TF_SRC_FILES := $(shell find . -type f -name *.tf)
vpath %.tf      examples/
vpath %.tfvars  examples/
vpath %.tfstate examples/

################################################################################
#: Generate OpenAPI clients
codegen: $(patsubst %/openapi.yml,%/client.gen.go,${API_SRC_FILES})
.PHONY: gen

%/client.gen.go: VPATH = internal/idmc
%/client.gen.go: \
		%/openapi.yml \
		%/codegen.yml \
		go.sum
	${CMD_OAPI_CODEGEN} \
			-config $(word 2,${^}) \
			${<}

################################################################################
#: Format module sources.
format: \
		.build/go-fmt.done \
		.build/tf-fmt.done
.PHONY: format

.build/go-fmt.done: \
		${GO_SRC_FILES} \
		${GO_TEST_FILES}
	go fmt $(foreach src_dir,${?D},./${src_dir}) \
	&& touch ${@}

.build/tf-fmt.done: \
		${TF_SRC_FILES}
	terraform fmt $(foreach src_dir,${?D},./${src_dir}) \
		-diff \
	&& touch ${@}

################################################################################
#: Lint project
lint:
	${CMD_GOLANGCI_LINT} run
.PHONY: lint

################################################################################
#: Tidy module dependencies
tidy: go.sum
.PHONY: tidy

# NOTE: touch needed at the end due to other files being manipulated.
go.sum: \
		go.mod
	go mod tidy -v && \
	touch ${@}

################################################################################
#: Compile module.
install: ${GOPATH}/bin/terraform-provider-idmc.exe
.PHONY: install

${GOPATH}/bin/terraform-provider-idmc.exe: \
		${GO_SRC_FILES} \
		go.sum \
		codegen
	go install

################################################################################
#: Run unit tests
test: .build/gotest.jsonl
.PHONY: test

.build/gotest.jsonl: \
		${GO_TEST_FILES} \
		| .build/
	go test -v -json ${GO_TEST_DIRS} \
		| tee ${@}

%_test.go: %.go
	touch ${@}
.NOTINTERMEDIATE: %_test.go

################################################################################
#: Generate documentation.
docs: docs/*
.PHONY: docs

docs/*: \
	docs/data-sources/* \
	docs/resources/* \
	docs/functions/*

docs/data-sources/* docs/resources/* docs/functions/* &: \
		${GO_SRC_FILES} \
		${TF_SRC_FILES} \
		*.md \
		| .build/
	${CMD_TFPLUGINDOCS} \
		generate \
		--provider-name idmc \
		--website-temp-dir .build/tfplugindocs
#		--providers-schema ??? (need to generate first)

# Not really true, tbh, but somewhat related.
examples/data-sources/idmc_%/data-source.tf: \
		internal/provider/%_data.go
	touch ${@}
examples/resources/idmc_%/resource.tf: \
		internal/provider/%.go
	touch ${@}
examples/provider/provider.tf: \
		internal/provider/provider.go
	touch ${@}

################################################################################
#: Run acceptance tests.
verify:
	TF_ACC=1 go test ./... -v ${TESTARGS} -timeout 120m
.PHONY: verify
