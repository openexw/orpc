package main

import (
	"github.com/openexw/orpc/server"
	"log"
)

func main() {
	s := server.NewServer()
	err := s.Server("tcp", ":8091")
	if err != nil {
		log.Fatalln("run orpc err:", err)
	}

	//reader := strings.NewReader("Hello1212")
	//type header [12]byte
	//
	//h := new(header)
	//io.ReadFull(reader, []byte(h))
	//fmt.Printf("%s", h)
}
