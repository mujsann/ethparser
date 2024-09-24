package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	parser "mujsann.com/ethparser/pkg"
)

const SERVER_URL = "http://localhost:8080"

func TestMain(m *testing.M) {
	go main()
	time.Sleep(2 * time.Second)
	code := m.Run()
	os.Exit(code)
}

func TestE2E(t *testing.T) {

	address := os.Getenv("TEST_ADDRESS")
	if address == "" {
		t.Fatal("TEST_ADDRESS environment variable is not set", address)
	}

	isValid, err := parser.IsValidAddress(address)
	if !isValid || err != nil {
		t.Fatal("TEST_ADDRESS environment variable was not validated", err)
	}

	testGetCurrentBlock(t)
	testSubscribe(t, address)
	testGetTransactions(t, address)

}

func testSubscribe(t *testing.T, address string) {

	body := SubscribeRequestBody{Address: address}
	bodyBytes, _ := json.Marshal(body)
	resp, err := http.Post(SERVER_URL+"/subscribe", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// test - subscribe without address in the body should return error
	body = SubscribeRequestBody{}
	bodyBytes, _ = json.Marshal(body)

	resp, err = http.Post(SERVER_URL+"/subscribe", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func testGetCurrentBlock(t *testing.T) {

	resp, err := http.Get(SERVER_URL + "/current-block")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// test that response body has current_block and current_block must not be empty
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	respData := GetCurrentBlockResponseBody{}
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		t.Fatal(err)
	}

	if respData.CurrentBlock <= 0 {
		t.Errorf("Expected 'current_block' to be greater than zero")
	}
}

func testGetTransactions(t *testing.T, address string) {

	url := fmt.Sprintf(SERVER_URL+"/transactions/%s", address)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// test - get transactions with invalid address should return error
	invalidAddress := "Invalid_address"
	url = fmt.Sprintf(SERVER_URL+"/transactions/%s", invalidAddress)
	resp, err = http.Get(url)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}
