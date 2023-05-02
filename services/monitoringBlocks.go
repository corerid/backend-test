package services

import (
	"encoding/json"
	"fmt"
	"github.com/corerid/backend-test/http"
	"github.com/corerid/backend-test/repositories"
	"github.com/corerid/backend-test/utils"
	"log"
	"math/big"
	"time"
)

type getLatestBlockNumberRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Id      int    `json:"id"`
}

type getLatestBlockNumberResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Id      int    `json:"id"`
}

type getBlockByNumberRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}
type getBlockByNumberResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		BaseFeePerGas   string `json:"baseFeePerGas"`
		Difficulty      string `json:"difficulty"`
		ExtraData       string `json:"extraData"`
		GasLimit        string `json:"gasLimit"`
		GasUsed         string `json:"gasUsed"`
		Hash            string `json:"hash"`
		LogsBloom       string `json:"logsBloom"`
		Miner           string `json:"miner"`
		MixHash         string `json:"mixHash"`
		Nonce           string `json:"nonce"`
		Number          string `json:"number"`
		ParentHash      string `json:"parentHash"`
		ReceiptsRoot    string `json:"receiptsRoot"`
		Sha3Uncles      string `json:"sha3Uncles"`
		Size            string `json:"size"`
		StateRoot       string `json:"stateRoot"`
		Timestamp       string `json:"timestamp"`
		TotalDifficulty string `json:"totalDifficulty"`
		Transactions    []struct {
			AccessList           []any  `json:"accessList,omitempty"`
			BlockHash            string `json:"blockHash"`
			BlockNumber          string `json:"blockNumber"`
			ChainID              string `json:"chainId"`
			From                 string `json:"from"`
			Gas                  string `json:"gas"`
			GasPrice             string `json:"gasPrice"`
			Hash                 string `json:"hash"`
			Input                string `json:"input"`
			MaxFeePerGas         string `json:"maxFeePerGas,omitempty"`
			MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas,omitempty"`
			Nonce                string `json:"nonce"`
			R                    string `json:"r"`
			S                    string `json:"s"`
			To                   string `json:"to"`
			TransactionIndex     string `json:"transactionIndex"`
			Type                 string `json:"type"`
			V                    string `json:"v"`
			Value                string `json:"value"`
		} `json:"transactions"`
		TransactionsRoot string `json:"transactionsRoot"`
		Uncles           []any  `json:"uncles"`
		Withdrawals      []struct {
			Address        string `json:"address"`
			Amount         string `json:"amount"`
			Index          string `json:"index"`
			ValidatorIndex string `json:"validatorIndex"`
		} `json:"withdrawals"`
		WithdrawalsRoot string `json:"withdrawalsRoot"`
	} `json:"result"`
}

func (s Service) MonitorBlockEthereum(specifiedAddress string) {
	jsonRPCServer := s.Config.JSONRPCServer

	runningBlockNumber, err := getLatestBlock(jsonRPCServer, s.Client)
	if err != nil {
		log.Fatal(err)
	}

	for {
		latestBlock, err := getLatestBlock(jsonRPCServer, s.Client)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Block number ", latestBlock)
		fmt.Println("Running Block number ", runningBlockNumber)

		// running > latest
		if runningBlockNumber.Cmp(latestBlock) == 1 {
			time.Sleep(5 * time.Second)
			continue
		}

		// getBlockByNumber
		block, err := getBlockByNumber(jsonRPCServer, s.Client, runningBlockNumber)
		if err != nil {
			log.Fatal(err)
		}

		// loop through transaction
		for _, transaction := range block.Result.Transactions {
			// check transaction from, to is specified address
			if transaction.From == specifiedAddress || transaction.To == specifiedAddress {
				fmt.Printf("Found transaction! %+v\n", transaction)
				transactionBytes, err := json.Marshal(transaction)
				if err != nil {
					return
				}

				blockNumber := new(big.Int)
				blockNumber.SetString(transaction.BlockNumber, 0)
				err = s.CreateTransactionRepo(repositories.Transaction{
					Address:         specifiedAddress,
					TransactionHash: transaction.Hash,
					BlockNumber:     blockNumber,
					Data:            transactionBytes,
				})
				if err != nil {
					log.Fatal(err)
				}
			}

		}

		runningBlockNumber.Add(runningBlockNumber, big.NewInt(1))
	}

}

func getLatestBlock(server string, client http.Client) (*big.Int, error) {
	reqBody, err := prepareBodyForGetLatestBlockNumber()
	if err != nil {
		return nil, err
	}
	req := http.NewRequest(reqBody, "application/json")
	response, err := client.Send(http.MethodPost, server, req)
	if err != nil {
		return nil, err
	}

	getLatestBlockNumber := &getLatestBlockNumberResponse{}
	err = http.ParseHTTPResponseToStruct(response, getLatestBlockNumber)
	if err != nil {
		return nil, err
	}

	blockNumber, err := utils.HexToBigInt(getLatestBlockNumber.Result)
	if err != nil {
		return nil, err
	}

	return blockNumber, nil
}

func getBlockByNumber(server string, client http.Client, blockNumber *big.Int) (getBlockByNumberResponse, error) {
	blockNumberHex, err := utils.BigIntToHex(blockNumber)
	if err != nil {
		return getBlockByNumberResponse{}, err
	}

	reqBody, err := prepareBodyForGetBlockByNumber(blockNumberHex)
	if err != nil {
		return getBlockByNumberResponse{}, err
	}

	req := http.NewRequest(reqBody, "application/json")
	response, err := client.Send(http.MethodPost, server, req)
	if err != nil {
		return getBlockByNumberResponse{}, err
	}
	if err != nil {
		return getBlockByNumberResponse{}, err
	}

	block := &getBlockByNumberResponse{}
	err = http.ParseHTTPResponseToStruct(response, block)
	if err != nil {
		return getBlockByNumberResponse{}, err
	}

	return *block, nil
}

func prepareBodyForGetLatestBlockNumber() ([]byte, error) {
	req := getLatestBlockNumberRequest{
		Jsonrpc: "2.0",
		Method:  "eth_blockNumber",
		Id:      0,
	}

	byteReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return byteReq, nil
}

func prepareBodyForGetBlockByNumber(blockNumberHex string) ([]byte, error) {
	req := getBlockByNumberRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params: []interface{}{
			blockNumberHex,
			true,
		},
		Id: 0,
	}

	byteReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return byteReq, nil
}
