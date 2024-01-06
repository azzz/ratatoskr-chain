/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// genesisCmd represents the genesis command
var genesisCmd = &cobra.Command{
	Use:   "genesis",
	Short: "Create genesis block on empty blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		bc := NewBlockchain()
		logger.Print("generating genesis block")
		err := bc.AddGenesisBlock(cmd.Context())
		logger.Print("genesis block generated")
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(genesisCmd)
}
