package main

import (
	"context"
	"github.com/openexw/orpc/client"
	"net"
	"sync"
)

type Args struct {
	A, B int
}

func main() {
	conn, err := net.Dial("tcp", ":8091")
	if err != nil {
		return
	}

	cli := client.NewClient(conn, client.WithIsTrace(true))

	wg := sync.WaitGroup{}
	//wg.Add(5)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			args := &Args{
				A: index,
				B: 2,
			}
			var reply string
			cli.Call(context.Background(), "Foo.Sum", args, &reply)
			println("A+B=", reply)
		}(i)
	}
	wg.Wait()
}
