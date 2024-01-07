package blockchain

import (
	"errors"
	"fmt"
	"github.com/azzz/ratatoskr/pkg/block"
	"github.com/azzz/ratatoskr/pkg/transaction"
)

type Blockchain struct {
	pow   ProofOfWork
	store Store
}

func (bc Blockchain) Iterator() Iterator {
	return Iterator{store: bc.store, head: bc.store.Tip()}
}

func (bc Blockchain) Tip() []byte {
	return bc.store.Tip()
}

func (bc Blockchain) AddBlock(transactions []transaction.Transaction) error {
	return bc.store.AddBlock(func(tip []byte) (block.Block, error) {
		newBlock := block.New(transactions, tip)
		return bc.pow.Sign(newBlock)
	})
}

func (bc Blockchain) addGenesis(address string) error {
	return bc.store.AddBlock(func(tip []byte) (block.Block, error) {
		if tip != nil {
			return block.Block{}, errors.New("genesis can be added only to an empty blockchain")
		}

		coinbase, err := transaction.NewCoinbaseTx(address, "")
		if err != nil {
			return block.Block{}, err
		}
		genesis := block.NewGenesis(coinbase)

		return bc.pow.Sign(genesis)
	})
}

func (bc Blockchain) GetBalance(address string) (uint64, error) {
	outs, err := bc.FindUTXO(address)
	if err != nil {
		return 0, fmt.Errorf("find unspent outputs: %w", err)
	}

	var sum uint64
	for _, out := range outs {
		sum += out.Value
	}

	return sum, nil
}

// FindUTXO returns all unspent outputs for the address.
func (bc Blockchain) FindUTXO(address string) ([]transaction.TxOutput, error) {
	var UTXOs []transaction.TxOutput
	unspent, err := bc.FindUnspentTransactions(address)

	if err != nil {
		return nil, fmt.Errorf("find unspent transactions: %w", err)
	}

	for _, tx := range unspent {
		for _, out := range tx.Vout {
			if out.CanUnclockWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs, nil
}

// FindUnspentTransactions returns all transactions with unspent Outputs (The outputs which are not used as inputs),
// unlockable by the address.
func (bc Blockchain) FindUnspentTransactions(address string) ([]transaction.Transaction, error) {
	var (
		iter         = bc.Iterator()
		unspent      = make(map[string]transaction.Transaction)
		spentOutputs = make(map[string]struct{})
	)

	for iter.Next() {
		if iter.Err() != nil {
			return nil, fmt.Errorf("retrieve block: %w", iter.Err())
		}

		for _, tx := range iter.Block().Transactions {
			for _, inTx := range tx.Vin {
				if !inTx.CanUnlockWith(address) {
					continue
				}

				key := fmt.Sprintf("%x;%d", inTx.TxID, inTx.Vout)
				spentOutputs[key] = struct{}{}
			}

			for outIdx, outTx := range tx.Vout {
				if !outTx.CanUnclockWith(address) {
					continue
				}

				key := fmt.Sprintf("%x;%d", tx.ID, outIdx)
				if _, ok := spentOutputs[key]; !ok {
					unspent[key] = tx
				}
			}
		}
	}

	result := make([]transaction.Transaction, 0, len(unspent))
	for _, tx := range unspent {
		result = append(result, tx)
	}

	return result, nil
}

// FindSpendableOutputs returns a minimal list of outputs needed to accumulate the required amount.
// First returning value: the accumulated amount
// Second returning value: map of transaction ID to output
func (bc Blockchain) FindSpendableOutputs(address string, amount uint64) ([]transaction.Output, error) {
	var (
		acc     uint64
		outputs = []transaction.Output{}
	)

	unspentTXs, err := bc.FindUnspentTransactions(address)
	if err != nil {
		return nil, fmt.Errorf("find unspent transactions: %w", err)
	}

mainLoop:
	for _, tx := range unspentTXs {
		for outIdx, out := range tx.Vout {
			if !out.CanUnclockWith(address) {
				continue
			}

			acc += out.Value
			outputs = append(outputs, transaction.Output{
				TxOutput: out,
				TxID:     tx.ID,
				Vout:     outIdx,
			})

			if acc >= amount {
				break mainLoop
			}
		}
	}

	return outputs, nil
}
