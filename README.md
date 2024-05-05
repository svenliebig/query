# query

a simple go package to create a URL query string from a struct, using [struct tags](https://pkg.go.dev/reflect#StructTag) to determine the key name and whether to include the field in the query string.

## install

```sh
go get github.com/svenliebig/query
```

## usage

You can use the `query` struct tag to specify the key name for the field in the query string. If the tag is not present, the field name will be used.

The function `query.Stringify` will return a URL query string from the struct, missing tags and unsupported kings will be ignored with that function.

If you want to react to missing tags or unsupported types, you can use the `query.StringifyE` function, which will report errors.

```go
package main

import (
    "fmt"
    "github.com/svenliebig/query"
)

type User struct {
    Name     string `query:"name"`
    Age      int    `query:"age"`
    Verified bool   `query:"verified"`
    Missing  string
    Complex  complex64
}

func main() {
    u := User{
        Name: "Alice",
        Age:  30,
        Verified: true,
    }

    q := query.Stringify(u)

    fmt.Println(q) // "name=Alice&age=30&verified=true"

    q, err := query.StringifyE(u)

    if err != nil {
        fmt.Println(err) // missing tag for field 'Missing'
    }

    fmt.Println(q) // ""
}
```

## options

There are some options you can use to customize the behavior of the `query.Stringify` and `query.StringifyE` functions.

### `query.SkipEmpty`

If you want to skip empty values in the query string, you can set the `query.SkipEmpty` option to `true`.

```go
package main

import (
    "fmt"
    "github.com/svenliebig/query"
)

type User struct {
    Name     string `query:"name"`
    Age      int    `query:"age"`
    Verified bool   `query:"verified"`
}

func main() {
    u := User{
        Name: "",
        Age:  0,
        Verified: false,
    }

    q := query.Stringify(u)

    fmt.Println(q) // "name=&age=0&verified=false"

    q = query.Stringify(u, query.SkipEmpty{})

    fmt.Println(q) // "verified=false"
}
```

The option will skip empty strings & zero values, but not `false` for `bool` fields.
