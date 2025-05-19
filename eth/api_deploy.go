package eth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// DeployStatusAPI provides RPC methods to track deployment status
type DeployStatusAPI struct {
	backend *Ethereum
}

// GetDeployStatus checks the status of a contract deployment transaction.
// Possible return values:
// - "success:<contract_address>"
// - "reverted"
// - "pending"
// - "not_found"
func (api *DeployStatusAPI) GetDeployStatus(ctx context.Context, txHash common.Hash) (string, error) {
	lookup, _, err := api.backend.BlockChain().GetTransactionLookup(txHash)
	if err != nil || lookup == nil {
		// fallback to mempool
		if tx := api.backend.TxPool().Get(txHash); tx != nil {
			return "pending", nil
		}
		return "not_found", nil
	}

	// Now try to get the receipt for the block
	receipts := api.backend.BlockChain().GetReceiptsByHash(lookup.BlockHash)
	if receipts == nil {
		return "not_found", nil
	}

	for _, r := range receipts {
		if r.TxHash == txHash {
			if r.Status == types.ReceiptStatusSuccessful && r.ContractAddress != (common.Address{}) {
				return fmt.Sprintf("success:%s", r.ContractAddress.Hex()), nil
			} else {
				return "reverted", nil
			}
		}
	}

	return "not_found", nil
}
