/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// addBlockCmd represents the addBlock command
var addBlockCmd = &cobra.Command{
	Use:   "add-block <string value>",
	Short: "adds a block to the blockchain",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing block value")
		}

		if len(args) > 1 {
			return fmt.Errorf("required 1 argument, give: %d", len(args))
		}

		bc := LoadBlockchain()
		logger.Printf("generating block, tip is: %x", bc.Tip())
		err := bc.AddBlock(args[0])
		if err != nil {
			return err
		}

		logger.Printf("now tip is %x", bc.Tip())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addBlockCmd)

}
