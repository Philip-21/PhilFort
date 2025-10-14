package cmd

import (
	"fmt"
	"log"

	"github.com/philip-21/philfort/internal/ssh"
	"github.com/spf13/cobra"
)

var (
	uploadUser       string
	uploadHost       string
	uploadKeyPath    string
	uploadLocalPath  string
	uploadRemotePath string
)

// uploadCmd handles uploading files to the remote server
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to a remote machine via SSH",
	Run: func(cmd *cobra.Command, args []string) {
		if uploadUser == "" || uploadHost == "" || uploadKeyPath == "" || uploadLocalPath == "" || uploadRemotePath == "" {
			log.Fatalf("PhilFort: missing required flags (user, host, key, local, remote)")
		}

		client, err := ssh.Connect(uploadUser,uploadHost,uploadKeyPath)
		if err != nil {
			log.Fatalf("PhilFort: failed to connect to server: %v", err)
		}
		defer client.Close()

		fmt.Printf("PhilFort: Uploading %s → %s@%s:%s\n", uploadLocalPath, uploadUser, uploadHost, uploadRemotePath)

		if err := ssh.UploadFile(client, uploadLocalPath, uploadRemotePath); err != nil {
			log.Fatalf("PhilFort: upload failed: %v", err)
		}

		fmt.Println("PhilFort: Upload complete ✅")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVarP(&uploadUser, "user", "u", "", "SSH username (e.g. ubuntu)")
	uploadCmd.Flags().StringVarP(&uploadHost, "host", "H", "", "SSH host (e.g. 1.2.3.4)")
	uploadCmd.Flags().StringVarP(&uploadKeyPath, "key", "k", "", "Path to private key (e.g. ~/.ssh/philfort)")
	uploadCmd.Flags().StringVarP(&uploadLocalPath, "local", "l", "", "Local file path to upload")
	uploadCmd.Flags().StringVarP(&uploadRemotePath, "remote", "r", "", "Remote destination path")
}

//philfort upload -u ubuntu -H 192.168.1.5 -k ~/.ssh/philfort -l ./data.txt -r /home/ubuntu/data.txt
