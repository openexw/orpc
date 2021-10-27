package testdata

type Sum int
type Args struct {
	A int
	B int
}

func (s *Sum) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}
