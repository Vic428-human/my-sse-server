package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/events", sseHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("unable to start server: %S", err.Error())
	}

}

func UseTickerWithDefer() time.Ticker {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	return ticker
}

func UseTickerWithContext(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			fmt.Println("Ticker1....", t)
		}
	}
}

/*
備註:
使用context包可以更优雅地管理Ticker的生命周期。
*/

// SseHandler Basic SSE handler
func SseHandler(w http.ResponseWriter, r *http.Request) {
	// Set necessary HTTP headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ticker1 := UseTickerWithDefer()

	cpuT := time.NewTicker(time.Second)
	defer cpuT.Stop()

	clientGone := r.Context().Done()

}
