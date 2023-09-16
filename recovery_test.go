package tin

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoveryHandler(t *testing.T) {

	tin := New()
	tin.Use(Recovery())

	tin.GET("/path", func(c *Context) {
		panic("this should be caught")
	})

	r := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	tin.router.ServeHTTP(w, r)

	res := w.Result()
	if res.Status != "500 "+http.StatusText(http.StatusInternalServerError) {
		t.Fatalf("Expect context to be exception: %s!=500", res.Status)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}
	if string(body) != `{"reason":"this should be caught","status":"error"}` {
		t.Fatalf("Invalid response '%s'", string(body))
	}
}

func TestRecoveryWithOutputAlreadyWritten(t *testing.T) {

	tin := New()
	tin.Use(Recovery())

	tin.GET("/path", func(c *Context) {
		c.JSON(200, "nope")
		panic("this should be caught")
	})

	r := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	tin.router.ServeHTTP(w, r)

	res := w.Result()
	if res.Status != "200 "+http.StatusText(http.StatusOK) {
		t.Fatalf("Expect context to be exception: %s!=200", res.Status)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}
	if string(body) != "\"nope\"" {
		t.Fatalf("Invalid response '%s'", string(body))
	}
}

func TestRecoveryWithAbort(t *testing.T) {

	tin := New()
	tin.Use(Recovery())

	tin.GET("/path", func(c *Context) {
		c.Writer.Header()["Access-Control-Allow-Origin"] = []string{"*"}
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("nope"))
		panic("this should be caught")
	})

	r := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	tin.router.ServeHTTP(w, r)

	res := w.Result()
	if res.Status != "400 "+http.StatusText(http.StatusBadRequest) {
		t.Fatalf("Expect context to be exception: %s!=200", res.Status)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}
	if string(body) != `{"reason":"nope","status":"error"}` {
		t.Fatalf("Invalid response '%s'", string(body))
	}
}

func TestRecoveryWithSse(t *testing.T) {

	tin := New()
	tin.Use(Recovery())

	tin.GET("/path", func(c *Context) {
		c.SSE()
		panic("this should be caught")
	})

	r := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	tin.router.ServeHTTP(w, r)

	res := w.Result()
	if res.Status != "500 "+http.StatusText(http.StatusInternalServerError) {
		t.Fatalf("Expect context to be exception: %s!=500", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}

	if string(body) != `{"reason":"this should be caught","status":"error"}` {
		t.Fatalf("Invalid response '%s'", string(body))
	}
}

func TestRecoveryWithSseOutputAlreadyWritten(t *testing.T) {

	tin := New()
	tin.Use(Recovery())

	tin.GET("/path", func(c *Context) {
		sse := c.SSE()
		sse.Send("hello", "world")
		panic("this should be caught")
	})

	r := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	tin.router.ServeHTTP(w, r)

	res := w.Result()
	if res.Status != "200 "+http.StatusText(http.StatusOK) {
		t.Fatalf("Expect context to be exception: %s!=200", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}

	expectedBody := "event: hello\ndata: \"world\"\n\n"
	if string(body) != expectedBody {
		t.Fatalf("Invalid response %s", expectedBody)
	}
}

func TestRecoveryHandlerWithNotification(t *testing.T) {

	notifierCalled := false
	tin := New()
	tin.Use(RecoveryWithNotification(func(c *Context, e interface{}) bool {
		if c.Request.URL.Path != "/path" {
			t.Fatalf("Invalid context path")

		}
		notifierCalled = true
		return true
	}))

	tin.GET("/path", func(c *Context) {
		panic("this should be caught")
	})

	r := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	tin.router.ServeHTTP(w, r)

	res := w.Result()
	if res.Status != "500 "+http.StatusText(http.StatusInternalServerError) {
		t.Fatalf("Expect context to be exception: %s!=500", res.Status)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}
	if string(body) != `{"reason":"this should be caught","status":"error"}` {
		t.Fatalf("Invalid response '%s'", string(body))
	}

	if !notifierCalled {
		t.Fatalf("Notifier has not been called")
	}
}
