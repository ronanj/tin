package tin

import (
	"io"
	"net/http/httptest"
	"testing"
)

func TestSSE_SendJSON(t *testing.T) {

	tin := New()
	req := httptest.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()

	tin.GET("/path", func(c *Context) {
		sse := c.SSE()
		sse.JSON(0)
	})

	tin.router.ServeHTTP(w, req)
	res := w.Result()
	if res.Status != "200 OK" {
		t.Fatalf("Expect context to be ok: %s!=200", res.Status)
	}

	if res.Header.Get("Content-Type") != "text/event-stream; charset=utf-8" {
		t.Fatalf("Content type not set: '%s'", res.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("Could not ready the body")
	}
	if string(body) != "event: data\ndata: 0\n\n" {
		t.Fatalf("Invalid response: '%s'", string(body))
	}

}
