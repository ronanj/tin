
# TIN: Tiny Gin

[![CircleCI](https://circleci.com/gh/ronanj/tin.svg?style=svg)](https://app.circleci.com/pipelines/github/ronanj/tin)


A dependency free and light-weight drop-in replacement for Gin framework with essential features

## Usage

Just need to import with 

``` go
import (
	gin "github.com/ronanj/tin"
)
```

## Built-in extensions

### CORS

```
import (
    gin "github.com/ronanj/tin"
)

r := gin.New()
gin.Use(gin.CORSMiddleware())
```


## Known limitations

### Path Matching

The router only supports termnial matching. For instance this will work:

```go
router.GET("/aaa/bbb/:server", func(c *tin.Context) {
	server := c.Param("server")
}

router.GET("/aaa/bbb/:server/:name", func(c *tin.Context) {
	server := c.Param("server")
	name := c.Param("name")
}
```

But this will panic:

```go
router.GET("/aaa/bbb/:server/ccc", func(c *tin.Context) {
	server := c.Param("server")
}
```
