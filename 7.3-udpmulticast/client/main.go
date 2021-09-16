package main

import (
	"fmt"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "224.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	defer conn.Close()
	buff := make([]byte, 1500)
	for {
		length, remote, err := conn.ReadFromUDP(buff)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%q: %s\n", remote, string(buff[:length]))
	}
}
