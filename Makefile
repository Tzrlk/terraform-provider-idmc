default: build

CMD_OAPI_CODEGEN ?= go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen

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

#: Compile module.
build: \
		codegen
	go install
.PHONY: build

#: Generate documentation.
docs:
	go generate
.PHONY: docs

#: Run acceptance tests.
verify:
	TF_ACC=1 go test ./... -v ${TESTARGS} -timeout 120m
.PHONY: verify
