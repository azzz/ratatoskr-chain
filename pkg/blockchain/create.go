package blockchain

import (
	"fmt"
	"github.com/azzz/ratatoskr/pkg/transaction"

	"github.com/azzz/ratatoskr/pkg/block"
	"github.com/azzz/ratatoskr/pkg/proofofwork"
	"go.etcd.io/bbolt"
)

func Create(db *bbolt.DB, address string) (Blockchain, error) {
	var (
		tip []byte
		pow = proofofwork.New()
	)

	err := db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucket(blocksBucket)
		if err != nil {
			return err
		}

		coinbase := transaction.NewCoinbaseTx(address, "")
		genesis := block.NewGenesis(coinbase)
		block, err := pow.Sign(genesis)
		if err != nil {
			return fmt.Errorf("proof-of-work: %w", err)
		}

		encoded, err := block.Serialize()
		if err != nil {
			return fmt.Errorf("serialize: %w", err)
		}

		if err := b.Put(block.Hash, encoded); err != nil {
			return fmt.Errorf("save block: %w", err)
		}

		if err := b.Put(tipKey, block.Hash); err != nil {
			return fmt.Errorf("save tip: %w", err)
		}

		tip = block.Hash

		return nil
	})

	return Blockchain{tip: tip, db: db, pow: pow}, err
}
