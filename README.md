# Snips

[![Build Status](https://travis-ci.org/yunify/snips.svg?branch=master)](https://travis-ci.org/yunify/snips)
[![Go Report Card](https://goreportcard.com/badge/github.com/yunify/snips)](https://goreportcard.com/report/github.com/yunify/snips)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/yunify/snips/blob/master/LICENSE)

A code generator for RESTful APIs.

## Introduction

Snips generates various code using API specifications in OpenAPI Specification (Swagger) v2.0 format.

Run `snips --help` to get to help messages of snips.

``` bash
$ snips --help
A code generator for RESTful APIs.

For example:
  $ snips -f ./specs/qingstor/api.json
          -t ./templates/qingstor/go \
          -o ./publish/qingstor-sdk-go/service
  $ ...
  $ snips --file=./specs/qingstor/api.json \
          --template=./templates/qingstor/ruby \
          --output=./publish/qingstor-sdk-ruby/lib/qingstor/sdk/service
  $ ...

Copyright (C) 2016 Yunify, Inc.

Usage:
  snips [flags]

Flags:
  -f, --file string       Specify the spec file.
      --format string     Specify the format of spec file. (default "OpenAPI-v2.0")
  -o, --output string     Specify the output directory.
  -t, --template string   Specify template files directory.
  -v, --version           Show version.
```

## Installation

Snips is a command line tool, and it's easy to have it installed. You can build
it from source code or download the precompiled binary directly.

### Install from Source Code

Snips requires Go 1.6 or later's vendor feature, the dependencies the project
used are included in the `vendor` directory. And we use [glide][glide link]
to manage project dependence.

``` bash
$ git clone git@github.com:yunify/snips.git
$ glide install
$ make install
```

___Notice:___ _You can also use Go 1.5 with the `GO15VENDOREXPERIMENT=1`._

### Download Precompiled Binary

0. Go to [releases tab][release link] and download the binary for your operating
system, for example [`snips-v0.0.7-darwin_amd64.tar.gz`][example download link].
0. Unarchive the downloaded file, and put the executable file `snips` into a
directory that in the `$PATH` environment variable, for example
`/usr/local/bin`.
0. Run `snips --help` to get started.

## SDK Development Workflow

Snips takes API specifications and template to generate lots of code for APIs,
then these generated code plus the handcrafted code makes up the SDK. Next,
we use scenario based test to make sure our SDKs are working properly, and
ensure their functional consistency.

```
+---------------------------------------------+--------------------+
|                                             |  Workflow Diagram  |
|       API Specifications                    +--------------------|
|                   +                                              |
|                   |                               Scenario       |
|     Templates     |               +------------->  Based         |
|         +         |      +------>SDKs             Testing        |
|         |         |      |        ^                  +           |
|         |         v      |        |                  |           |
|         +-----> Snips+---+        |                  |           |
|                                   |                  |           |
|                                   |                  v           |
|       Handcraft+------------------+               Publish        |
|                                                                  |
+------------------------------------------------------------------+
```

#### Add an SDK for Another Programing Language

0. Create handcraft files of SDK, to handle configuration, network request and etc.
0. Writing templates for API code.
0. Generate code using snips.
0. Running tests.
0. Publish.

#### Update an Exists SDK

0. Change handcraft files (if needed) and update the API specs.
0. Regenerate code.
0. Running tests.
0. Publish.

### Example

Let's take Go SDK for QingStor ([`qingstor-sdk-go`][qingstor-sdk-go link]) for
example.

#### Prepare

- `./specs/qingstor`: Refer to [the QingStor API specifications][api specs link]
- `./test/features`: Refer to [the QingStor SDK test scenarios][sdk test scenarios link]

___Tips:___ _Include these files as git submodule._

#### Procedures

0. Create template files which will be used to generate API code in `./template`.
0. Generate code using snips, and format the generated code.

    ``` bash
    $ snips --version
    snips version 0.0.9
    $ snips --service=qingstor \
            --service-api-version=latest \
            --spec="./specs" \
            --template="./template" \
            --output="./service"
    Loaded templates from ./template
    4 template(s) detected.
    Loaded service QingStor (2016-01-06) from ./specs

    Generating to: ./service/qingstor.go
    Generating to: ./service/object.go
    Generating to: ./service/bucket.go
    Generating to: ./service/types.go

    Everything looks fine.
    $ gofmt -w .
    ```

0. Implement test scenarios in `./test`.

    ``` bash
    $ ls ./test
    bucket.go                 config.yaml.example       test_config.yaml
    bucket_acl.go             features                  test_config.yaml.example
    bucket_cors.go            main.go                   utils.go
    bucket_external_mirror.go object.go                 vendor
    bucket_policy.go          object_multipart.go
    config.yaml               service.go
    ```

0. Running scenarios test, and pass all tests.

    ``` bash
    $ pushd ./test
    $ go run *.go
    ...
    38 scenarios (38 passed)
    84 steps (84 passed)
    1m2.408357076s
    $ popd
    ```

0. Every time the QingStor API changes, just update the `specs/qingstor` and
`./test/features` submodule and regenerate code. Then add/change test for this
API change, and rerun the online test to make sure the SDK is working properly.

## References

- [QingStor API Specs][api specs link]
- [QingStor SDK Test Scenarios][sdk test scenarios link]

## Contributing

Please see [_`Contributing Guidelines`_](./CONTRIBUTING.md) of this project
before submitting patches.

## LICENSE

The Apache License (Version 2.0, January 2004).

[glide link]: https://glide.sh
[qingstor-sdk-go link]: https://github.com/yunify/qingstor-sdk-go
[api specs link]: https://github.com/yunify/qingstor-api-specs
[sdk test scenarios link]: https://github.com/yunify/qingstor-sdk-test-scenarios

[release link]: https://github.com/yunify/snips/releases
[example download link]: https://github.com/yunify/snips/releases/download/v0.0.7/snips-v0.0.7-darwin_amd64.tar.gz
