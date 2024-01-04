package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/azzz/ratatoskr/pkg"
)

func main() {
	var (
		t time.Time
	)

	pow := pkg.NewSimpleHashCash(0)
	blockchain := pkg.NewBlockchain(pow)

	ctx := context.Background()

	t = time.Now()
	log.Println("Initialize BlockChain")
	if err := blockchain.Init(ctx); err != nil {
		panic(err)
	}
	log.Printf("BlockChain initialized in %s\n", time.Now().Sub(t))

	t = time.Now()
	log.Println("Add a block")
	if err := blockchain.Add(ctx, "Hello World"); err != nil {
		panic(err)
	}
	log.Printf("Block added in %s\n", time.Now().Sub(t))

	genesis, _ := blockchain.Block(0)
	block, _ := blockchain.Block(1)

	fmt.Printf("genesis: %s\n", genesis)
	fmt.Printf("block: %s\n", block)
}
