# Sim

[![Go](https://github.com/kitproj/sim/actions/workflows/go.yml/badge.svg)](https://github.com/kitproj/sim/actions/workflows/go.yml)
[![goreleaser](https://github.com/kitproj/sim/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/kitproj/sim/actions/workflows/goreleaser.yml)

## Why

Make the dev loop crazy fast.

## What

Sim is straight-forward API simulation tool that's tiny, fast, secure and scalable.

Most of today's API mocking tools run in virtual machines such as the JVM or NPM. Sim is a single binary with zero dependencies.

- It's orders of magnitude smaller binary and memory usage. Which much lower CPU usage. Each process can simulation multiple APIs. 
- Running on Kubernetes? Three pods could simulate every API in your organization with high-availability.
- 
Sim doesn't just mock APIs, it allows you to specify scripts for each API operation which have access to a simple
key-value database that allows APIs to save state between requests.

Sim was written with extensive help from AI.

## Install

Like `jq`, `sim` is a tiny (16Mb) standalone binary. You can download it from
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

As Docker image:

```bash
docker run --rm -ti -v my-directory:/apis ghcr.io/kitproj/sim /sim /apis
```

## Usage

Create a directory containing files named `*.yaml`. 

Then run:

```bash
sim my-directory

```

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

Scripting

```yaml
openapi: 3.0.0
info:
  title: Teapot API
  version: 1.0.0
servers:
  - url: http://localhost:4040
# This script is executed when the spec is loaded.
# You can use it to specify global variables and functions that are available to all scripts.
x-sim-script: |
  var status = 418
paths:
  /teapot:
    get:
      # This script is executed whenever the request is serviced.
      # The last variable is the response object.
      x-sim-script: |
        response = {
           "status": status,
           "headers": {
             "Teapot": "true"
           },
           "body": { "message": "I'm a teapot" }
         }
      responses:
        '200':
          description: OK
```

Scripting with a database:

```yaml
openapi: 3.0.0
info:
  title: Document API
  version: 1.0.0
servers:
  - url: http://localhost:4040
paths:
  /documents:
    post:
      x-sim-script: |
        var uuid = randomUUID();
        db.put("/documents/" + uuid, request.body)
        response = {
          "status": 201,
            "headers": {
                "Location": "/documents/" + uuid
            }
        }
      responses:
        '200':
          description: OK
    get:
      x-sim-script: |
        response = {
          "body": db.list("/documents")
        }
      responses:
        '200':
          description: OK
  /documents/{id}:
    get:
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      x-sim-script: |
        var document = db.get("/documents/" + request.pathParams.id);
        if (document) {
          response = {
            "body": document
          }
        } else {
          response = {
              "status": 404,
          }
        }
      responses:
        '200':
          description: OK

    delete:
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      x-sim-script: |
        db.delete("/documents/" + request.pathParams.id)
        response = {}
      responses:
        '204':
          description: OK
```

## Reference

In you script you have access to the following:

### `request`

An object containing the HTTP request, e.g.

```json
{
  "method": "PUT",
  "path": "/documents/bar",
  "pathParams": {
    "id": "bar"
  },
  "queryParams": {
    "foo": "bar"
  },
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "baz": "qux"
  }
}
```

For example:

```javascript
var body = request.body;
```


### `randomUuid()`

A function that generates a random UUID. For example:


```javascript
var uuid = randomUUID();
```

### `db`

A service that allows you to persist and access data:

```javascript
// get an object, maybe null
var value = db.get(key);
// put an object (idempotent)
var existed = db.put(key, value);
// delete an object (idempotent)
var deleted = db.delete(key)
// return an array of all objects
var list = db.list(keyPrefix);
```

For example:

```javascript
var obj = db.get("/my-api/" + id)
if (obj) {
    response = {status: 404}
} else {
    response = {body: obj}
}
```