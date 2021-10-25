package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	//s := server.NewServer()
	//s.Server("tcp", ":8091")

	reader := strings.NewReader("Hello1212")
	type header [12]byte

	h := new(header)
	io.ReadFull(reader, []byte(h))
	fmt.Printf("%s", h)
}
