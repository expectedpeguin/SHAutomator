package sshhandler

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

func ReadScriptFile(scriptFile string) ([]string, error) {
	if scriptFile == "" {
		return nil, nil
	}

	content, err := ioutil.ReadFile(scriptFile)
	if err != nil {
		return nil, err
	}

	scriptLines := strings.FieldsFunc(string(content), func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	var commands []string
	for _, line := range scriptLines {
		if line != "" {
			commands = append(commands, line)
		}
	}
	return commands, nil
}

func ReadServersFile(filename string) ([]ServerDetails, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var servers []ServerDetails
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 3 {
			server := ServerDetails{
				Host:     fields[0],
				Username: fields[1],
				Password: fields[2],
			}
			if fields[2] == "keyfile" && len(fields) == 4 {
				server.KeyFile = fields[3]
			}
			servers = append(servers, server)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return servers, nil
}
