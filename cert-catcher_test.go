package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type FakeNodeGetter struct {
	NodeName string
}

func (fn *FakeNodeGetter) GetNode(r *http.Request) (string, error) {
	return fn.NodeName, nil
}

var server *httptest.Server

func TestMain(m *testing.M) {
	fmt.Println("Set up stuff for tests here")

	var store Store
	store = MemStore{}

	var ng NodeGetter
	ng = &FakeNodeGetter{}

	handler := createHandler(ng, store)
	server = httptest.NewServer(handler)
	defer server.Close()

	exitVal := m.Run()
	os.Exit(exitVal)
}

func Test_Ping(t *testing.T) {
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

func Test_Sh(t *testing.T) {
	resp, err := http.Get(server.URL + "/sh")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	headerVal := resp.Header.Get("Content-Type")
	expected := "application/x-sh"
	if headerVal != expected {
		t.Errorf("incorrect header expected: %s got: %s", expected, headerVal)
	}
}
