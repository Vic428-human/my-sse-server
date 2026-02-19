package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {
	http.HandleFunc("/events", SseHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("unable to start server: %s", err.Error()) // log.Fatal 不支援 format string
	}

}

// ticker.Stop() 在 function return 的那一刻就被呼叫了
// 所以呼叫者拿到的 ticker 一出場就已經被 Stop() 掉
func UseTicker() *time.Ticker {
	return time.NewTicker(time.Second)
}

// SseHandler Basic SSE handler
func SseHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	memT := UseTicker()
	defer memT.Stop()

	cpuT := UseTicker()
	defer cpuT.Stop()

	clientGone := r.Context().Done()
	rc := http.NewResponseController(w)
	for {
		select {
		case <-clientGone:
			fmt.Println("Client has disconnet")
			return
		case <-memT.C:
			m, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("unable to get mem: %v", err)
				return
			}
			/*  兩個換行符 \n\n，這是 SSE 協議用來標示「一個事件結束」的必要格式
			event:mem
			data:Total: 16777216000,Used: 8388608000,Perc: 50.00%
			*/

			/*
							前端 JavaScript 接收到的內容
							const es = new EventSource("/events");

				es.addEventListener("mem", (event) => {
				    console.log(event.type);  // "mem"
				    console.log(event.data);  // "Total: 16777216000,Used: 8388608000,Perc: 50.00%" (event.data 會是純字串)
				});
			*/
			fmt.Fprintf(w, "event: mem\ndata: Total: %d,Used: %d,Perc: %.2f%%\n\n", m.Total, m.Used, m.UsedPercent)
			rc.Flush()
		case <-cpuT.C:
			c, err := cpu.Times(false)
			if err != nil {
				log.Printf("unable to get cpu: %v", err)
				return
			}
			if _, err := fmt.Fprintf(w, "event:cpu\ndata:User: %.2f, Sys: %.2f, Idle: %.2f\n\n", c[0].User, c[0].System, c[0].Idle); err != nil {
				log.Printf("unable to write: %s", err.Error())
				return
			}
			rc.Flush()
		}
	}
}
