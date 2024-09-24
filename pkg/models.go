package parser

type Transaction struct {
	Type                 string          `json:"type"`
	Nonce                string          `json:"nonce"`
	To                   string          `json:"to"`
	From                 string          `json:"from"`
	Gas                  string          `json:"gas"`
	Hash                 string          `json:"hash"`
	Value                string          `json:"value"`
	Input                string          `json:"input"`
	MaxPriorityFeePerGas string          `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string          `json:"maxFeePerGas"`
	MaxFeePerBlobGas     string          `json:"maxFeePerBlobGas"`
	AccessList           []RPCAccessList `json:"accessList"`
	ChainID              string          `json:"chainId"`
	YParity              string          `json:"yParity"`
	BlockNumber          string          `json:"blockNumber"`
	BlockHash            string          `json:"blockHash"`
	TransactionIndex     string          `json:"transactionIndex"`
	R                    string          `json:"r"`
	S                    string          `json:"s"`
	V                    string          `json:"v"`
}

type RPCAccessList struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

type RPCBlock struct {
	Transactions []Transaction `json:"transactions"`
	Timestamp    string        `json:"timestamp"`
	Number       string        `json:"number"`
}

type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RPCResponse struct {
	Result string    `json:"result"`
	Error  *RPCError `json:"error"`
	ID     int       `json:"id"`
}

type RPCBatchRequest []RPCRequest

type RPCBlockResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	ID      int         `json:"id"`
	Error   *RPCError   `json:"error"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
