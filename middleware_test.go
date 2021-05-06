package tin

import (
	"io/ioutil"
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
		c.JSON(200, H{"status": "ok"})
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

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}
	if string(body) != `{"status":"ok"}` {
		t.Fatal("Invalid response", string(body))
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

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}
	if string(body) != `` {
		t.Fatal("Invalid response", string(body))
	}

}

func Test404(t *testing.T) {

	tin := New()
	tin.Use(CORSMiddleware())

	req := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	ctx := tin.newContext(w, req, nil)

	tin.handle(func(c *Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.JSON(404, H{"status": "error"})
	}, ctx)

	res := w.Result()
	if res.Status != "404 Not Found" {
		t.Fatalf("Expect context to be aborted: %s!=404", res.Status)
	}
	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("CORS not set")
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}
	if string(body) != `{"status":"error"}` {
		t.Fatal("Invalid response", string(body))
	}

}
