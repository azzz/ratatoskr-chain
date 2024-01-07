/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "See balance by an address",
	RunE: func(cmd *cobra.Command, args []string) error {
		bc := LoadBlockchain()

		address, err := cmd.Flags().GetString("address")
		if err != nil {
			return err
		}

		balance, err := bc.GetBalance(address)
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(w, "address\tamount\n")
		fmt.Fprintf(w, "%s\t%d\n", address, balance)
		w.Flush()

		return nil
	},
}

func init() {
	balanceCmd.Flags().String("address", "", "Blockchain address")
	balanceCmd.MarkFlagRequired("address")

	rootCmd.AddCommand(balanceCmd)
}
