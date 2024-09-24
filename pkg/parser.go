package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []Transaction
}

type ParserService struct {
	mu sync.Mutex

	CurrentBlock         string // hexadecimal
	Subscribers          map[string]bool
	RPC_URL              string
	TransactionDaysLimit int64
}

const RPC_URL = "https://ethereum-rpc.publicnode.com"

func NewParser(txLimit int64) ParserService {

	return ParserService{
		mu: sync.Mutex{},

		Subscribers:          make(map[string]bool),
		RPC_URL:              RPC_URL,
		TransactionDaysLimit: txLimit,
	}
}

// GetCurrentBlock retrieves the last parsed block
func (p *ParserService) GetCurrentBlock() int {
	block, err := p.fetchCurrentBlock()
	if err != nil {
		return 0
	}
	return int(block)
}

// Adds an address to the list of observers
func (p *ParserService) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.Subscribers[address]; exists {
		return false
	}
	p.Subscribers[address] = true
	return true
}

// Gets all transactions for a given address from recent blocks
func (p *ParserService) GetTransactions(address string) []Transaction {
	return p.pollEthTransactions(int64(0), int64(p.GetCurrentBlock()), address)
}

// Makes a call to get the current block
func (p *ParserService) fetchCurrentBlock() (int64, error) {
	body := RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		ID:      generateID(),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(RPC_URL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return 0, err
	}

	var respData RPCResponse
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		log.Printf("Error unmarshalling JSON-RPC response: %v", err)
		return 0, err
	}

	if respData.Error != nil {
		return 0, fmt.Errorf("RPC Error: %s", respData.Error.Message)
	}

	block, err := convertHexToInt64(respData.Result)
	if err != nil {
		return 0, err
	}

	p.CurrentBlock = respData.Result
	return block, nil
}
