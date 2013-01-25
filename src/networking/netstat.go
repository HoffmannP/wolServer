package networking

import (
	"net"
	"time"
)

const stat_port = "139" // all computers reply on port 139 (Samba/Windos File Share)

func Netstat(IP string) bool {
	protocol := "tcp"
	_, err := net.DialTimeout(protocol, IP+":"+stat_port, 1*time.Second)
	if err != nil {
		return false
	}
	return true
}
