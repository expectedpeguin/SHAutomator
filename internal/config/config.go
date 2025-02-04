package config

import (
	"flag"
	"runtime"
)

type Config struct {
	Host          string
	Port          int
	Username      string
	Password      string
	KeyFile       string
	ScriptFile    string
	ServersFile   string
	MaxConcurrent int
}

func ParseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Host, "host", "", "SSH host")
	flag.IntVar(&cfg.Port, "port", 22, "SSH port")
	flag.StringVar(&cfg.Username, "username", "", "SSH username")
	flag.StringVar(&cfg.Password, "password", "", "SSH password")
	flag.StringVar(&cfg.KeyFile, "keyfile", "", "Path to private key file")
	flag.StringVar(&cfg.ScriptFile, "script", "", "Path of the script to run")
	flag.StringVar(&cfg.ServersFile, "servers", "", "File containing server details")
	flag.IntVar(&cfg.MaxConcurrent, "concurrent", runtime.NumCPU(), "Maximum concurrent connections")

	flag.Parse()
	return cfg
}
