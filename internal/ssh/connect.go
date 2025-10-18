package ssh

import (
	"fmt"
	"os"

	gossh "golang.org/x/crypto/ssh"
)

// Connect creates an SSH client and session using key-based authentication
func ConnectWithSession(host, user, keyPath string) (*gossh.Client, *gossh.Session, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read private key: %v", err)
	}

	signer, err := gossh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	config := &gossh.ClientConfig{
		User:            user,
		Auth:            []gossh.AuthMethod{gossh.PublicKeys(signer)},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(), // TODO: Replace with secure verification
	}

	client, err := gossh.Dial("tcp", host, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial SSH: %v", err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("failed to create SSH session: %v", err)
	}

	return client, session, nil
}

// DefaultTerminalModes defines a sane terminal mode setup for interactive shells
func DefaultTerminalModes() gossh.TerminalModes {
	return gossh.TerminalModes{
		gossh.ECHO:          1,     // enable echoing
		gossh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		gossh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
}

// RequestPty requests a pseudo-terminal for the SSH session
func RequestPty(session *gossh.Session, modes gossh.TerminalModes) error {
	return session.RequestPty("xterm", 80, 40, modes)
}
