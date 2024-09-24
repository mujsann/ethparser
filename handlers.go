package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	parser "mujsann.com/ethparser/pkg"
)

type SubscribeRequestBody struct {
	Address string `json:"address"`
}
type GetCurrentBlockResponseBody struct {
	CurrentBlock int `json:"current_block"`
}

// Handler for getting the current block
// Method: Get
func (app *App) getCurrentBlockHandler(w http.ResponseWriter, _ *http.Request) {
	log.Print("request to get current block")
	block := app.svc.GetCurrentBlock()
	json.NewEncoder(w).Encode(&GetCurrentBlockResponseBody{
		CurrentBlock: block,
	})
}

// Handler for subscribing to updates
// Method: Post
// Body: {address: <string>}
func (app *App) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	bodyData := &SubscribeRequestBody{}
	err = json.Unmarshal(body, bodyData)
	if err != nil {
		http.Error(w, "Could not unmarshal request body", http.StatusBadRequest)
		return
	}

	address := bodyData.Address

	// check if address is valid
	if address == "" {
		http.Error(w, "Could not add the adress: address not provided", http.StatusBadRequest)
	}

	response := app.svc.Subscribe(address)

	if !response {
		http.Error(w, "Address was already added to subscribers", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler for getting all transactions
// Route: /transactions/<address>
// Method: Get
// Params: {address: <string>}
func (app *App) getTransactionsHandler(w http.ResponseWriter, r *http.Request) {

	// extract address from the url path
	address := r.URL.Path[len("/transactions/"):]

	// check if address is valid
	isValid, err := parser.IsValidAddress(address)

	if !isValid {
		http.Error(w, fmt.Sprintf("could not get transactions: %s", err), http.StatusBadRequest)
		return
	}

	// get transactions from the parser service
	response := app.svc.GetTransactions(address)
	if response == nil {
		http.Error(w, "failed to get transactions for this address", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
