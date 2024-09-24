package parser

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"math/rand"
)

var (
	mu sync.Mutex
	r  = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// a simple id generator
func generateID() int {
	mu.Lock()
	defer mu.Unlock()
	return r.Intn(1000000)
}

// isValidAddress checks if the given Ethereum address is valid by querying its balance
// if the address is valid, it will return a balance (even if it's zero). If the address is invalid, it will return an error.
func IsValidAddress(address string) (bool, error) {

	if address == "" {
		return false, fmt.Errorf("address should not be empty")
	}

	cleanAdd := strings.TrimPrefix(address, "0x")
	_, err := hex.DecodeString(cleanAdd)
	if err != nil {
		return false, fmt.Errorf("address %s is not a hex string: %s", address, err)
	}

	body := RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBalance",
		Params:  []interface{}{address, "latest"},
		ID:      generateID(),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(RPC_URL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	respData := RPCResponse{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return false, err
	}

	return true, nil
}

func convertHexToInt64(data string) (int64, error) {
	return strconv.ParseInt(data[2:], 16, 64)
}
