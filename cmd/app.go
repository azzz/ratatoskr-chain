package cmd

import (
	"path"

	"github.com/azzz/ratatoskr/pkg/blockchain"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

func CreateBlockchain(address string) blockchain.Blockchain {
	db, err := bbolt.Open(path.Join(config.DataDir, "database.db"), 0600, nil)
	cobra.CheckErr(err)

	bc, err := blockchain.CreatePersistent(db, address)
	cobra.CheckErr(err)

	return bc
}

func LoadBlockchain() blockchain.Blockchain {
	db, err := bbolt.Open(path.Join(config.DataDir, "database.db"), 0600, nil)
	cobra.CheckErr(err)

	bc, err := blockchain.LoadPersistent(db)
	cobra.CheckErr(err)

	return bc
}
