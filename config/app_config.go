package config

import (
	"flag"
	"log"
	"strings"
)

type AppConfig struct {
	Protocol string
	PemFile  string
	KeyFile  string
	Port     int
}

func NewDefaultConfig() *AppConfig {
	var pemPath, keyPath, proto string
	var port int

	flag.StringVar(&pemPath, "pem", "server.pem", "path to pem file")
	flag.StringVar(&keyPath, "key", "server.key", "path to key file")
	flag.StringVar(&proto, "proto", "http", "Proxy protocol (http or https)")
	flag.IntVar(&port, "port", 8888, "Port to listen on")

	flag.Parse()

	if !strings.EqualFold(proto, "http") && !strings.EqualFold(proto, "https") {
		log.Fatal("Protocol must be either http or https")
	}

	return &AppConfig{
		PemFile:  pemPath,
		KeyFile:  keyPath,
		Protocol: proto,
		Port:     port,
	}
}
