package tin

import (
	"io/ioutil"
	"log"
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
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}
	log.Println("string(body)", string(body))
	if string(body) != `{"reason":"this should be caught","status":"error"}` {
		t.Fatalf("Invalid response '%s'", string(body))
	}
}

func TestRecoveryWithOuputAlreadyWritten(t *testing.T) {

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
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}
	log.Println("string(body)", string(body))
	if string(body) != "\"nope\"" {
		t.Fatalf("Invalid response '%s'", string(body))
	}
}
