package cmd

import (
	"fmt"
	"log"

	"github.com/philip-21/philfort/internal/ssh"
	"github.com/spf13/cobra"
)

var (
	downloadUser      string
	downloadHost      string
	downloadKeyPath   string
	downloadLocalPath string
	downloadRemotePath string
)

// downloadCmd handles downloading files from the remote server
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from a remote machine via SSH",
	Run: func(cmd *cobra.Command, args []string) {
		if downloadUser == "" || downloadHost == "" || downloadKeyPath == "" || downloadLocalPath == "" || downloadRemotePath == "" {
			log.Fatalf("PhilFort: missing required flags (user, host, key, local, remote)")
		}

		client, err := ssh.Connect(downloadUser, downloadHost, downloadKeyPath)
		if err != nil {
			log.Fatalf("PhilFort: failed to connect to server: %v", err)
		}
		defer client.Close()

		fmt.Printf("PhilFort: Downloading %s@%s:%s → %s\n", downloadUser, downloadHost, downloadRemotePath, downloadLocalPath)

		if err := ssh.DownloadFile(client, downloadRemotePath, downloadLocalPath); err != nil {
			log.Fatalf("PhilFort: download failed: %v", err)
		}

		fmt.Println("PhilFort: Download complete ✅")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&downloadUser, "user", "u", "", "SSH username (e.g. ubuntu)")
	downloadCmd.Flags().StringVarP(&downloadHost, "host", "H", "", "SSH host (e.g. 1.2.3.4)")
	downloadCmd.Flags().StringVarP(&downloadKeyPath, "key", "k", "", "Path to private key (e.g. ~/.ssh/philfort)")
	downloadCmd.Flags().StringVarP(&downloadLocalPath, "local", "l", "", "Local path to save file")
	downloadCmd.Flags().StringVarP(&downloadRemotePath, "remote", "r", "", "Remote file path to download")
}

//philfort download -u ubuntu -H 192.168.1.5 -k ~/.ssh/philfort -r /home/ubuntu/data.txt -l ./data.txt
