package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const batchSize = 1000 // Number of blocks to fetch in each batch
const rateLimit = 10

// generate a batch of block numbers
func (ep *ParserService) createBatch(first, last int64) []int64 {
	var b []int64
	for i := first; i <= last; i++ {
		b = append(b, i)
	}
	return b
}
func (p *ParserService) fetchBlocksByNumbers(blockNumbers []int64, fullTx bool) ([]*RPCBlock, error) {
	var batchReq RPCBatchRequest

	// data for batch request
	for _, blockNumber := range blockNumbers {
		batchReq = append(batchReq, RPCRequest{
			JSONRPC: "2.0",
			Method:  "eth_getBlockByNumber",
			Params:  []interface{}{fmt.Sprintf("0x%x", blockNumber), fullTx},
			ID:      generateID(),
		})
	}

	jsonData, err := json.Marshal(batchReq)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(p.RPC_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var batchResponse []RPCBlockResponse
	err = json.NewDecoder(resp.Body).Decode(&batchResponse)
	if err != nil {
		return nil, err
	}

	var blocks []*RPCBlock
	for _, res := range batchResponse {
		if res.Result != nil {
			blockData := res.Result
			blockBytes, _ := json.Marshal(blockData)

			var block RPCBlock
			err = json.Unmarshal(blockBytes, &block)
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, &block)
		}
	}

	return blocks, nil
}

// poll transactions in parallel
func (p *ParserService) pollEthTransactions(firstBlock, lastBlock int64, address string) []Transaction {

	var wg sync.WaitGroup
	ch := make(chan *RPCBlock, batchSize)
	doneCh := make(chan struct{})

	// use rate limiting
	rateLimiter := time.NewTicker(time.Second / rateLimit)

	var allTransactions []Transaction

	// worker to process blocks and filter transactions by to and from
	go func() {
		defer close(doneCh)
		for block := range ch {
			blockTimeHex := block.Timestamp

			// convert blocktime from hex to int64
			timestamp, err := convertHexToInt64(blockTimeHex)
			if err != nil {
				log.Printf("Error converting block timestamp: %v", err)
				continue
			}

			blockTime := time.Unix(timestamp, 0)

			// check if the block is older than the TransactionDaysLimit
			if p.TransactionDaysLimit > 0 && blockTime.Before(time.Now().Add(-time.Duration(p.TransactionDaysLimit)*24*time.Hour)) {
				log.Printf("Stopping processing when found block %s with date (time: %v) is older  than the TransactionDaysLimit %d day(s)", block.Number, blockTime, p.TransactionDaysLimit)
				return
			}

			for _, tx := range block.Transactions {
				if tx.From == address || tx.To == address {
					allTransactions = append(allTransactions, tx)
				}
			}

		}
		doneCh <- struct{}{}
	}()

	// request blocks in batches
	for blockNumber := lastBlock; blockNumber >= firstBlock; blockNumber -= batchSize {

		wg.Add(1)
		go func(start, end int64) {
			defer wg.Done()

			// rate limiting the request
			<-rateLimiter.C

			blocks, err := p.fetchBlocksByNumbers(p.createBatch(start, end), true)
			if err != nil {
				log.Printf("Error fetching blocks %d to %d: %v", start, end, err)
				return
			}

			for _, block := range blocks {
				select {
				case ch <- block:
				case <-doneCh:
					return
				}
			}

		}(blockNumber, min(blockNumber+batchSize-1, lastBlock))
	}

	// close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
		rateLimiter.Stop()
	}()

	<-doneCh

	return allTransactions

}
