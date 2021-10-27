package main

import (
	"github.com/openexw/orpc/server"
	"github.com/openexw/orpc/testdata"
	"log"
)

func main() {
	s := server.NewServer()
	s.Register(new(testdata.Sum))
	s.Register(new(testdata.Profile))
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
