// This example demonstrates the importance of using timeouts in HTTP requests to avoid hanging goroutines.
package main

import (
	// "context"
	"context"
	"fmt"
	"net/http"
	"runtime"
	"tier-1/http-client/assets"
	"time"
)

func noTimeoutRequest(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("[Bad Practice] error: %w", err)
		return nil, err
	}
	if res != nil {
		defer res.Body.Close()
	}
	fmt.Println("[Bad Practice]: Request completed successfully!")
	return res, nil
}

func timeoutRequest(url string) (*http.Response, error) {
	client := assets.NewSecureHttpClient(3 * time.Second)

	ctx := context.Background()
	res, err := client.Get(ctx, url)

	if err != nil {
		fmt.Println("[Good Practice] error: %w", err)
		return nil, err
	}

	if res != nil {
		defer res.Body.Close()
	}

	fmt.Println("[Good Practice]: Request completed successfully!")
	return res, nil
}

func main() {
	server := assets.CreateHangingServer()

	fmt.Printf("=== Test server runing %s ===\n", server.URL)
	fmt.Printf("Inital Goroutines: %d\n\n", runtime.NumGoroutine())

	// ---SCENARY 1: BAD PRACTICE (NO TIMEOUT) ---
	fmt.Println("--- Executing Scenary 1 (Bad practice: no Timeout) ---")

	// Lanzamos la petición en una goroutine para que no bloquee este hilo principal
	go func() {
		_, err := noTimeoutRequest(server.URL)
		if err != nil {
			fmt.Println("Error in bad practice request:", err)
		}
	}()

	// Wait 2 seconds to let the bad practice request hang
	time.Sleep(2 * time.Second)
	fmt.Printf("Live Goroutines due bad practices: %d\n\n", runtime.NumGoroutine())

	// ---SCENARY 2: GOOD PRACTICE (WITH TIMEOUT) ---
	fmt.Println("--- Executing Scenary 2 (Good practice: with Timeout) ---")

	start := time.Now()
	resp, err := timeoutRequest(server.URL)
	fmt.Printf("Live Goroutines due good practices: %d\n\n", runtime.NumGoroutine())
	if err != nil {
		fmt.Println("Error in good practice request:", err)
		fmt.Printf("Request canceled after: %s\n", time.Since(start))
	} else {
		fmt.Println("Good practice request completed with status code:", resp.StatusCode)
	}

	time.Sleep(2 * time.Second) // Wait a bit to let the goroutine finish
	fmt.Printf("Live Goroutines after good practice: %d\n", runtime.NumGoroutine())
	fmt.Println("Note: The bad practice request is still hanging and will not complete, while the good practice request was canceled after the timeout.")
	defer server.Close()
}
