package tin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type SSEContext struct {
	*Context
	flusher http.Flusher
}

func (t *Context) SSE() *SSEContext {

	flusher, ok := t.Writer.(http.Flusher)

	if !ok {
		http.Error(t.Writer, "Streaming unsupported!", http.StatusInternalServerError)
		return &SSEContext{Context: t}
	}

	t.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	t.Writer.Header().Set("Cache-Control", "no-cache")
	t.Writer.Header().Set("Connection", "keep-alive")
	t.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	return &SSEContext{
		flusher: flusher,
		Context: t,
	}
}

func (t *SSEContext) Error(err error) error {

	return t.send("error", err.Error())
}

func (t *SSEContext) Event(event string) error {

	return t.send(event, nil)
}

func (t *SSEContext) JSON(v interface{}) error {

	return t.send("data", v)
}

func (t *SSEContext) Send(event string, v interface{}) error {
	return t.send(event, v)
}

func (t *SSEContext) send(event string, v interface{}) error {

	if t.clientGone || t.flusher == nil {
		return errors.New("client is gone")
	}

	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(t.Writer, "event: %s\n", event)
	if err != nil {
		return err
	}

	t.Context.headerWritten = true

	_, err = fmt.Fprintf(t.Writer, "data: %s\n\n", string(body))
	if err != nil {
		return err
	}

	if t.flusher != nil {
		t.flusher.Flush()
	}
	return nil
}

func (t *SSEContext) Gone() <-chan struct{} {
	return t.Request.Context().Done()
}
