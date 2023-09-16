package tin

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		ctx.activateRecovery = true
	}
}

func RecoveryWithNotification(notifier func(ctx *Context, e interface{}) bool) HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			ctx.activateRecovery = true
			ctx.recoveryNotifier = notifier
		}()
	}
}

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

func (t *Tin) Use(middleware HandlerFunc) {
	t.middlewares = append(t.middlewares, middleware)
}

/* Allow only one middleware to start with */
func (C *Context) Next() {

}

func (c *Context) AbortWithStatus(status int) {
	c.Writer.WriteHeader(status)
	c.isAborted = true
}

func (t *Tin) applyMiddleware(ctx *Context) bool {

	for _, mw := range t.middlewares {
		mw(ctx)
		if ctx.isAborted {
			return false
		}
	}
	return true
}

func (t *Tin) handle(handle func(c *Context), ctx *Context) {

	handle(ctx)

}
