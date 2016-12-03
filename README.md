# Snips

<span style="display: inline-block">
[![Build Status](https://travis-ci.org/yunify/snips.svg?branch=master)]
(https://travis-ci.org/yunify/snips)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)]
(https://github.com/yunify/snips/blob/master/LICENSE)
</span>

A code generator for QingCloud & QingStor SDKs.

## Introduction

Snips generates various code for QingCloud and QingStor SDKs with API
specifications in OpenAPI Specification (Swagger) v2.0 format.

Run `snips --help` to get to help messages of snips.

``` bash
$ snips --help
A code generator for QingCloud & QingStor SDKs.
It is used to generate code from our public APIs currently.

For example:
  $ snips -m QingStor -n latest \
          -s ./specs -t ./templates/qingstor/go \
          -o ./publish/qingstor-sdk-go/service
  $ ...
  $ snips --service=QingStor \
          --service-api-version=latest \
          --spec=./specs \
          --template=./templates/qingstor/ruby \
          --output=./publish/qingstor-sdk-ruby/lib/qingstor/sdk/service
  $ ...

Copyright (C) 2016 Yunify, Inc.

Usage:
  snips [flags]

Flags:
  -o, --output string                Output files directory.
  -i, --service-api-version string   Service API version to use. (default "latest")
  -m, --service string               Service to use.
  -s, --spec string                  Spec files directory.
      --spec-format string           Format of spec file. (default "Swagger-v2.0")
  -t, --template string              Template files directory.
  -v, --version                      Show version.
```

## Installation

Snips is a command line tool, and it's easy to have it installed. You can build
it from source code or download the precompiled binary directly.

### Install from Source Code

Snips requires Go 1.6 or later's vendor feature, the dependencies the project
used are included in the `vendor` directory. And we use [govendor]
[govender link] to manage project dependence.

``` bash
go get -u github.com/yunify/snips
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
+------------------------------------------------+--------------------+
|                                                |  Workflow Diagram  |
|       API Specifications                       +--------------------|
|                   +                                                 |
|                   |                               Scenario          |
|     Templates     |               +------------->  Based            |
|         +         |      +------>SDKs             Testing           |
|         |         |      |        ^                  +              |
|         |         v      |        |                  |              |
|         +-----> Snips+---+        |                  |              |
|                                   |                  |              |
|                                   |                  v              |
|       Handcraft+------------------+               Publish           |
|                                                                     |
+---------------------------------------------------------------------+
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

There's a directory that refer to the QingStor API specifications
(`./specs/qingstor`), and a directory (`./templates`) that consisted of
templates files which will be used to generate API code, and a directory
(`./test/features`) refer to the test scenario definitions.

Every time the QingStor API changes, just update the `specs/qingstor` submodule
with `make update` and run regenerate code with `make generate`. Then add/change
test for this API change, and the final step is to run the online test with
`make test` to make sure the SDK is working properly.

## Contributing

Please see [_`Contributing Guidelines`_](./CONTRIBUTING.md) of this project
before submitting patches.

## LICENSE

The Apache License (Version 2.0, January 2004).

[govender link]: https://github.com/kardianos/govendor
[qingstor-sdk-go link]: https://github.com/yunify/qingstor-sdk-go

[release link]: https://github.com/yunify/snips/releases
[example download link]: https://github.com/yunify/snips/releases/download/v0.0.7/snips-v0.0.7-darwin_amd64.tar.gz
