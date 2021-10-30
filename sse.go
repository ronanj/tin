package tin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type tinSSEContext struct {
	*Context
	flusher http.Flusher
}

func (t *Context) SSE() *tinSSEContext {

	flusher, ok := t.Writer.(http.Flusher)

	if !ok {
		http.Error(t.Writer, "Streaming unsupported!", http.StatusInternalServerError)
		return &tinSSEContext{Context: t}
	}

	t.Writer.Header().Set("Content-Type", "text/event-stream")
	t.Writer.Header().Set("Cache-Control", "no-cache")
	t.Writer.Header().Set("Connection", "keep-alive")
	t.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	return &tinSSEContext{
		flusher: flusher,
		Context: t,
	}
}

func (t *tinSSEContext) Error(err error) error {

	return t.send("error", err.Error())
}

func (t *tinSSEContext) Event(event string) error {

	return t.send(event, nil)
}

func (t *tinSSEContext) JSON(v interface{}) error {

	return t.send("data", v)
}

func (t *tinSSEContext) Send(event string, v interface{}) error {
	return t.send(event, v)
}

func (t *tinSSEContext) send(event string, v interface{}) error {

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

	_, err = fmt.Fprintf(t.Writer, "data: %s\n\n", string(body))
	if err != nil {
		return err
	}

	if t.flusher != nil {
		t.flusher.Flush()
	}
	return nil
}

func (t *tinSSEContext) Gone() <-chan struct{} {
	return t.Request.Context().Done()
}
