package grpc_helper

import "google.golang.org/grpc"

type GRPCRequest struct {
	Address string
}

func (g *GRPCRequest) CallInsecureBlocking(call func(*grpc.ClientConn) (interface{}, error)) (interface{}, error) {
	conn, err := grpc.Dial(g.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return call(conn)
}
