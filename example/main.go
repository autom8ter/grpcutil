package main

import (
	"context"
	"fmt"
	grpc_util "github.com/autom8ter/grpc-util"
	"github.com/autom8ter/grpc-util/example/gen/go/example"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
)

//Hello is a HelloServiceServer
type Hello struct{}

//Implement hello api server interface
func (h Hello) Hello(ctx context.Context, request *example.HelloRequest) (*example.HelloResponse, error) {
	return &example.HelloResponse{
		Response: fmt.Sprintf("Greetings! You said: %s", request.Text),
	}, nil
}

func main() {
	hellosvc := &Hello{}
	server, err := grpc_util.NewServer(
		context.Background(),
		grpc_util.WithClientProxyDialOptions(grpc.WithInsecure()),
		grpc_util.WithPort(8080),
		grpc_util.WithServerOptions(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_validator.UnaryServerInterceptor(),
				grpc_recovery.UnaryServerInterceptor(),
			)),
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_validator.StreamServerInterceptor(),
				grpc_recovery.StreamServerInterceptor(),
			)),
		),
		grpc_util.WithProxyOptions(runtime.WithIncomingHeaderMatcher(grpc_util.MappedHeaderMatcherFunc(map[string]bool{
			"authorization": true,
			"Authorization": true,
		}))),
		grpc_util.WithProxyServiceRegistration(func(ctx context.Context, mux *runtime.ServeMux, host string, opts ...grpc.DialOption) {
			if err := example.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, host, opts); err != nil {
				log.Fatal(err.Error())
			}
		}),
		grpc_util.WithGRPCServiceRegistration(func(server *grpc.Server) {
			example.RegisterHelloServiceServer(server, hellosvc)
		}))
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := server.Run(context.Background()); err != nil {
		log.Fatal(err.Error())
	}
}
