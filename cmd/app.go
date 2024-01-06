package cmd

import (
	"path"

	"github.com/azzz/ratatoskr/pkg/chain"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

func NewBlockchain() chain.Blockchain {
	db, err := bbolt.Open(path.Join(config.DataDir, "database.db"), 0600, nil)
	cobra.CheckErr(err)

	bc, err := chain.NewBlockchain(db, chain.NewSimpleHashCash(24))
	cobra.CheckErr(err)

	return bc
}
