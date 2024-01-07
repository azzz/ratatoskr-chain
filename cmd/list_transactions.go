/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

// listTransactionsCmd represents the listTransactions command
var listTransactionsCmd = &cobra.Command{
	Use:   "list-transactions",
	Short: "List all transactions",
	Run: func(cmd *cobra.Command, args []string) {
		bc := LoadBlockchain()
		iter := bc.Iterator()

		for iter.Next() {
			cobra.CheckErr(iter.Err())
			b := iter.Block()

			for _, tx := range b.Transactions {
				ts := time.UnixMicro(int64(b.Timestamp))
				// block header
				fmt.Printf("%s: %x\n", ts.Format(time.DateTime), b.Hash)
				// tx header
				fmt.Printf("Tx: \t%x\n", tx.ID)

				fmt.Printf("\tInputs:\n")
				for _, inTx := range tx.Vin {
					if inTx.TxID == nil {
						fmt.Printf("\t\tcoinbase\n")
					} else {
						fmt.Printf("\t\t%x->%d\n", inTx.TxID, inTx.Vout)
					}
				}

				fmt.Printf("\tOutputs:\n")

				for outIdx, outTx := range tx.Vout {
					fmt.Printf("\t\t%d: Value=%d, Address=%q\n", outIdx, outTx.Value, outTx.ScriptPubKey)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listTransactionsCmd)
}
