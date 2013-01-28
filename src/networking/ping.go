// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/* shameless recycling net/ipraw_test.go */

package networking

import (
	"bytes"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	ICMP4_ECHO_REQUEST = 8
	ICMP4_ECHO_REPLY   = 0
	ICMP4_ECHO_X_CODE  = 0
)

func Ping(IP string) (bool, error) {
	if os.Getuid() != 0 {
		return false, net.UnknownNetworkError("must be root to ping")
	}

	wait := 100 * time.Millisecond
	network := "ip4:icmp"
	seqnum := rand.Int() & (1<<16 - 1)
	id := os.Getpid() & 0xffff
	echo := newICMPv4EchoRequest(id, seqnum, 128, []byte("KT online Ping"))

	c, err := net.ListenPacket(network, "")
	if err != nil {
		return false, err
	}
	c.SetDeadline(time.Now().Add(wait))
	defer c.Close()

	ra, err := net.ResolveIPAddr(network, IP)
	if err != nil {
		return false, err
	}

	_, err = c.WriteTo(echo, ra)
	if err != nil {
		return false, err
	}

	reply := make([]byte, 256)
	for {
		_, _, err = c.ReadFrom(reply)
		if err != nil {
			if e, v := err.(*net.OpError); v && e.Timeout() {
				return false, nil
			}
			return false, err
		}
		rpType, rpCode, rid, rseqnum := parseICMPEchoReply(reply)
		if rpType == ICMP4_ECHO_REPLY &&
			rpCode == ICMP4_ECHO_X_CODE &&
			rid == id &&
			rseqnum == seqnum {
			return true, nil
		}
	}
	panic("You should not bee here!")
}

func newICMPv4EchoRequest(id, seqnum, msglen int, filler []byte) []byte {
	b := newICMPInfoMessage(id, seqnum, msglen, filler)
	b[0] = ICMP4_ECHO_REQUEST
	b[1] = ICMP4_ECHO_X_CODE

	// calculate ICMP checksum
	cklen := len(b)
	s := uint32(0)
	for i := 0; i < cklen-1; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if cklen&1 == 1 {
		s += uint32(b[cklen-1])
	}
	s = (s >> 16) + (s & 0xffff)
	s = s + (s >> 16)
	// place checksum back in header; using ^= avoids the
	// assumption the checksum bytes are zero
	b[2] ^= uint8(^s & 0xff)
	b[3] ^= uint8(^s >> 8)

	return b
}

func newICMPInfoMessage(id, seqnum, msglen int, filler []byte) []byte {
	b := make([]byte, msglen)
	copy(b[8:], bytes.Repeat(filler, (msglen-8)/len(filler)+1))
	b[0] = 0                    // type
	b[1] = 0                    // code
	b[2] = 0                    // checksum
	b[3] = 0                    // checksum
	b[4] = uint8(id >> 8)       // identifier
	b[5] = uint8(id & 0xff)     // identifier
	b[6] = uint8(seqnum >> 8)   // sequence number
	b[7] = uint8(seqnum & 0xff) // sequence number
	return b
}

func parseICMPEchoReply(b []byte) (pType, pCode, id, seqnum int) {
	pType = int(b[0])
	pCode = int(b[1])
	id = int(b[4])<<8 | int(b[5])
	seqnum = int(b[6])<<8 | int(b[7])
	return
}
