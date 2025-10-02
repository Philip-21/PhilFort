package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "philfort",
	Short: "PhilFort CLI - SSH + SFTP automation",
	Long:  "PhilFort CLI allows you to generate keys, connect to a VM or Server, and transfer files securely.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
