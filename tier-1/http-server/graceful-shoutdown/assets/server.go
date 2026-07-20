package assets

import (
	"fmt"
	"net/http"
	"time"
)

func NewServer() *http.Server {
	//router
	mux := http.NewServeMux()

	//fast endpoint
	mux.HandleFunc("/fast", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[SERVER]: Init fast request 2 sec")
		time.Sleep(2 * time.Second)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Fast OK"))
		fmt.Println("[SERVER]: Fast request ended")
	})

	//slow request
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[SERVER]: Init slow request 8 sec")
		time.Sleep(8 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Slow OK"))
		fmt.Println("[SERVER]: slow request ender")
	})

	srv := &http.Server{
		Addr:    ":9999",
		Handler: mux,
	}

	return srv
}
