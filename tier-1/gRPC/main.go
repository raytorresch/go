package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tier-1/gRPC/assets"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//Launch servers HTTP gRCP
	go assets.CreateHTTPServer()
	go assets.CreategRPCServer()
	time.Sleep(200 * time.Millisecond) // give both listeners time to bind before the client dials

	//  crypt data to test
	hashHex := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	// crypto signature Base64 (Plain text)
	signatureBase64 := "MEQCIGS1fX8291a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c82a9e8b=="

	// Transform signature Base64/Hex to RAW BYTES (to binari gRPC)
	rawSignatureBytes, _ := hex.DecodeString("3044022024b57d7f36f756b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b102202a9e8b")

	requestNum := 1000

	fmt.Printf("INIT BENCK MARK (%d Requests) ===\n\n", requestNum)

	// =========================================================================
	// TEST 1: HTTP/1.1 + JSON
	// =========================================================================
	clientHTTP := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10, // Actived connections reuse
		},
	}

	payloadStruct := assets.CryptoRequestJSON{Hash: hashHex, SignatureB64: signatureBase64}
	jsonBytes, _ := json.Marshal(payloadStruct)

	fmt.Printf("[HTTP/JSON] Payload size: %d bytes\n", len(jsonBytes))

	startJSON := time.Now()
	var totalBytesHTTP int64

	for i := 0; i < requestNum; i++ {
		req, _ := http.NewRequest("POST", "http://127.0.0.1:8081/validate", bytes.NewReader(jsonBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := clientHTTP.Do(req)
		if err != nil {
			fmt.Println("Error HTTP:", err)
			return
		}

		// Leemos y descartamos para mantener el socket TCP libre
		body, _ := io.ReadAll(resp.Body)
		totalBytesHTTP += int64(len(jsonBytes) + len(body))
		resp.Body.Close()
	}

	durationHTTP := time.Since(startJSON)

	// =========================================================================
	// TEST 2: gRPC (HTTP/2 + Protobuf Binario Puro)
	// =========================================================================
	conn, err := grpc.Dial("127.0.0.1:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	payloadBinario := append([]byte(hashHex), rawSignatureBytes...)
	fmt.Printf("[Binary gRPC] Payload sizes: %d bytes (Without Base64 overload)\n\n", len(payloadBinario))

	startGRPC := time.Now()
	var totalBytesGRPC int64

	ctx := context.Background()
	for i := 0; i < requestNum; i++ {
		var reply []byte
		// Remote gRPC invocation over multiplex connection HTTP/2
		err := conn.Invoke(ctx, "/crypto.CryptoService/ValidateBinarySgnature", payloadBinario, &reply, grpc.ForceCodec(assets.RawCodec{}))
		if err != nil {
			fmt.Println("Error gRPC:", err)
			return
		}
		totalBytesGRPC += int64(len(payloadBinario) + len(reply))
	}
	durationGRPC := time.Since(startGRPC)

	// =========================================================================
	// RESULTS & TRAFFIC COMPARATION
	// =========================================================================
	fmt.Println("================== WEB RESULTS ==================")
	fmt.Printf("Total time HTTP/1.1 JSON: %v (%.2f ms/req)\n", durationHTTP, float64(durationHTTP.Milliseconds())/float64(requestNum))
	fmt.Printf("Total time gRPC HTTP/2:  %v (%.2f ms/req)\n", durationGRPC, float64(durationGRPC.Milliseconds())/float64(requestNum))

	fmt.Println("\n WEB TRAFFIC in WIRE (1,000 requests):")
	fmt.Printf("HTTP/JSON Transfered: ~%.2f KB\n", float64(totalBytesHTTP)/1024)
	fmt.Printf("gRPC Transfered:      ~%.2f KB\n", float64(totalBytesGRPC)/1024)

	ahorroPercent := (1.0 - (float64(totalBytesGRPC) / float64(totalBytesHTTP))) * 100
	fmt.Printf("\nReal Band Width Save with gRPC + Pure Binary: %.2f%%\n", ahorroPercent)
}
