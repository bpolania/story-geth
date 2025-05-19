package eth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// RawCall represents a single RPC call
type RawCall struct {
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

// BatchExecAPI provides the execBatch RPC
type BatchExecAPI struct {
	backend *Ethereum
}

// ExecBatch executes a list of supported RPC calls in order
func (api *BatchExecAPI) ExecBatch(ctx context.Context, calls []RawCall) ([]interface{}, error) {
	results := make([]interface{}, 0, len(calls))

	for _, call := range calls {
		var res interface{}

		switch call.Method {
		case "eth_getBalance":
			var args [2]string

			paramsBytes, _ := json.Marshal(call.Params)
			if err := json.Unmarshal(paramsBytes, &args); err != nil {
				res = map[string]string{"error": "invalid params"}
			} else {
				bal, err := api.ethGetBalance(ctx, args[0], args[1])
				if err != nil {
					res = map[string]string{"error": err.Error()}
				} else {
					res = bal
				}
			}

		default:
			res = map[string]string{"error": "unsupported method: " + call.Method}
		}

		results = append(results, res)
	}

	return results, nil
}

// ethGetBalance emulates the eth_getBalance RPC call
func (api *BatchExecAPI) ethGetBalance(ctx context.Context, addr string, blockTag string) (string, error) {
	address := common.HexToAddress(addr)

	var blockNumber rpc.BlockNumber
	if err := blockNumber.UnmarshalJSON([]byte(`"` + blockTag + `"`)); err != nil {
		return "", err
	}

	header, err := api.backend.APIBackend.HeaderByNumber(ctx, blockNumber)
	if err != nil {
		return "", err
	}

	state, err := api.backend.BlockChain().StateAt(header.Root)
	if err != nil {
		return "", err
	}

	balance := state.GetBalance(address)
	return fmt.Sprintf("0x%x", balance), nil
}
