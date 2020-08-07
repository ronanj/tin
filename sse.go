package tin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type tinSSEContext struct {
	*tinContext
	flusher http.Flusher
}

func (t *tinContext) SSE() *tinSSEContext {

	flusher, ok := t.w.(http.Flusher)

	if !ok {
		http.Error(t.w, "Streaming unsupported!", http.StatusInternalServerError)
		return &tinSSEContext{tinContext: t}
	}

	t.w.Header().Set("Content-Type", "text/event-stream")
	t.w.Header().Set("Cache-Control", "no-cache")
	t.w.Header().Set("Connection", "keep-alive")
	t.w.Header().Set("Access-Control-Allow-Origin", "*")

	return &tinSSEContext{
		flusher:    flusher,
		tinContext: t,
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

func (t *tinSSEContext) send(event string, v interface{}) error {

	if t.clientGone || t.flusher == nil {
		return errors.New("Client is gone")
	}

	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(t.w, "event: %s\n", event)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(t.w, "data: %s\n\n", string(body))
	if err != nil {
		return err
	}

	if t.flusher != nil {
		t.flusher.Flush()
	}
	return nil
}
