# **dj**: For cleaner and safer Go code

<a href="https://github.com/go-dj/dj/actions/workflows/check.yml"><img src="https://github.com/ProtonMail/go-proton-api/actions/workflows/check.yml/badge.svg?branch=master" alt="CI Status"></a>
<a href="https://pkg.go.dev/github.com/go-dj/dj"><img src="https://pkg.go.dev/badge/github.com/ProtonMail/go-proton-api" alt="GoDoc"></a>
<a href="https://goreportcard.com/report/github.com/go-dj/dj"><img src="https://goreportcard.com/badge/github.com/ProtonMail/go-proton-api" alt="Go Report Card"></a>
<a href="LICENSE"><img src="https://img.shields.io/github/license/go-dj/dj.svg" alt="License"></a>

`dj` is a Go 1.18+ library that makes it easy to write clean and safe Go code.

# Install
Simply `go get` the library:

```
go get github.com/go-dj/dj
```

# Usage
Import `dj` and use it like so:

```go
import (
    "fmt"

    "github.com/go-dj/dj"
)

func main() {
    vals := dj.RangeN(5)

    fmt.Println(vals) // [0 1 2 3 4]

    doubled := dj.MapEach(vals, func(val int) int {
        return val * 2
    })

    fmt.Println(doubled) // [0 2 4 6 8]
}
```