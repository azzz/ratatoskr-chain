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
		addrs, err := cmd.Flags().GetString("address")
		cobra.CheckErr(err)

		bc := CreateBlockchain(addrs)

		logger.Printf("done, tip: %x", bc.Tip())
	},
}

func init() {
	createCmd.Flags().String("address", "", "Blockchain address")
	createCmd.MarkFlagRequired("address")

	rootCmd.AddCommand(createCmd)
}
