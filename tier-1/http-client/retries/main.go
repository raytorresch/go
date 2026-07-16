package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"tier-1/http-client/retries/assets"
	"time"
)

func main() {
	server := assets.NewServer()

	jsonPayload := []byte(`{"hsm_key":"master_key_1"}`)

	ctx := context.Background()

	err := requestWithRetries(ctx, server.URL, jsonPayload)
	if err != nil {
		fmt.Println("Final error, Limit attempts reached: ", err)
	}
}

func requestWithRetries(ctx context.Context, url string, payload []byte) error {
	client := &http.Client{Timeout: 2 * time.Second}

	maxRetries := 3
	baseBackoff := 1 * time.Second

	for try := 1; try <= maxRetries; try++ {

		//For each try new reader from same byte slices
		//NopCloser add empty Close() method for interface compliance
		bodyReader := bytes.NewReader(payload)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, io.NopCloser(bodyReader))

		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)

		//Verifing response to determinate if retry is necessary (503 response)
		retring := err != nil || (resp != nil && resp.StatusCode == http.StatusServiceUnavailable)

		if !retring {
			//200 ok cleanup and exit
			if resp != nil {
				defer resp.Body.Close()
				_, _ = io.Copy(io.Discard, resp.Body)
			}
			return nil
		}

		//backoff
		if try < maxRetries {
			if resp != nil {
				resp.Body.Close()
			}

			// Exponential Backoff : base * 2^(intento-1) -> 1s, 2s...
			waitingTime := baseBackoff * (1 << (uint(try) - 1))

			// jitter: add mileseconds between 0 and 200
			jitter := time.Duration(rand.Intn(200)) * time.Millisecond
			totalTime := waitingTime + jitter

			fmt.Printf("[Client] try %d failed. Retring en %v... \n", try, totalTime.Round(time.Millisecond))
			time.Sleep(totalTime)
		}
	}

	return fmt.Errorf("Reach allowed Attempts limit: %d", maxRetries)
}
