# cfsvcenv &middot; [![Travis-CI](https://travis-ci.org/govau/cfsvcenv.svg)](https://travis-ci.org/govau/cfsvcenv) [![GoDoc](https://godoc.org/github.com/govau/cfsvcenv?status.svg)](http://godoc.org/github.com/govau/cfsvcenv) [![Report card](https://goreportcard.com/badge/github.com/govau/cfsvcenv)](https://goreportcard.com/report/github.com/govau/cfsvcenv)

`cfsvcenv` is a Go package for binding Cloud Foundry service credentials
(including user-provided services) to environment variables.

## Rationale

Your Go program uses Cloud Foundry services to provide configuration to it.
Cloud Foundry does this by injecting the configuration into the `VCAP_SERVICES`
environment variable.

However, you want to keep your program CF-agnostic, and you also want to keep
development simple

`cfsvcenv` binds the CF service credentials into the OS's environment so that
your Go program does not need any special treatment when used within CF.

## Usage and examples

* [Read the documentation](https://godoc.org/github.com/govau/cfsvcenv)
* [Example](https://godoc.org/github.com/govau/cfsvcenv/tree/master/example)

## Caveats

You cannot define environment variables in the global scope because these have
not yet been bound. This is because the `init()` function that does the binding
has not yet run.

For example this won't work:

```go
package main

import (
	"fmt"
	"os"
)

var serviceAPIKey = os.Getenv("SERVICE_API_KEY")

func main() {
	fmt.Println(serviceAPIKey) // Will be empty
}
```

## Development

```sh
go test -race ./...
```
