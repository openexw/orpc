package main

import (
	"context"
	"fmt"
	"github.com/openexw/orpc/client"
	"github.com/openexw/orpc/testdata"
	"net"
)

type Args struct {
	A, B int
}

func main() {
	conn, err := net.Dial("tcp", ":8091")
	if err != nil {
		return
	}

	cli := client.NewClient(conn, client.WithIsTrace(false))
	args := &Args{
		A: 12,
		B: 2,
	}
	var reply int
	cli.Call(context.Background(), "Sum.Add", args, &reply)
	println(args.A, "+", args.B, "=", reply)

	profileArgs := &testdata.Profile{
		Name: "Jack",
		Age:  18,
		Sex:  1,
	}
	var profileReply testdata.Profile
	cli.Call(context.Background(), "Profile.AddProfile", profileArgs, &profileReply)
	fmt.Printf("Profile.AddProfile##argv:%+v, resp: %+v", profileArgs, &profileReply)
}
