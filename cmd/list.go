/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all blocksh",
	RunE: func(cmd *cobra.Command, args []string) error {
		bc := NewBlockchain()
		iter := bc.Iterator()

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(w, "ts\tnonce\thash\tvalue\n")

		for iter.Next() {
			if iter.Err() != nil {
				return iter.Err()
			}
			b := iter.Block()
			ts := time.UnixMicro(int64(b.Timestamp))
			fmt.Fprintf(w,
				"%s\t%d\t%x\t%s\n",
				ts.Format(time.DateTime),
				b.Nonce,
				b.Hash,
				b.Data)
		}

		w.Flush()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
