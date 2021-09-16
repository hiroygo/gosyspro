package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("udp4", "224.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	interval := 5 * time.Second
	for now := range time.Tick(interval) {
		if _, err := conn.Write([]byte(now.String())); err != nil {
			panic(err)
		}
		fmt.Println("Tick: ", now.String())
	}
}
