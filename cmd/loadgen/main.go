package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{Timeout: 3 * time.Second}
	endpoints := []struct {
		method string
		url    string
	}{
		{"GET", "http://localhost:8080/health"},
		{"GET", "http://localhost:8080/api/v1/products/"},
		{"GET", "http://localhost:8080/api/v1/products/999"},
		{"GET", "http://localhost:8080/api/v1/users"},
		{"GET", "http://localhost:8080/metrics"},
	}

	fmt.Println("Load generator running... Ctrl+C to stop")
	for {
		e := endpoints[rand.Intn(len(endpoints))]
		req, _ := http.NewRequest(e.method, e.url, bytes.NewBuffer(nil))
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
		}
		time.Sleep(time.Duration(50+rand.Intn(150)) * time.Millisecond)
	}
}