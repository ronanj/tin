package tin

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {

	tin := New()
	req := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	ctx := newContext(w, req, nil)

	tin.handle(func(c *Context) {
		c.JSON(404, H{"status": "error"})
	}, ctx)

	res := w.Result()
	if res.Status != "404 Not Found" {
		t.Fatalf("Expect context to be aborted: %s!=404", res.Status)
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Fatal("Content type not set", res.Header.Get("Content-Type"))
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

func TestCustomHeader(t *testing.T) {

	tin := New()
	req := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	ctx := newContext(w, req, nil)

	tin.handle(func(c *Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.JSON(404, H{"status": "error"})
	}, ctx)

	res := w.Result()
	if res.Status != "404 Not Found" {
		t.Fatalf("Expect context to be aborted: %s!=404", res.Status)
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Fatal("Content type not set")
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
