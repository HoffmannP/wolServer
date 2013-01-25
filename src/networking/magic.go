package networking

import (
	"net"
)

const wol_port = "7"

func SendMagicPacket(macAddr string, bcastAddr string) error {

	packet, err := constructMagicPacket(macAddr)
	if err != nil {
		return err
	}

	c, err := net.Dial("udp", bcastAddr+":"+wol_port)
	if err != nil {
		return err
	}
	defer c.Close()

	written, err := c.Write(packet)

	if (err != nil) || (written != len(packet)) {
		return err
	}

	return nil
}

func constructMagicPacket(macAddr string) ([]byte, error) {
	macBytes, err := net.ParseMAC(macAddr)
	if err != nil {
		return nil, err
	}

	b := []uint8{255, 255, 255, 255, 255, 255}
	for i := 0; i < 16; i++ {
		b = append(b, macBytes...)
	}
	return b, err
}
