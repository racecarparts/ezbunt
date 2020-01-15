# ezbunt

ezbunt is a Go module for conveniently interacting with a [buntdb](https://github.com/tidwall/buntdb).  Functions include: creating a data file, storing simple data types, and complex types as JSON objects.

## Documentation

[godoc.org/github.com/racecarparts/ezbunt](https://godoc.org/github.com/racecarparts/ezbunt)

## Usage

via: `go.mod` (go module):

```text
require (
    ...
    github.com/racecarparts/ezbunt v0.1.2
)
```


or via: `go get`

```bash
$ go get -u github.com/racecarparts/ezbunt
```

## Example

```go
package main

import (
    "fmt"
    ez "github.com/racecarparts/ezbunt"
)

func main() {
    ez := ez.New("data.file")
    ez.WriteKeyVal("my:1234", "sharona")

    pairs, err := ez.GetPairs("my")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(pairs)

    val, err := ez.GetVal("my:1234")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(val)
}
```
