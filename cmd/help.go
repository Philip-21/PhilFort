package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Display help information about PhilFort commands",
	Long: `PhilFort CLI - Secure SSH + SFTP automation tool.

Available Commands:
  gen-key     Generate SSH key pair for authentication
  connect     Connect to a remote SSH server
  upload      Upload a file to a remote server
  download    Download a file from a remote server

Usage Examples:
  philfort gen-key --output ~/.ssh/philfort --pub ~/.ssh/philfort.pub
  philfort connect --host 192.168.1.20 --user ubuntu --key ~/.ssh/philfort
  philfort upload --host 192.168.1.20 --src ./local.txt --dst /home/ubuntu/local.txt
  philfort download --host 192.168.1.20 --src /home/ubuntu/app.log --dst ./app.log

Run 'philfort <command> --help' for more information on a specific command.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)
}
