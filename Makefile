.PHONY: mocks
mocks:
	mockery --config=./config/.mockery.yml

.PHONY: lint
lint:
	golangci-lint run --config=./config/.golangci.yml

.PHONY: generate
generate:
	@go generate ./...

UNIT_TEST_PACKAGES := $(shell go list ./... | grep -v /test)
.PHONY: unit-test
unit-test:
	@go test -race $(UNIT_TEST_PACKAGES)

.PHONY: unit-test-cov
unit-test-cov:
ifeq (,$(UNIT_COV_FILE))
	@echo "\033[0;31mUNIT_COV_FILE not found. Include it and retry. \033[0m"
	@exit 1
else
	@go test -race -coverprofile=$(UNIT_COV_FILE) -covermode=atomic $(UNIT_TEST_PACKAGES) -coverpkg=./...
endif

# TODO: fill it.
# .PHONY: integration-test
# integration-test:
#