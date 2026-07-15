package assets

import (
	"net/http"
	"net/http/httptest"
)

func CreateServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "SUCCESS")`))
	}))

	return server
}
