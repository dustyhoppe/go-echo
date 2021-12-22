package main

import (
	"github.com/dustyhoppe/go-echo/config"
	"github.com/dustyhoppe/go-echo/proxy"
)

func main() {

	config := config.NewDefaultConfig()

	proxyServer := proxy.NewProxyServer(config, nil)

	proxyServer.Start()
}
