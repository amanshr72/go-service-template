package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type sendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type sendEmailResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/send", handleSend)

	// Simulate occasional vendor failure — testing/retry/error-handling logic
	mux.HandleFunc("POST /v1/send-flaky", handleFlakySend)

	log.Println("Mock notification server running on :8089")
	log.Fatal(http.ListenAndServe(":8089", mux))
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	var req sendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Basic contract validation
	if req.To == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "to field required"})
		return
	}

	log.Printf("mock: received email request to=%s subject=%s", req.To, req.Subject)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(sendEmailResponse{
		MessageID: randomID(),
		Status:    "queued",
	})
}

// handleFlakySend randomly fails ~30% of the time — simulates real-world vendor instability
func handleFlakySend(w http.ResponseWriter, r *http.Request) {
	if rand.Float64() < 0.3 {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "vendor temporarily unavailable"})
		return
	}
	handleSend(w, r)
}

func randomID() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 10)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = chars[r.Intn(len(chars))]
	}
	return "msg_" + string(b)
}
