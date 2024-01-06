package blockchain

import (
	"errors"

	"github.com/azzz/ratatoskr/pkg/proofofwork"
	"go.etcd.io/bbolt"
)

func Load(db *bbolt.DB) (Blockchain, error) {
	var tip []byte

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		if b == nil {
			return errors.New("blockchain is not initialized")
		}

		tip = b.Get(tipKey)

		return nil
	})

	if err != nil {
		return Blockchain{}, err
	}

	if tip == nil {
		return Blockchain{}, errors.New("blockchain is not initialized")
	}

	return Blockchain{tip: tip, db: db, pow: proofofwork.New()}, err
}
