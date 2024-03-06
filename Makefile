.PHONY: mocks
mocks:
	mockery --config=./config/.mockery.yml

UNIT_TEST_PACKAGES := $(shell go list ./... | grep -v /test)
.PHONY: unit-test
unit-test:
	@go test -race $(UNIT_TEST_PACKAGES)

# TODO: fill it.
# .PHONY: integration-test
# integration-test:
#