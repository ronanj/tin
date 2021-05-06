package tin

import (
	"net/http/httptest"
	"testing"
)

func TestNoAbort(t *testing.T) {

	tin := New()
	tin.Use(CORSMiddleware())

	req := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	ctx := tin.newContext(w, req, nil)

	called := false
	tin.handle(func(c *Context) {
		called = true
	}, ctx)

	if !called {
		t.Fatal("Handler should be called")
	}

	res := w.Result()
	if res.Status != "200 OK" {
		t.Fatalf("Expect context to be aborted: %s!=200", res.Status)
	}
	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("CORS not set")
	}

}

func TestAbort(t *testing.T) {

	tin := New()
	tin.Use(CORSMiddleware())

	req := httptest.NewRequest("OPTIONS", "/path", nil)
	w := httptest.NewRecorder()
	ctx := tin.newContext(w, req, nil)

	called := false
	tin.handle(func(c *Context) {
		called = true
	}, ctx)

	if called {
		t.Fatal("Handler should not be called")
	}

	res := w.Result()
	if res.Status != "204 No Content" {
		t.Fatalf("Expect context to be aborted: %s!=204", res.Status)
	}

	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("CORS not set")
	}

}

func Test404(t *testing.T) {

	tin := New()
	tin.Use(CORSMiddleware())

	req := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	ctx := tin.newContext(w, req, nil)

	tin.handle(func(c *Context) {
		c.JSON(404, H{"status": "error"})
	}, ctx)

	res := w.Result()
	if res.Status != "404 Not Found" {
		t.Fatalf("Expect context to be aborted: %s!=404", res.Status)
	}
	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("CORS not set")
	}

}
