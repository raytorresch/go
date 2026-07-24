// Lab: connection reuse in net/http.
//
// Runs the same 150-request burst against two clients that only differ
// in idle connection pool size (MaxIdleConns / MaxIdleConnsPerHost):
// a "default" client with a tiny pool (2/2) and a "custom" client with
// a large one (100/100). An httptrace.ClientTrace tags each connection
// as new or reused, so the printed metrics show how pool size directly
// controls how many TCP connections get reused across rounds instead
// of being redialed.
package main

import (
	"context"
	"fmt"
	"io"
	"sync"
	"tier-1/http-client/connection-reuse/assets"
	"time"
)

func makeRequest(client assets.HttpClient, serverURL string) {
	var wg sync.WaitGroup
	totalRequests := 150

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := client.Get(context.Background(), serverURL)
			if err != nil {
				fmt.Printf("[Goroutine %d] Error: %v\n", i, err)
				return
			}

			// Clean up the response body to allow connection reuse
			defer resp.Body.Close()
			_, _ = io.Copy(io.Discard, resp.Body)
			// if err != nil {
			// 	fmt.Printf("[Goroutine %d] Error reading response body: %v\n", i, err)
			// 	return
			// }

			// fmt.Printf("[Goroutine %d] Request completed successfully!, Reques Code: %d \n", i, resp.StatusCode)

		}()
	}
	wg.Wait()
}

func simulation(client assets.HttpClient, metrics *assets.MetricConnection, serverURL string) {
	// First requests round to establish connections
	metrics.Reset()
	makeRequest(client, serverURL)
	roundOne := fmt.Sprintf("Round 1: New Connections: %d, Reused Connections: %d", metrics.NewCon, metrics.ReuseCon)

	time.Sleep(50 * time.Millisecond) // Wait a bit before the next round

	// Second requests round to test connection reuse
	metrics.Reset()
	makeRequest(client, serverURL)
	roundTwo := fmt.Sprintf("Round 2: New Connections: %d, Reused Connections: %d", metrics.NewCon, metrics.ReuseCon)

	fmt.Printf("%s, %s \n", roundOne, roundTwo)
}

func main() {
	// Create a test server that simulates a real server
	server := assets.CreateServer()

	// Create a "default" custom transport with connection reuse settings
	defaulTransport := assets.NewTransport(2, 2, 90*time.Second)
	defaultMetrics := &assets.MetricConnection{}
	defaultClient := assets.NewSecureHttpClient(3*time.Second, defaulTransport, defaultMetrics)

	fmt.Println("=== Transaction with default client ===")
	simulation(defaultClient, defaultMetrics, server.URL)

	// Create a "custom" custom transport with connection reuse settings
	customTransport := assets.NewTransport(100, 100, 90*time.Second)
	customMetrics := &assets.MetricConnection{}
	customClient := assets.NewSecureHttpClient(3*time.Second, customTransport, customMetrics)

	fmt.Println("=== Transacitons with costume client")

	simulation(customClient, customMetrics, server.URL)

	defer server.Close()
}
