package main

import (
	"github.com/openexw/orpc/server"
	"github.com/openexw/orpc/testdata"
	"log"
)

func main() {
	// 实例化 server
	s := server.NewServer()
	// 注册 service
	s.Register(new(testdata.Sum))
	s.Register(new(testdata.Profile))
	// 启动 server
	err := s.Server("tcp", ":8091")
	if err != nil {
		log.Fatalln("run orpc err:", err)
	}
}
