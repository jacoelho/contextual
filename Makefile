# disable default rules
.SUFFIXES:
MAKEFLAGS+=-r -R
DATE = $(shell date +%Y%m%d%H%M%S)

export GOBIN=$(CURDIR)/bin

default: test

.PHONY: test
test:
	go test -race -v ./...

.PHONY: tidy
tidy:
	go mod tidy && go mod verify

.PHONY: ci-tidy
ci-tidy:
	go mod tidy
	git status --porcelain go.mod go.sum || { echo "Please run 'go mod tidy'."; exit 1; }

.PHONY: lint
lint:
	docker run -it --rm -v $(CURDIR):/app -w /app golangci/golangci-lint golangci-lint run

