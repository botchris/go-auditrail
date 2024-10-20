# Go Auditrail

[![go test](https://github.com/botchris/go-auditrail/actions/workflows/go-test.yml/badge.svg)](https://github.com/botchris/go-auditrail/actions/workflows/go-test.yml)
[![golangci-lint](https://github.com/botchris/go-auditrail/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/botchris/go-auditrail/actions/workflows/golangci-lint.yml)

Provides a simple interface for logging audit events in Go with support for multiple backends such as
file, Elasticsearch, and more.

## Installation

```bash
go get github.com/auditrail/go-auditrail
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/botchris/go-auditrail"
)

func main() {
    // Create a new client
    client := auditrail.NewFileLogger(os.Stdout)

    // Log a message
    entry := auditrail.NewEntry("john", "order.deleted", "ordersService")

    if err := client.Log(context.TODO(), entry); err != nil {
        fmt.Println(err.Error())
    }

}
```
