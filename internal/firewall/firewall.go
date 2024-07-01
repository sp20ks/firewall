package firewall

import (
	"fmt"
	"net"
	"strings"
)

func FilterConnection(conn net.Conn, rules []string) {
	remotrAddr := conn.RemoteAddr().String()
	fmt.Println(remotrAddr)
	allowed := false

	for _, rule := range rules {
		fmt.Println(rule)
		if strings.Contains(remotrAddr, rule) {
			allowed = true
			break
		}
	}

	if allowed {
		fmt.Printf("Allowed connection from %s\n", remotrAddr)
	} else {
		fmt.Printf("Blocked conntection from %s\n", remotrAddr)
		conn.Close()
	}
}
