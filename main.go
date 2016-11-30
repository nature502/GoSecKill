package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var once sync.Once

// Payload defines seckill request body.
type payload struct {
	// id reprensets a sku.
	Id string `json:"id"`
}

func secKillHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read the body into a string for json decoding
	payload := &payload{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(data, payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// let's create a job with the payload
	job := Job{payload: payload}
	// Push the work onto the queue.

	select {
	case JobQueue <- job:
	// Keep request for waiting to long, when JobQueue is full.
	case <-time.After(time.Duration(100) * time.Millisecond):
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Sold out.")
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Get tickets.")
}

func init() {
	once.Do(func() {
		// Set max workers as 100
		dispatcher := NewDispatcher(100)
		dispatcher.Run()
	})
}

func main() {
	http.HandleFunc("/", secKillHandler)
	http.ListenAndServe(":8080", nil)
}
