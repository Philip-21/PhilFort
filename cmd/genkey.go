package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"philfort/internal/ssh"
)

var privateKeyPath string
var publicKeyPath string

var genKeyCmd = &cobra.Command{
	Use:   "gen-key",
	Short: "Generate SSH key pair for PhilFort",
	Run: func(cmd *cobra.Command, args []string) {
		// Expand "~" in paths
		privateKeyPathExpanded, _ := filepath.Abs(privateKeyPath)
		publicKeyPathExpanded, _ := filepath.Abs(publicKeyPath)

		err := ssh.GenerateKeyPair(privateKeyPathExpanded, publicKeyPathExpanded, 4096)
		if err != nil {
			log.Fatalf("PhilFort: Failed to generate keys: %v", err)
		}

		fmt.Printf("PhilFort: Keys successfully generated!\n")
		fmt.Printf("Private key: %s\n", privateKeyPathExpanded)
		fmt.Printf("Public key:  %s\n", publicKeyPathExpanded)
	},
}

func init() {
	rootCmd.AddCommand(genKeyCmd)

	genKeyCmd.Flags().StringVarP(&privateKeyPath, "output", "o", "./philfort", "Private key output path")
	genKeyCmd.Flags().StringVarP(&publicKeyPath, "pub", "p", "./philfort.pub", "Public key output path")
}
