# Sim

[![Go](https://github.com/kitproj/sim/actions/workflows/go.yml/badge.svg)](https://github.com/kitproj/sim/actions/workflows/go.yml)
[![goreleaser](https://github.com/kitproj/sim/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/kitproj/sim/actions/workflows/goreleaser.yml)

## Why

Make the dev loop crazy fast.

## What

Sim is straight forward API simulation tools that's tiny and fast secure and scalable.

Most of today's API mocking tools run in virtual machines such as the JVM or NPM. Sim is written in Golang and leans on
standard libraries:

It's orders of magnitude smaller binary and memory usage. Which much lower CPU usage. Each process can simulation
multiple APIs. Running on Kubernetes? Three pods could simulate every API in your organization with
high-availability.

Sim doesn't just mock APIs, it allows you to specify scripts for each API operation and back it with a simple disk
storage.

Sim was written with extensive help from AI.

## Install

Like `jq`, `sim` is a tiny (8Mb) standalone binary. You can download it from
the [releases page](https://github.com/kitproj/sim/releases/latest).

If you're on MacOS, you can use `brew`:

```bash
brew tap kitproj/sim --custom-remote https://github.com/kitproj/sim
brew install sim
```

Otherwise, you can use `curl`:

```bash
curl -q https://raw.githubusercontent.com/kitproj/sim/main/install.sh | sh
```

We do not support `go install`.

## Usage

Simulations are described by their API specification. For simple mocking, specify your examples in the OpenAPI spec:

```yaml
openapi: 3.0.0
info:
  title: Hello API
  version: 1.0.0
servers:
  - url: http://localhost:8080
paths:
  /hello:
    get:
      responses:
        '200':
          description: OK
          content:
            application/json:
              example: { "message": "Hello, world!" }
```

Then run:

```bash
sim api-specs
```

```yaml
openapi: 3.0.0
info:
  title: Teapot API
  version: 1.0.0
servers:
  - url: http://localhost:4040
paths:
  /teapot:
    get:
      x-sim-script: |
        response = {
           "status": 418,
           "headers": {
             "Teapot": "true"
           },
           "body": { "message": "I'm a teapot" }
         }
      responses:
        '200':
          description: OK
```
