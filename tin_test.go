package tin_test

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ronanj/tin"
)

func TestTin_Shutdown(t *testing.T) {
	ti := tin.Default()

	rand.Seed(time.Now().Unix())
	socketPath := fmt.Sprintf("/tmp/tin-socket-%d.sock", rand.Int())
	defer os.Remove(socketPath)
	ready := make(chan error)

	t.Logf("Using socket %s", socketPath)

	go func() {
		// This is a blocking call if no error
		err := ti.Run("unix:" + socketPath)
		ready <- err
	}()

	select {
	case <-time.After(time.Millisecond * 200):
		// Assume the server has been started correctly

	case err := <-ready:
		if err != nil {
			t.Fatalf("Run() failed: %v", err)
		}

	}

	httpUnixClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}

	// Send a request to the server
	resp, err := httpUnixClient.Get("http://unix:" + socketPath + "/status")
	if err != nil {
		t.Errorf("GET() failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET() failed: %v", resp.StatusCode)
	}

	gotErr := ti.Shutdown(context.Background())
	if gotErr != nil {
		t.Errorf("Shutdown() failed: %v", gotErr)
		return
	}

	// Send a request to the server
	resp, err = httpUnixClient.Get("http://unix:" + socketPath + "/status")
	if err == nil {
		t.Errorf("GET() should have failed after shutdown")
	}
}
