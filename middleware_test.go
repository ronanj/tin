package tin

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidMethod(t *testing.T) {

	tin := New()
	tin.Use(CORSMiddleware())

	called := false
	tin.GET("/path", func(c *Context) {
		c.JSON(200, H{"status": "ok"})
		called = true
	})

	r := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	tin.router.ServeHTTP(w, r)

	if !called {
		t.Error("Handler should be called")
	}

	res := w.Result()
	if res.Status != "200 OK" {
		t.Errorf("Expect context to be aborted: %s!=200", res.Status)
	}
	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS not set")
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error("Could not ready the body")
	}
	if string(body) != `{"status":"ok"}` {
		t.Error("Invalid response", string(body))
	}
}

func TestOptions(t *testing.T) {

	tin := New()
	tin.Use(CORSMiddleware())

	called := false
	tin.GET("/path", func(c *Context) {
		c.JSON(200, H{"status": "ok"})
		called = true
	})

	r := httptest.NewRequest("OPTIONS", "/path", nil)
	w := httptest.NewRecorder()
	tin.router.ServeHTTP(w, r)

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

func TestInvalidMethod(t *testing.T) {

	tin := New()
	tin.Use(CORSMiddleware())

	called := false
	tin.GET("/path", func(c *Context) {
		c.JSON(200, H{"status": "ok"})
		called = true
	})

	r := httptest.NewRequest("POST", "/path", nil)
	w := httptest.NewRecorder()
	tin.router.ServeHTTP(w, r)

	if called {
		t.Fatal("Handler should not be called")
	}

	res := w.Result()
	if res.Status != "405 Method Not Allowed" {
		t.Fatalf("Expect context to be aborted: %s!=405", res.Status)
	}

	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("CORS not set")
	}

}

func TestMethodReturns404(t *testing.T) {

	tin := New()
	tin.Use(CORSMiddleware())

	tin.GET("/path", func(c *Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.JSON(404, H{"status": "error"})
	})

	r := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	tin.router.ServeHTTP(w, r)

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

func TestInvalidRoute(t *testing.T) {

	tin := New()
	tin.Use(CORSMiddleware())

	tin.GET("/path", func(c *Context) {
	})

	r := httptest.NewRequest("GET", "/nope", nil)
	w := httptest.NewRecorder()
	tin.router.ServeHTTP(w, r)

	res := w.Result()
	if res.Status != "404 "+http.StatusText(http.StatusNotFound) {
		t.Fatalf("Expect context to be aborted: %s!=404", res.Status)
	}
	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("CORS should be set even for invalid path")
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}
	if !strings.HasPrefix(string(body), `404`) {
		t.Fatalf("Invalid response '%s'", string(body))
	}

}
