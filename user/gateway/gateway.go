package gateway

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type ServerConfig []struct {
	name        string
	addr        string
	registerFun func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}

func EnableGateway(c context.Context, gatewayAddr string, serverConfig ServerConfig) {
	c, cancel := context.WithCancel(c)
	defer cancel()
	mux := runtime.NewServeMux()
	for _, s := range serverConfig {
		err := s.registerFun(
			c, mux, s.addr,
			[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		)
		if err != nil {
			log.Printf("cannot register service %s: %v", s.name, err)
		}
	}
	log.Printf("grpc gateway started at %s", gatewayAddr)
	err := http.ListenAndServe(gatewayAddr, mux)
	if err != nil {
		log.Println("listen and serve error")
		return
	}
}
