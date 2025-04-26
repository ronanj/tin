
# TIN: Tiny Gin

[![Go Audit](https://github.com/ronanj/tin/actions/workflows/audit.yml/badge.svg)](https://github.com/ronanj/tin/actions/workflows/audit.yml)

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

```go
import (
    gin "github.com/ronanj/tin"
)

r := gin.New()
gin.Use(gin.CORSMiddleware())
```

### SSE (Server Side Events)

```go
import (
    gin "github.com/ronanj/tin"
    "time"
)

r := gin.New()

r.GET("/stream", func(c *gin.Context) {
    sse := c.SSE()
    sse.Event("init")
    for i:=1; i<=n; i++ {
        if err := sse.JSON(gin.H{"counter": i}); err!=nil {
            // The client is gone, so just abort?
            break
        }
        time.Sleep(time.Second)
    }
    sse.Event("done")
})
```


## Known limitations

### Path Matching

The router support complex matching, but the matches should be added by order of priority. Eg, in the example below, `/aaa/bbb/:server/ccc` must be added before `/aaa/bbb/:server/:name`.

```go
router.GET("/aaa/bbb/:server/ccc", func(c *tin.Context) {
    server := c.Param("server")
}

router.GET("/aaa/bbb/:server/:name", func(c *tin.Context) {
    server := c.Param("server")
    name := c.Param("name")
}

router.GET("/aaa/bbb/:server", func(c *tin.Context) {
	server := c.Param("server")
}
```
