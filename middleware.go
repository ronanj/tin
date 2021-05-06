package tin

/*

https://stackoverflow.com/questions/29418478/go-gin-framework-cors

router = gin.New()
router.Use(CORSMiddleware())

*/

func CORSMiddleware() HandlerFunc {
	return func(c *Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

/* Allow only one middleware to start with */

func (C *Context) Next() {

}

func (c *Context) AbortWithStatus(status int) {
	c.Writer.WriteHeader(status)
	c.isAborted = true
}

func (t *Tin) handle(handle func(c *Context), ctx *Context) {

	for _, mw := range t.middlewares {
		mw(ctx)
		if ctx.isAborted {
			return
		}
	}

	handle(ctx)

}
