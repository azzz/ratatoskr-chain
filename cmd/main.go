package main

import (
	"context"
	"os"
	"path"

	"github.com/azzz/ratatoskr/pkg/chain"
	"go.etcd.io/bbolt"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dbpath := path.Join(home, ".ratatoskr.db")
	db, err := bbolt.Open(dbpath, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var (
		ctx = context.Background()
	)

	blockchain, err := chain.NewBlockChainFromState(ctx, db)
	if err != nil {
		panic(err)
	}

	_ = blockchain
}
