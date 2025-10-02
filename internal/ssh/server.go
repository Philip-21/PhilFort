package ssh

// Package ssh server provides the server-side implementation for PhilFort.
//
// The SSH server runs on the remote machine (e.g., a VM) and is responsible for:
//   - Authenticating clients using their private/public key pairs:
//       * The client signs the connection with its private key
//       * The server verifies the signature against the authorized public keys
//   - Accepting and executing remote shell commands
//   - Supporting secure file transfer through SFTP
//
// Note: This is the server-side only. The client (running on a laptop or other machine) 
// uses its private key to connect to the server. The server checks if the matching 
// public key is present in its authorized_keys file before granting access.


import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/sftp"
	gossh "golang.org/x/crypto/ssh"
)

var (
	mu             sync.RWMutex
	authorizedKeys = make(map[string]gossh.PublicKey) // global key store
)

// loadPrivateKey loads the server host private key from disk
func loadPrivateKey(path string) (gossh.Signer, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	signer, err := gossh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return signer, nil
}

// loadAuthorizedKeys loads all authorized client public keys into memory.
// Supports OpenSSH-style "authorized_keys" with multiple lines.
func loadAuthorizedKeys(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open authorized keys: %w", err)
	}
	defer file.Close()

	mu.Lock()
	defer mu.Unlock()
	authorizedKeys = make(map[string]gossh.PublicKey)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // skip empty lines and comments
		}
		pubKey, _, _, _, err := gossh.ParseAuthorizedKey([]byte(line))
		if err != nil {
			log.Printf("Skipping invalid key: %v", err)
			continue
		}
		authorizedKeys[string(pubKey.Marshal())] = pubKey
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading authorized keys: %w", err)
	}

	log.Printf("Loaded %d authorized keys", len(authorizedKeys))
	return nil
}

// reloadAuthorizedKeys allows reloading keys at runtime (for revocation/updating).
func reloadAuthorizedKeys(path string) error {
	return loadAuthorizedKeys(path)
}

// isAuthorized checks if a given public key is in the authorized list.
func isAuthorized(key gossh.PublicKey) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := authorizedKeys[string(key.Marshal())]
	return ok
}

// StartServer boots an SSH server with key authentication, command exec, and SFTP
func StartServer(addr, hostKeyPath, authorizedKeyPath string) error {
	// Load persistent host key
	hostKey, err := loadPrivateKey(hostKeyPath)
	if err != nil {
		return fmt.Errorf("could not load host key: %w", err)
	}

	// Load authorized keys into memory
	if err := loadAuthorizedKeys(authorizedKeyPath); err != nil {
		return fmt.Errorf("could not load authorized keys: %w", err)
	}

	// Configure server
	config := &gossh.ServerConfig{
		PublicKeyCallback: func(conn gossh.ConnMetadata, key gossh.PublicKey) (*gossh.Permissions, error) {
			if !isAuthorized(key) {
				return nil, fmt.Errorf("unauthorized key for user %s", conn.User())
			}

			// OPTIONAL: Enforce username → key mapping (per-user authorized_keys)
			userKeysFile := filepath.Join(filepath.Dir(authorizedKeyPath), conn.User()+".pub")
			if _, err := os.Stat(userKeysFile); err == nil {
				// if per-user file exists, enforce user check
				userKeys := make(map[string]gossh.PublicKey)
				if err := loadUserKeys(userKeysFile, userKeys); err == nil {
					if _, ok := userKeys[string(key.Marshal())]; !ok {
						return nil, fmt.Errorf("key valid but not authorized for user %s", conn.User())
					}
				}
			}

			log.Printf("Login successful for user %s from %s", conn.User(), conn.RemoteAddr())
			return nil, nil
		},
	}
	config.AddHostKey(hostKey)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	log.Printf("SSH server listening on %s", addr)

	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept incoming connection: %v", err)
			continue
		}

		// Perform SSH handshake
		sshConn, chans, reqs, err := gossh.NewServerConn(nConn, config)
		if err != nil {
			log.Printf("failed handshake: %v", err)
			continue
		}
		log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

		go gossh.DiscardRequests(reqs)
		go handleChannels(chans, sshConn)
	}
}

// loadUserKeys loads per-user key file (optional enhancement).
func loadUserKeys(path string, store map[string]gossh.PublicKey) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	pubKey, _, _, _, err := gossh.ParseAuthorizedKey(data)
	if err != nil {
		return err
	}
	store[string(pubKey.Marshal())] = pubKey
	return nil
}

func handleChannels(chans <-chan gossh.NewChannel, conn *gossh.ServerConn) {
	for newChannel := range chans {
		switch newChannel.ChannelType() {
		case "session":
			channel, requests, err := newChannel.Accept()
			if err != nil {
				log.Printf("could not accept channel: %v", err)
				continue
			}
			go handleSession(channel, requests)

		case "subsystem":
			channel, requests, err := newChannel.Accept()
			if err != nil {
				log.Printf("could not accept subsystem channel: %v", err)
				continue
			}
			go func(in <-chan *gossh.Request) {
				for req := range in {
					if req.Type == "subsystem" && string(req.Payload[4:]) == "sftp" {
						log.Printf("Starting SFTP session for %s", conn.User())
						req.Reply(true, nil)

						server, err := sftp.NewServer(channel)
						if err != nil {
							log.Printf("SFTP server error: %v", err)
							channel.Close()
							return
						}
						if err := server.Serve(); err == io.EOF {
							server.Close()
							log.Printf("SFTP session closed")
						} else if err != nil {
							log.Printf("SFTP server completed with error: %v", err)
						}
					} else {
						req.Reply(false, nil)
					}
				}
			}(requests)

		default:
			newChannel.Reject(gossh.UnknownChannelType, "unsupported channel type")
		}
	}
}

func handleSession(channel gossh.Channel, requests <-chan *gossh.Request) {
	for req := range requests {
		switch req.Type {
		case "exec":
			cmd := string(req.Payload[4:])
			log.Printf("Executing command for client: %s", cmd)

			// SECURITY: server process user executes this command
			// Best practice: run this server under an unprivileged account.
			command := exec.Command("sh", "-c", cmd)
			command.Stdout = channel
			command.Stderr = channel

			if err := command.Run(); err != nil {
				io.WriteString(channel, fmt.Sprintf("Error: %s\n", err))
			}

			req.Reply(true, nil)
			channel.Close()
		default:
			req.Reply(false, nil)
		}
	}
}
