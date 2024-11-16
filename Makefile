LOCAL_BIN:=$(CURDIR)/bin

PHONY: test
test:
	go test -race -count 1 ./...

PHONY: generate
generate: .genmock

PHONY: .genmock
.genmock: .bin-mock
	$(info $(shell printf "\033[34;1m▶\033[0m") go generate-mocks...)
	@for f in $(shell find internal -name 'genmock.go'| sort -u); do \
		PATH=$(LOCAL_BIN):"$(PATH)" go generate $$f; \
	done

PHONY: .bin-mock
.bin-mock:
	$(info $(shell printf "\033[34;1m▶\033[0m") Installing mockery...)
	GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v2@v2.47.0