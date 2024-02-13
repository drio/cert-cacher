package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FakeNodeGetter struct {
	NodeName string
}

func (fn *FakeNodeGetter) GetNode(r *http.Request) (string, error) {
	return fn.NodeName, nil
}

func Test_Ping(t *testing.T) {
	var store Store
	store = MemStore{}

	var ng NodeGetter
	ng = &FakeNodeGetter{}

	handler := createHandler(ng, store)
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expected := "pong"
	if string(body) != expected {
		t.Errorf("Expected response body to be %s, got %s", expected, body)
	}

}
