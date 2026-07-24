package assets

import (
	"net/http"
	"net/http/httptest"
)

// Create a test server that simulates a hanging request (never responds)
func CreateHangingServer() *httptest.Server {
	Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a heavy process or deadlock: it waits forever
		<-r.Context().Done()
	}))
	return Server
}
