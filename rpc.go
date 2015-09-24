package armor

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type RPC struct {
	config *Config
}

func newRPC(config *Config) *RPC {
	return &RPC{config: config}
}

func (r *RPC) Run(rpcInit func(*grpc.Server)) {
	runOn := fmt.Sprintf("%s:%s",
		r.config.GetString("rpc.interface"),
		r.config.GetStringWithDefault("rpc.port", "46060"),
	)

	lis, err := net.Listen("tcp", runOn)
	if err != nil {
		log.Fatalf("-X gRPC failed to listen: %v", err)
	}

	s := grpc.NewServer()
	rpcInit(s)
	go s.Serve(lis)
	log.Printf("-> gRPC running on %s", runOn)
}
