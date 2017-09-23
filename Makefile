SHELL := /bin/bash

VERSION=$(shell cat constants/version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)
DIRS_WITHOUT_VENDOR=$(shell ls -d */ | grep -vE "vendor")
PKGS_WITHOUT_VENDOR=$(shell go list ./... | grep -v "/vendor/")

.PHONY: help
help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  all           to check, build, test and release snips"
	@echo "  check         to vet and lint snips"
	@echo "  build         to create bin directory and build snips"
	@echo "  test          to run test"
	@echo "  test-coverage to test with coverage"
	@echo "  install       to install snips to ${GOPATH}/bin"
	@echo "  uninstall     to uninstall snips"
	@echo "  release       to build and release snips"
	@echo "  clean         to clean build and test files"

.PHONY: all
all: check build release clean unit-test unit-coverage

.PHONY: check
check: vet lint

.PHONY: vet
vet:
	@echo "go tool vet, on snips packages"
	@go tool vet -all ${DIRS_WITHOUT_VENDOR}
	@echo "ok"

.PHONY: lint
lint:
	@echo "golint, on snips packages"
	@lint=$$(for pkg in ${PKGS_WITHOUT_VENDOR}; do golint $${pkg}; done); \
	 if [[ -n $${lint} ]]; then echo "$${lint}"; exit 1; fi
	@echo "ok"

.PHONY: build
build:
	@echo "build snips"
	mkdir -p ./bin
	go build -o ./bin/snips .
	@echo "ok"

.PHONY: test
test:
	@echo "run test"
	go test -v ${PKGS_WITHOUT_VENDOR}
	@echo "ok"

.PHONY: test-coverage
test-coverage:
	@echo "run test with coverage"
	for pkg in ${PKGS_WITHOUT_VENDOR}; do \
		output="coverage$${pkg#github.com/yunify/snips}"; \
		mkdir -p $${output}; \
		go test -v -cover -coverprofile="$${output}/profile.out" $${pkg}; \
		if [[ -e "$${output}/profile.out" ]]; then \
			go tool cover -html="$${output}/profile.out" -o "$${output}/profile.html"; \
		fi; \
	done
	@echo "ok"

.PHONY: install
install: build
	@if [[ -z "${GOPATH}" ]]; then echo "ERROR: $GOPATH not found."; exit 1; fi
	@echo "Installing into ${GOPATH}/bin/snips..."
	@cp ./bin/snips ${GOPATH}/bin/snips
	@echo "ok"

.PHONY: uninstall
uninstall:
	@if [[ -z "${GOPATH}" ]]; then echo "ERROR: $GOPATH not found."; exit 1; fi
	@echo "Uninstalling snips..."
	rm -f ${GOPATH}/bin/snips
	@echo "ok"

.PHONY: release
release:
	@echo "release snips"
	mkdir -p ./release
	@echo "for Linux"
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux/snips .
	mkdir -p ./release
	tar -C ./bin/linux/ -czf ./release/snips-v${VERSION}-linux_amd64.tar.gz snips
	@echo "for macOS"
	mkdir -p ./bin/linux
	GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/snips .
	tar -C ./bin/darwin/ -czf ./release/snips-v${VERSION}-darwin_amd64.tar.gz snips
	@echo "for Windows"
	mkdir -p ./bin/windows
	GOOS=windows GOARCH=amd64 go build -o ./bin/windows/snips.exe .
	cd ./bin/windows/ && zip ../../release/snips-v${VERSION}-windows_amd64.zip snips.exe
	@echo "ok"

.PHONY: clean
clean:
	rm -rf ./bin
	rm -rf ./release
	rm -rf ./coverage
