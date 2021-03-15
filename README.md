# Kubernetes API documentation generator

`k8s-api-docgen` is a tool that reads Go source files to extract metadata
about Kubernetes API custom resources defined by the developer in [Godoc](https://blog.golang.org/godoc)
and to produce different kinds of output - such as JSON or Markdown.

Its main use case is to automatically generate documentation in Markdown
format to be used within a CI/CD pipeline of a software written in Go
for Kubernetes, such as an operator.

## Requirements

To compile this software you need:

- a working [Go compiler](https://golang.org/);

- a `make` installation (the one in MacOSX and the GNU ones are known to be
  working correctly).

The code has been tested on MacOS X and Linux, but should resonably work in other
platforms too.

## How to compile the code

Just use the Makefile:

    $ make
    [...]

## Usage instructions

The tool needs to have access to the Go source files for your Kubernetes API.
These are the files that will also be used by
[controller-gen](https://book.kubebuilder.io/reference/controller-gen.html) to
generate the CRD.

Supposing these files are reachable at `../operator/api/v1` you can extract the
documentation in JSON format via:

    $ ./bin/k8s-api-docgen ../operator/api/v1/*types.go

The JSON stream will be written to *standard output*. Should you desire to
create a file, you can use the `-o` option as follows:

    $ ./bin/k8s-api-docgen -o documentation.json ../operator/api/v1/*types.go

## Copyright

`k8s-api-docgen` is distributed under Apache License 2.0.
Copyright (C) 2021 EnterpriseDB Corporation.
Some portions Copyright 2016 The prometheus-operator Authors.
