/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the genesis command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Print("creating blockchain")
		bc := CreateBlockchain()
		logger.Printf("done, tip: %x", bc.Tip())
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
