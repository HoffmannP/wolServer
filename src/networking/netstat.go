package networking

import (
	"net"
	"time"
)

const stat_port = "135" // most windows computers should reply on port 135

func Netstats(IP string, ports []string) (bool, error) {
	for _, port := range ports {
		if status, err := Netstat(IP, port); err != nil {
			return false, err
		} else {
			if status {
				return true, nil
			}
		}
	}
	return false, nil
}

func Netstat(IP, port string) (bool, error) {
	if port == "" {
		port = stat_port
	}
	protocol := "tcp"
	_, err := net.DialTimeout(protocol, IP+":"+port, 1*time.Second)
	if err != nil {
		if e, ok := err.(*net.OpError); ok {
			if e.Timeout() {
				return false, nil
			}
			if e.Err.Error() == "connection refused" {
				return true, nil
			}
			if e.Err.Error() == "no route to host" {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}
