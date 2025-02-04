package sshhandler

import (
	"context"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	client *ssh.Client
	config *ssh.ClientConfig
}

func NewClient(server ServerDetails) (*Client, error) {
	auth, err := GetAuthMethod(server.Password, server.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth method: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            server.Username,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &Client{
		client: client,
		config: config,
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) ExecuteCommands(ctx context.Context, commands []string) error {
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Setup pipes for stdin, stdout, stderr
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	// Handle command execution
	errChan := make(chan error, 1)
	go func() {
		for _, cmd := range commands {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				if _, err := fmt.Fprintln(stdin, cmd); err != nil {
					errChan <- err
					return
				}
			}
		}
		stdin.Close()
		errChan <- session.Wait()
	}()

	// Handle output
	go io.Copy(io.Discard, stdout)
	go io.Copy(io.Discard, stderr)

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
