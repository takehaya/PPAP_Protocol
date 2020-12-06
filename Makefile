NAME := ppap

#brunch name version
VERSION := $(shell git rev-parse --abbrev-ref HEAD)

PKG_NAME=$(shell basename `pwd`)

LDFLAGS := -ldflags="-s -w  -X \"github.com/takehaya/PPAP_Protocol/pkg/version.Version=$(VERSION)\" -extldflags \"-static\""
SRCS    := $(shell find . -type f -name '*.go')

.DEFAULT_GOAL := build
build: $(SRCS)
	go build $(LDFLAGS) -o ./bin/$(NAME) ./cmd/$(NAME)

.PHONY: crun
crun:
	go run $(LDFLAGS) ./cmd/$(NAME) --mode client --gateway "172.27.1.2"

.PHONY: sbrun
sbrun:
	go run $(LDFLAGS) ./cmd/$(NAME) --mode server --gateway "172.27.2.1" --bindaddr "172.27.1.2" --have1 "Pen!" --have2 "Apple!"

.PHONY: scrun
scrun:
	go run $(LDFLAGS) ./cmd/$(NAME) --mode server --gateway "172.27.2.2" --have1 "Pen!" --have2 "Pineapple!"


## lint
.PHONY: lint
lint:
	@for pkg in $$(go list ./...): do \
		golint --set_exit_status $$pkg || exit $$?; \
	done

.PHONY: codecheck
codecheck:
	test -z "$(gofmt -s -l . | tee /dev/stderr)"
	go vet ./...

.PHONY: clean
clean:
	rm -rf bin

.PHONY: install
install:
	go install $(LDFLAGS) ./cmd/$(NAME)
