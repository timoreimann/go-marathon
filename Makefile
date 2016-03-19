#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-02-10 15:35:14 +0000 (Tue, 10 Feb 2015)
#
#  vim:ts=2:sw=2:et
#
HARDWARE=$(shell uname -m)
VERSION=$(shell awk '/const Version/ { print $$4 }' version.go | sed 's/"//g')
DEPS=$(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
PACKAGES=$(shell go list ./...)
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr

.PHONY: test examples authors changelog

build:
	go build

authors:
	git log --format='%aN <%aE>' | sort -u > AUTHORS

deps:
	@echo "--> Installing build dependencies"
	@go get -d -v ./... $(DEPS)

lint:
	@echo "--> Running golint"
	@which golint 2>/dev/null ; if [ $$? -eq 1 ]; then \
		go get -u github.com/golang/lint/golint; \
	fi
	@golint .

vet:
	@echo "--> Running go vet ./..."
	@go vet ./...

cover:
	@echo "--> Running go test --cover"
	@go test --cover

format:
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

test: deps
	@echo "--> Running go tests"
	@go test -v
	@$(MAKE) vet
	@$(MAKE) cover

changelog: release
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B > changelog
