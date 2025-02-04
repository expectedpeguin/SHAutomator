package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"SSHAutomator/internal/config"
	"SSHAutomator/internal/executor"
	"SSHAutomator/internal/sshhandler"
)

func main() {
	cfg := config.ParseFlags()

	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		cancel()
	}()

	// Read commands
	commands, err := sshhandler.ReadScriptFile(cfg.ScriptFile)
	if err != nil {
		log.Fatalf("Error reading script file: %v", err)
	}

	// Get server details
	servers, err := getServers(cfg)
	if err != nil {
		log.Fatalf("Error getting server details: %v", err)
	}

	// Execute commands
	if err := executeCommands(ctx, cfg, servers, commands); err != nil {
		log.Fatalf("Error executing commands: %v", err)
	}
}

func getServers(cfg *config.Config) ([]sshhandler.ServerDetails, error) {
	if cfg.ServersFile != "" {
		return sshhandler.ReadServersFile(cfg.ServersFile)
	}

	if cfg.Host == "" || cfg.Username == "" {
		return nil, fmt.Errorf("host and username are required when not using servers file")
	}

	return []sshhandler.ServerDetails{{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		KeyFile:  cfg.KeyFile,
	}}, nil
}

func executeCommands(ctx context.Context, cfg *config.Config, servers []sshhandler.ServerDetails, commands []string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(servers))

	// Create execution pool
	pool := executor.NewPool(ctx, cfg.MaxConcurrent)

	// Execute commands for each server
	for _, server := range servers {
		wg.Add(1)
		server := server // Create new variable for goroutine

		pool.Submit(func() {
			defer wg.Done()

			client, err := sshhandler.NewClient(server)
			if err != nil {
				errChan <- fmt.Errorf("failed to create SSH client for %s: %v", server.Host, err)
				return
			}
			defer client.Close()

			if err := client.ExecuteCommands(ctx, commands); err != nil {
				errChan <- fmt.Errorf("command execution failed on %s: %v", server.Host, err)
			}
		})
	}

	// Wait for all executions to complete
	wg.Wait()
	close(errChan)

	// Collect errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors during execution", len(errors))
	}

	return nil
}
