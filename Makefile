NAME := soapforce
SRCS := $(shell find . -type d -name vendor -prune -o -type f -name "*.go" -print)
VERSION := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\"" 
DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL := bin/$(NAME) 

.PHONY: test
test: glide
	@go test -cover -v `glide novendor`

.PHONY: install
install: bin/$(NAME)
	mv bin/$(NAME) /usr/local/bin/$(NAME)

.PHONY: clean
clean:
	@rm -rf bin/*
	@rm -rf vendor/*
	@rm -rf dist/*

.PHONY: format
format: import
	-@goimports -w $(SRCS)
	@gofmt -w $(SRCS)

.PHONY: import
import:
	go get golang.org/x/tools/cmd/goimports

.PHONY: cross-build
cross-build:
	-@goimports -w $(SRCS)
	@gofmt -w $(SRCS)
	@for os in darwin linux windows; do \
	    for arch in amd64 386; do \
	        GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo \
	        -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
	    done; \
	done

.PHONY: bin/$(NAME) 
bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

.PHONY: dist
dist:
	@cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) cp ../completions/zsh/_$(NAME) {} \; && \
	$(DIST_DIRS) tar zcf $(NAME)-$(VERSION)-{}.tar.gz {} \;

.PHONY: dist
docker-build:
	docker build . -t $(NAME)

.PHONY: run
run:
	go run salesforce/main.go
