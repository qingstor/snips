# Snips

<span style="display: inline-block">
[![Build Status](https://travis-ci.org/yunify/snips.svg?branch=master)](https://travis-ci.org/yunify/snips)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/yunify/snips/blob/master/LICENSE)
</span>

A code generator for QingCloud & QingStor SDKs.

## Introduction

`Snips` is a command line tool that allows you to generate QingCloud & QingStor SDKs with API specifications, code templates and code base.

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

## LICENSE

The Apache License (Version 2.0, January 2004).
