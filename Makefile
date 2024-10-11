SHELL = /bin/sh

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
CMD_TERRAFORM     ?= terraform

################################################################################

EXE_OUT := $(realpath ${GOPATH}/bin/terraform-provider-idmc.exe)

# We know that all of the api dirs that'll contain openapi specs and generated
# code are in a set location, and always start with their major version. I could
# also use a wildcard for the source files, but doing it this way enforces an
# expectation of the openapi spec existing in our discovered dirs.
API_SRC_DIRS  := $(sort $(wildcard internal/idmc/v*/))
API_SRC_FILES := $(addsuffix openapi.yml,${API_SRC_DIRS})
API_CFG_FILES := $(addsuffix codegen.yml,${API_SRC_DIRS})
API_GEN_FILES := $(addsuffix client.gen.go,${API_SRC_DIRS})
API_TPL_FILES := $(sort $(wildcard internal/idmc/templates/*.go.tmpl))

# Sadly I'd have to maintain a manual list of all the source package dirs here
# if handling them only with wildcard would work.
GO_SRC_FILES := $(patsubst ./%,%,$(shell find . -type f -name '*.go' -not -name '*_test.go' -not -name '*.gen.go'))
GO_SRC_DIRS  := $(sort $(dir ${GO_SRC_FILES}))

# Tests are always in the same dirs as the src files, so we can just use that.
# Any actual test files should be able to be wildcarded from those dirs.
GO_TEST_FILES := $(wildcard $(addsuffix *_test.go,${GO_SRC_DIRS}))

# Our terraform sources are in very specific directories, so using wildcard like
# this is a lot more efficient than shelling-out to 'find'.
TF_SRC_DIRS_DAT := $(wildcard examples/data-sources/*/)
TF_SRC_DIRS_RES := $(wildcard examples/resources/*/)
TF_SRC_DIRS_FUN := $(wildcard examples/functions/*/)
TF_SRC_DIRS  := $(sort \
	examples/provider/ \
	${TF_SRC_DIRS_DAT} \
	${TF_SRC_DIRS_RES} \
	${TF_SRC_DIRS_FUN} \
)
TF_SRC_FILES := $(wildcard $(addsuffix *.tf,${TF_SRC_DIRS}))
TF_LOG_FILES_DAT := $(addsuffix terraform.jsonl,${TF_SRC_DIRS_DAT})
TF_LOG_FILES_RES := $(addsuffix terraform.jsonl,${TF_SRC_DIRS_RES})
TF_LOG_FILES_FUN := $(addsuffix terraform.jsonl,${TF_SRC_DIRS_FUN})
TF_LOG_FILES := $(sort \
	examples/provider/terraform.jsonl \
	${TF_LOG_FILES_DAT} \
	${TF_LOG_FILES_RES} \
	${TF_LOG_FILES_FUN} \
)

DOC_OUT_TPLS  := $(wildcard templates/*.md.tmpl)
DOC_OUT_FILES := $(sort \
	docs/index.md \
	$(patsubst %/,%.md,$(patsubst examples/data-sources/idmc_%,docs/data-sources/%,${TF_SRC_DIRS_DAT})) \
	$(patsubst %/,%.md,$(patsubst examples/resources/idmc_%,docs/resources/%,${TF_SRC_DIRS_DAT})) \
	$(patsubst %/,%.md,$(patsubst examples/functions/idmc_%,docs/functions/%,${TF_SRC_DIRS_DAT})) \
)

#: Used to debug variable resolution.
debug:
	@echo "${DOC_OUT_FILES}" | tr ' ' '\n'
.PHONY: debug

################################################################################
#: Generate OpenAPI clients
codegen: ${API_GEN_FILES}
.PHONY: codegen

${API_GEN_FILES}: %/client.gen.go: \
		%/openapi.yml \
		%/codegen.yml \
		${API_TPL_FILES} \
		tools/tools.go \
		go.sum
	go generate ${@D}/api.go

################################################################################
#: Format module sources.
format: \
		.build/go-fmt.done \
		.build/tf-fmt.done
.PHONY: format

.build/go-fmt.done: \
		${GO_SRC_FILES} \
		${GO_TEST_FILES} \
		| .build/
	go fmt $(addprefix ./,${?D}) \
	&& touch ${@}

.build/tf-fmt.done: \
		${TF_SRC_FILES} \
		| .build/
	terraform fmt -diff $(addprefix ./,${?D}) \
	&& touch ${@}

################################################################################
#: Lint project
lint: .build/checkstyle.xml
.PHONY: lint

.build/checkstyle.xml: \
		.golangci.yml \
		${GO_SRC_FILES} \
		${GO_TEST_FILES} \
		| .build/
	${CMD_GOLANGCI_LINT} run \
		--config ${<} \
		--issues-exit-code 2 \
		--allow-parallel-runners \
		--out-format checkstyle \
	| tee ${@}

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
install: \
		${EXE_OUT}
.PHONY: install

${EXE_OUT}: \
		${GO_SRC_FILES} \
		${API_GEN_FILES} \
		go.sum
	go install

################################################################################
#: Run unit tests
test: .build/gotest-unit.jsonl
.PHONY: test

################################################################################
#: Run acceptance tests.
verify: .build/gotest-acc.jsonl
.PHONY: verify

# Acceptance tests need a little extra env.
.build/gotest-acc.jsonl: TF_ACC = 1

# Both test types are executed the same.
.build/gotest-%.jsonl: \
		${GO_TEST_FILES} \
		| .build/
	go test -v -json ./... \
    		| tee ${@}

# Enforce a dependency to ensure individual tests are invalidated when their
# corresponding source files are too. Gotta tell make not to delete them as
# temporary intermediate files as well, since this pattern would cause that
# assumption.
%_test.go: %.go
	touch ${@}
.NOTINTERMEDIATE: %_test.go

################################################################################
#: Run examples
examples: \
	${TF_LOG_FILES}
.PHONY: examples

TF_INPUT         = 0
TF_IN_AUTOMATION = true
TF_LOG           = json
TF_PROVIDER_LOG  = json
TF_LOG_PATH      = terraform.jsonl

# <editor-fold desc="This isn't working for some reason">

#define TPL_TF_DEPS =
#${1}/terraform.jsonl: \
#		${1}/${2}.tf \
#		${1}/${2}.tftest.hcl \
#		${1}/local_override.tf \
#		${1}/main.tf \
#		${EXE_OUT}
#	@rm -f ${@}
#	@export TF_LOG_PATH=${@}
#	@echo "Testing $(dir ${@})"
#	${CMD_TERRAFORM} -chdir=$(dir ${@}) test
#${1}/${2}.tf: \
#		internal/provider/$(patsubst idmc_%,%,$(notdir ${dir}))${3}.go
#	touch $${@}
#endef

#$(eval $(call ${TPL_TF_DEPS},examples/provider/,provider))
#$(foreach dir,${TF_SRC_DIRS_DAT},$(eval $(call ${TPL_TF_DEPS},${dir},data-source,_data)))
#$(foreach dir,${TF_SRC_DIRS_RES},$(eval $(call ${TPL_TF_DEPS},${dir},resource)))
#$(foreach dir,${TF_SRC_DIRS_FUN},$(eval $(call ${TPL_TF_DEPS},${dir},function)))

# </editor-fold>

.PRECIOUS: ${TF_LOG_FILES}
${TF_LOG_FILES}: %/terraform.jsonl: \
		$(wildcard %/*.tftest.hcl) \
		$(wildcard %/*.tf) \
		%/local_override.tf \
		${EXE_OUT}
	@rm -f ${@}
	@export TF_LOG_PATH=${@}
	@echo -e "\nTesting $(dir ${@})"
	${CMD_TERRAFORM} -chdir=$(dir ${@}) test

# This ensures that any local auth settings can be immediately updated in each
# of the examples. Not worth extracting this pattern into a re-usable variable,
# however, since it's already covered elsewhere.
$(addsuffix local_override.tf,${TF_SRC_DIRS}): %/local_override.tf: \
		examples/local_override.tf
	cp ${<} ${@}

################################################################################
#: Generate documentation.
docs: ${DOC_OUT_FILES}
.PHONY: docs

${DOC_OUT_FILES} &: \
		${GO_SRC_FILES} \
		${TF_SRC_FILES} \
		${DOC_OUT_TPLS} \
		| .build/
	${CMD_TFPLUGINDOCS} \
		generate \
		--provider-name idmc \
		--rendered-provider-name IDMC \
		--website-temp-dir .build/tfplugindocs
