package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"tier-1/http-server/graceful-shoutdown/assets"
	"time"
)

func makeRequest(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("[CLIENT %s]: Network error: %v\n", url, err)
		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("[CLIENT %s] Server response: Code %d - Body: %s\n", url, resp.StatusCode, string(body))

}

func main() {
	//wiaiting for srv
	srv := assets.NewServer()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[SERVER] ListenAndServe error: %v\n", err)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	// client simulation
	var wg sync.WaitGroup
	wg.Add(2)
	go makeRequest("http://localhost:9999/fast", &wg)
	go makeRequest("http://localhost:9999/slow", &wg)

	time.Sleep(1 * time.Second)

	// Init Graceful shoutdown
	ctxGraceContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	err := srv.Shutdown(ctxGraceContext)
	shoutdownDuration := time.Since(start).Round(time.Microsecond)

	// results
	if err != nil {
		//slow request
		fmt.Printf("\n[MAIN] Failed Shutdownd after %v. Err: %v\n", shoutdownDuration, err)
	} else {
		fmt.Printf("\n[MAIN]  Successful  Shutdown after %v\n", shoutdownDuration)
	}

	// New request verifing server rejection
	fmt.Println("[MAIN] Request post-shutdown...")
	resp, errReq := http.Get("http://localhost:9999/rapida")
	if errReq != nil {
		fmt.Printf("[New Client] Refused conection: %v\n", errReq)
	} else {
		resp.Body.Close()
	}

	wg.Wait()
	fmt.Println("[MAIN]: Process finished")
}
