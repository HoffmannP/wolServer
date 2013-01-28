package networking

import (
	"net"
)

const wol_port = "9"

func SendMagicPacket(MAC string, broadcast string) error {
	broadcast = broadcast + ":" + wol_port
	packet, err := constructMagicPacket(MAC)
	if err != nil {
		return err
	}

	c, err := net.Dial("udp4", broadcast)
	if err != nil {
		return err
	}
	defer c.Close()

	written, err := c.Write(packet)
	// println("Sending magic packet to", broadcast, "with", MAC)
	if err != nil {
		return err
	}
	if written < len(packet) {
		return net.UnknownNetworkError("package not completely send")
	}
	return nil
}

func constructMagicPacket(MAC string) ([]byte, error) {
	macBytes, err := net.ParseMAC(MAC)
	if err != nil {
		return nil, err
	}

	b := []uint8{255, 255, 255, 255, 255, 255}
	for i := 0; i < 16; i++ {
		b = append(b, macBytes...)
	}
	return b, err
}
