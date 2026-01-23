package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_get(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	expected := testData{Name: "test", Value: 42}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := NewClient()
	var result testData
	err := client.get(server.URL, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != expected.Name || result.Value != expected.Value {
		t.Errorf("got %+v, want %+v", result, expected)
	}
}

func TestClient_get_error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient()
	var result struct{}
	err := client.get(server.URL, &result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClient_getBytes(t *testing.T) {
	expected := []byte(`{"test": "data"}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(expected)
	}))
	defer server.Close()

	client := NewClient()
	result, err := client.getBytes(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(result) != string(expected) {
		t.Errorf("got %s, want %s", result, expected)
	}
}
