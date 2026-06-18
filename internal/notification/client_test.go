// client_test.go — tests our HTTP client logic against an in-process mock server.
// httptest.Server is a real listening server, but managed entirely inside the test.
package notification

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_SendEmail_Success(t *testing.T) {
	// Spin up a temporary real HTTP server just for this test
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/send", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var req SendEmailRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "user@test.com", req.To)

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(SendEmailResponse{
			MessageID: "msg_123",
			Status:    "queued",
		})
	}))
	defer server.Close() // always shut down, frees the port

	client := NewClient(server.URL) // point our client at the fake server
	resp, err := client.SendEmail(SendEmailRequest{
		To: "user@test.com", Subject: "Welcome", Body: "Hi there",
	})

	assert.NoError(t, err)
	assert.Equal(t, "msg_123", resp.MessageID)
}

func TestClient_SendEmail_VendorError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // simulate vendor outage
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.SendEmail(SendEmailRequest{To: "x@t.com", Subject: "s", Body: "b"})

	assert.Error(t, err) // our client correctly surfaces vendor failure
}

func TestClient_SendEmail_MalformedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json")) // simulate vendor returning garbage
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.SendEmail(SendEmailRequest{To: "x@t.com", Subject: "s", Body: "b"})

	assert.Error(t, err)
}
