# Go Auditrail

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

    "github.com/auditrail/go-auditrail"
)

func main() {
    // Create a new client
    client := auditrail.NewFileLogger(os.Stdout)

    // Log a message
    entry := auditrail.NewEntry("john", "order.deleted", "ordersService)

    if err := client.Log(context.TODO(), entry); err != nil {
        fmt.Println(err.Error())
    }

}
```
