package server

import (
	"firewall/internal/config"
	"firewall/internal/firewall"
	"fmt"
	"log"
	"net"
)

func Start(cfg *config.Config) error {
	addr := cfg.HTTPServer.Address
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}

	defer listener.Close()
	log.Printf("Server listening on address %s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go firewall.FilterConnection(conn, cfg.AllowedIps)
	}
}
