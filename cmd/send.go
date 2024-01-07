/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/azzz/ratatoskr/pkg/transaction"

	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send money between addresses",
	Run: func(cmd *cobra.Command, args []string) {
		sender, err := cmd.Flags().GetString("from")
		cobra.CheckErr(err)

		receiver, err := cmd.Flags().GetString("to")
		cobra.CheckErr(err)

		amount, err := cmd.Flags().Getuint("amount")
		cobra.CheckErr(err)

		bc := LoadBlockchain()
		outputs, err := bc.FindSpendableOutputs(sender, amount)
		cobra.CheckErr(err)

		tx, err := transaction.NewUTXOTransaction(sender, receiver, amount, outputs)
		cobra.CheckErr(err)

		err = bc.AddBlock([]transaction.Transaction{tx})
		cobra.CheckErr(err)

		fmt.Println("successfully sent")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().String("from", "", "Sender address")
	sendCmd.Flags().String("to", "", "Receiver address")
	sendCmd.Flags().uint("amount", 0, "Amount")

	sendCmd.MarkFlagRequired("from")
	sendCmd.MarkFlagRequired("to")
	sendCmd.MarkFlagRequired("amount")
}
