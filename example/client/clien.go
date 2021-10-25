package main

import "net"

func main() {
	_, err := net.Dial("tcp", ":8091")
	if err != nil {
		return
	}
}
