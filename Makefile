SHELL := /bin/bash

.PHONY: all check vet lint build test coverage install uninstall release clean

VERSION=$(shell cat metadata/version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)
DIRS_WITHOUT_VENDOR=$(shell ls -d */ | grep -vE "vendor")
PKGS_WITHOUT_VENDOR=$(shell go list ./... | grep -v "/vendor/")

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  all           to check, build, test and release snips"
	@echo "  check         to vet and lint snips"
	@echo "  build         to create bin directory and build snips"
	@echo "  unit-test     to run test"
	@echo "  unit-coverage to test with coverage"
	@echo "  install       to install snips to /usr/local/bin/snips"
	@echo "  uninstall     to uninstall snips"
	@echo "  release       to build and release snips"
	@echo "  clean         to clean build and test files"

all: check build release clean unit-test unit-coverage

check: vet lint

vet:
	@echo "go tool vet, on snips packages"
	@go tool vet -all ${DIRS_WITHOUT_VENDOR}
	@echo "ok"

lint:
	@echo "golint, on snips packages"
	@lint=$$(for pkg in ${PKGS_WITHOUT_VENDOR}; do golint $${pkg}; done); \
	 if [[ -n $${lint} ]]; then echo "$${lint}"; exit 1; fi
	@echo "ok"

build:
	@echo "build snips"
	mkdir -p ./bin
	go build -o ./bin/snips .
	@echo "ok"

unit-test:
	@echo "run test"
	go test -v ${PKGS_WITHOUT_VENDOR}
	@echo "ok"

unit-coverage:
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

install: build
	@echo "install snips to /usr/local/bin/snips"
	cp ./bin/snips /usr/local/bin/snips
	@echo "ok"

uninstall:
	@echo "delete /usr/local/bin/snips"
	rm -f /usr/local/bin/snips
	@echo "ok"

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
	tar -C ./bin/windows/ -czf ./release/snips-v${VERSION}-windows_amd64.tar.gz snips.exe
	@echo "ok"

clean:
	rm -rf ./bin
	rm -rf ./release
	rm -rf ./coverage
