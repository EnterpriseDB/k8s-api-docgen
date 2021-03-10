# Kubernetes API documentation generator

## Requirements

To compile this software you need:

- a working [Go compiler](https://golang.org/);

- a `make` installation (the one in MacOSX and the GNU ones are known to be
  working correctly).

The code has been tested on MacOSX and Linux but should resonably work in other
platforms too.

## How to compile the code

Just use the Makefile:

    $ make
    [...]

## Usage instructions

The tool need to have access to the Go source files for your Kubernetes API.
These are the files that will be also used by
[controller-gen](https://book.kubebuilder.io/reference/controller-gen.html) to
generate the CRD.

Supposing those files are reachable at `../operator/api/v1` you can extract the
documentation in JSON format via:

    $ ./bin/k8s-api-docgen ../operator/api/v1/*types.go

The JSON stream will be written to the standard output. Should you desire to
create a file you can:

    $ ./bin/k8s-api-docgen -o documentation.json ../operator/api/v1/*types.go
