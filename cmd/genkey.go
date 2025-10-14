package cmd

//private and public key generation

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/philip-21/philfort/internal/ssh"
	"github.com/spf13/cobra"
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

//philfort gen-key -o ~/.ssh/philfort -p ~/.ssh/philfort.pub
//philfort gen-key --output ~/.ssh/philfort --pub ~/.ssh/philfort.pub
