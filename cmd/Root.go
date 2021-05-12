package cmd

import (
	"fmt"
	"os"

	"github.com/sirloso/bcd_dropbox_importer/transferer"
	"github.com/spf13/cobra"
)

var source string
var destination string
var rename bool
var verbose bool

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer from one location to the other",
	Long:  `transfer files from external drive to s drive`,
	Run: func(cmd *cobra.Command, args []string) {
		transferer.Transfer(source, destination, verbose, rename)
	},
}

func Execute() {
	var rootCmd = &cobra.Command{Use: "bcd_mover"}
	rootCmd.AddCommand(transferCmd)

	transferCmd.Flags().StringVarP(&source, "source", "s", "", "Source directory to read from")
	transferCmd.Flags().StringVarP(&destination, "destination", "d", "", "Destination directory to read to")
	transferCmd.Flags().BoolP("rename", "r", true, "Rename files in source directory")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
