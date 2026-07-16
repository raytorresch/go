package assets

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
)

func NewServer() *httptest.Server {
	var tries int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tries := atomic.AddInt64(&tries, 1)

		body, _ := io.ReadAll(r.Body)

		if tries <= 2 {
			fmt.Printf("[SERVER] Try %d: Payload recibed `%s`. 503(Overload)", tries, string(body))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		fmt.Printf("[SERVER] Try %d: Payload recibed `%s`. 200 OK", tries, string(body))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))

	return server
}
