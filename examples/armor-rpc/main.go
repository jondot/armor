package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jondot/armor"
	pb "github.com/jondot/armor/examples/armor-rpc/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

type HelloServer struct {
	armor *armor.Armor
}

func (s *HelloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	s.armor.Metrics.Inc("rpc.say_hello")
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	m := armor.New("rpc", "1.0")
	r := m.GinRouter()

	r.GET("/", func(c *gin.Context) {
		defer m.Metrics.Timed("timed.request", time.Now())
		m.Metrics.Inc("foobar")
		c.String(200, "hello world")
	})

	m.RunWithRPC(r,
		func(s *grpc.Server) {
			pb.RegisterGreeterServer(s, &HelloServer{m})
		})
}
