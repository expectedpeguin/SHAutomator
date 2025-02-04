# SSHAutomator

SSHAutomator is a high-performance Go program that provides concurrent execution of scripts and commands on multiple remote servers over SSH, with built-in connection pooling and error handling.

## Features

- **Concurrent Execution**: Execute commands on multiple servers simultaneously with controlled concurrency
- **Connection Pooling**: Smart connection management to prevent resource exhaustion
- **Robust Error Handling**: Comprehensive error reporting and graceful failure management
- **Context-Aware**: Support for timeouts and graceful cancellation
- **Flexible Authentication**: Support for both password and key-based authentication
- **Resource Efficient**: Optimized memory usage and connection handling
- **Configuration Options**: Extensive command-line options for customization

## Installation

```bash
go install github.com/yourusername/SSHAutomator@latest
```

## Usage

The program accepts the following command-line arguments:

```bash
sshautomator -script <script.txt> [options]
```

### Options

- `-host`: SSH host to connect to
- `-port`: SSH port (default: 22)
- `-username`: SSH username
- `-password`: SSH password (optional if using key file)
- `-keyfile`: Path to private key file (optional)
- `-script`: Path to the script file containing commands
- `-servers`: Path to server list file
- `-concurrent`: Maximum number of concurrent connections (default: number of CPU cores)

### Authentication Methods

You can authenticate using either a password or a private key:

```bash
# Password authentication
sshautomator -script commands.txt -host example.com -username user -password pass

# Key-based authentication
sshautomator -script commands.txt -host example.com -username user -keyfile ~/.ssh/id_rsa
```

### Server List Format

When using the `-servers` option, the server list file should be formatted as follows:

```text
host1 username password
host2 username keyfile /path/to/key
host3 username password
```

Each line represents one server with space-separated fields:
1. Host address
2. Username
3. Authentication type ("password" or "keyfile")
4. Password or path to keyfile (depending on authentication type)

### Example Usage

Single server execution:
```bash
sshautomator -script script.txt -host example.com -port 22 -username myuser -password mypassword
```

Multiple servers with concurrent execution:
```bash
sshautomator -script script.txt -servers serverlist.txt -concurrent 5
```

## Script File Format

The script file should contain one command per line:

```bash
cd /var/log
ls -la
grep "error" application.log
```

## Error Handling

The program provides detailed error reporting:
- Connection failures
- Authentication errors
- Command execution failures
- Timeout issues

Error messages include the affected server and specific failure reason.

## Performance Considerations

- Uses connection pooling to manage resources efficiently
- Controlled concurrency to prevent overwhelming target servers
- Efficient memory usage with proper cleanup
- Context-aware execution with timeout support

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

1. Clone the repository
```bash
git clone https://github.com/yourusername/SSHAutomator.git
```

2. Install dependencies
```bash
go mod tidy
```

3. Run tests
```bash
go test ./...
```

### Code Style

- Follow standard Go coding conventions
- Use `gofmt` for code formatting
- Include tests for new features
- Update documentation as needed

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Security Considerations

- The program uses `ssh.InsecureIgnoreHostKey()` for host key verification. In production environments, consider implementing proper host key verification.
- Passwords and sensitive data are handled securely in memory.
- Connection timeouts are implemented to prevent hanging connections.
- Consider using key-based authentication instead of passwords when possible.

## Troubleshooting

Common issues and solutions:

1. Connection timeouts
   - Check network connectivity
   - Verify SSH service is running on target
   - Check firewall settings

2. Authentication failures
   - Verify credentials
   - Check key file permissions
   - Ensure key format is correct

3. Execution errors
   - Verify command syntax
   - Check user permissions on target
   - Review server logs for details