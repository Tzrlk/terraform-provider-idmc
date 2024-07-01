default: build

#: Generate OpenAPI clients
codegen: internal/idmc/client.gen.go
.PHONY: gen

internal/idmc/client.gen.go: \
		tools/oapi-codegen.yml \
		internal/idmc/apiv2.yml \
		internal/idmc/apiv3.yml
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
			--config=${<} \
			internal/idmc/*.yml

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
