package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// unic type for context
type ctxKey string

const requestIDKey ctxKey = "request_id"

// --- 1. TRACEABILITY AND LOGIN MIDDLEWARE ---
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Request ID (prod must use UUID)
		reqID := fmt.Sprintf("req-%d", time.Now().UnixNano())

		// Saving Request ID in Contexto to handlers/logs use it
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		r = r.WithContext(ctx)

		log.Printf("[%s] -> %s %s", reqID, r.Method, r.URL.Path)

		// Ejecutamos el siguiente handler en la cadena
		next.ServeHTTP(w, r)

		// Executed after response
		log.Printf("[%s] <- Completed: %v", reqID, time.Since(start))
	})
}

// --- 2. RECOVERY MIDDLEWARE ---
func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC RECOVERED]: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// --- HANDLER ---
func HandlerHello(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(requestIDKey).(string)

	fmt.Fprintf(w, "¡Hello! your Request ID is: %s\n", reqID)
}

func main() {
	mux := http.NewServeMux()

	// Handler normal
	mux.HandleFunc("/hello", HandlerHello)

	// chained middlewares:
	// flow: PanicRecovery -> Logging -> HandlerHello
	protectedHandle := PanicRecoveryMiddleware(LoggingMiddleware(mux))

	fmt.Println("Server runing :8080")
	http.ListenAndServe(":8080", protectedHandle)
}
