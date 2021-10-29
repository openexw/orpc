package main

import (
	"context"
	"fmt"
	"github.com/openexw/orpc/client"
	"github.com/openexw/orpc/internal/protocol"
	"github.com/openexw/orpc/testdata"
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

	cli := client.NewClient(conn, client.WithIsTrace(false), client.WithSerializeType(protocol.JSON))

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			profileArgs := &testdata.Profile{
				Name: "Jack",
				Age:  uint8(18 + index),
				Sex:  1,
			}
			var profileReply testdata.Profile
			cli.Call(context.Background(), "Profile.AddProfile", profileArgs, &profileReply)
			fmt.Printf("Profile.AddProfile##argv:%+v, resp: %+v\n", profileArgs, &profileReply)
		}(i)
	}

	wg.Wait()
}
