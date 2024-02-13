package main

import (
	"bytes"
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

func createServer(storeType string) *httptest.Server {
	var store Store
	store = MemStore{}
	if storeType == "disk" {
		store = DiskStore{}
	}
	store.Init()

	var ng NodeGetter
	ng = &FakeNodeGetter{NodeName: "foobar"}

	handler := createHandler(ng, store)
	server = httptest.NewServer(handler)
	return server
}

func TestPing(t *testing.T) {
	server = createServer("mem")
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

func TestSh(t *testing.T) {
	server = createServer("mem")
	defer server.Close()

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

func TestGetInvalidPath(t *testing.T) {
	server = createServer("mem")
	defer server.Close()

	resp, err := http.Get(server.URL + "/foo")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGetWhenEmpty(t *testing.T) {
	server = createServer("mem")
	defer server.Close()

	resp, err := http.Get(server.URL + "/cert")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	resp, err = http.Get(server.URL + "/key")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestAddCertAndReadWithDiskStore(t *testing.T) {
	generateCertFiles("foobar")
	server = createServer("disk")
	defer func() {
		server.Close()
		rmCerts("foobar")
	}()

	resp, err := http.Get(server.URL + "/cert")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("1 Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	resp, err = http.Get(server.URL + "/key")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("2 Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestDays(t *testing.T) {
	generateCertFiles("foobar")
	server = createServer("disk")
	defer func() {
		server.Close()
		rmCerts("foobar")
	}()

	resp, err := http.Get(server.URL + "/days")
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

	expected := "364"
	if string(body) != expected {
		t.Errorf("Expected response body to be [%s], got [%s]", expected, body)
	}
}

func TestAddCertAndReadWithMemStore(t *testing.T) {
	runTestAddCertAndReadWithMemStore(t, "cert")
	runTestAddCertAndReadWithMemStore(t, "key")
}

func runTestAddCertAndReadWithMemStore(t *testing.T, ext string) {
	generateCertFiles("foobar")
	server = createServer("mem")
	defer func() {
		server.Close()
		rmCerts("foobar")
	}()

	file, err := os.Open(fmt.Sprintf("foobar.%s", ext))
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(buf.Bytes())
	url := server.URL + fmt.Sprintf("/%s", ext)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", reader)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	resp, err = http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fileReader := bytes.NewReader(buf.Bytes())
	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		t.Errorf("error reading from file buffer %v", err)
	}
	if !bytes.Equal(bodyBytes, fileBytes) {
		t.Errorf("the data I got back from the server does not match what I sent")
	}
}
