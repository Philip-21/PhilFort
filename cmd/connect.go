package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/philip-21/philfort/internal/ssh"
	"github.com/spf13/cobra"
)

var (
	host    string
	user    string
	keyPath string
	command string
)

// connectCmd establishes an SSH connection to a remote host and optionally executes a command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect securely to a remote PhilFort SSH server",
	Long: `Connect to a remote server using key-based authentication. 
Optionally execute a command using the --cmd flag.
Example:
  philfort connect --host 192.168.1.20 --user ubuntu --key ~/.ssh/philfort --cmd "ls -la"`,
	Run: func(cmd *cobra.Command, args []string) {
		if host == "" || user == "" || keyPath == "" {
			log.Fatalf("PhilFort: --host, --user, and --key are required")
		}

		client, session, err := ssh.ConnectWithSession(host, user, keyPath)
		if err != nil {
			log.Fatalf("PhilFort: failed to connect to server: %v", err)
		}
		defer client.Close()
		defer session.Close()

		if command != "" {
			output, err := session.CombinedOutput(command)
			if err != nil {
				log.Fatalf("PhilFort: command failed: %v", err)
			}
			fmt.Println(string(output))
		} else {
			// Interactive shell
			session.Stdout = os.Stdout
			session.Stderr = os.Stderr
			session.Stdin = os.Stdin

			modes := ssh.DefaultTerminalModes()
			if err := ssh.RequestPty(session, modes); err != nil {
				log.Fatalf("PhilFort: PTY request failed: %v", err)
			}

			if err := session.Shell(); err != nil {
				log.Fatalf("PhilFort: unable to start shell: %v", err)
			}

			session.Wait()
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringVar(&host, "host", "", "Target host (e.g. 192.168.1.20:22)")
	connectCmd.Flags().StringVar(&user, "user", "", "Username for SSH authentication")
	connectCmd.Flags().StringVar(&keyPath, "key", "", "Path to private key for authentication")
	connectCmd.Flags().StringVar(&command, "cmd", "", "Optional command to execute remotely")
}

//interactive connection
//philfort connect --host 192.168.1.10:22 --user ubuntu --key ~/.ssh/philfort

//single command
//philfort connect --host 192.168.1.10:22 --user ubuntu --key ~/.ssh/philfort --cmd "uptime"
