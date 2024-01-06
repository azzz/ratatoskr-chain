/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/azzz/ratatoskr/pkg/proofofwork"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all blocksh",
	RunE: func(cmd *cobra.Command, args []string) error {
		pow := proofofwork.New()
		bc := LoadBlockchain()
		iter := bc.Iterator()

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(w, "timestamp\thash\tvalid\tvalue\n")

		for iter.Next() {
			if iter.Err() != nil {
				return iter.Err()
			}
			b := iter.Block()
			ts := time.UnixMicro(int64(b.Timestamp))
			fmt.Fprintf(w,
				"%s\t%s\t%t\t%q\n",
				ts.Format(time.DateTime),
				b.String(),
				pow.Validate(b),
				b.Data,
			)
		}

		w.Flush()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
