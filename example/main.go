package main

import (
	"context"
	"fmt"
	"github.com/autom8ter/grpcutil"
	"github.com/autom8ter/grpcutil/example/gen/go/example"
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
	server, err := grpcutil.NewServer(
		context.Background(),
		grpcutil.WithClientProxyDialOptions(grpc.WithInsecure()),
		grpcutil.WithPort(8080),
		grpcutil.WithServerOptions(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_validator.UnaryServerInterceptor(),
				grpc_recovery.UnaryServerInterceptor(),
			)),
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_validator.StreamServerInterceptor(),
				grpc_recovery.StreamServerInterceptor(),
			)),
		),
		grpcutil.WithProxyOptions(runtime.WithIncomingHeaderMatcher(grpcutil.MappedHeaderMatcherFunc(map[string]bool{
			"authorization": true,
			"Authorization": true,
		}))),
		grpcutil.WithProxyServiceRegistration(func(ctx context.Context, mux *runtime.ServeMux, host string, opts ...grpc.DialOption) {
			if err := example.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, host, opts); err != nil {
				log.Fatal(err.Error())
			}
		}),
		grpcutil.WithGRPCServiceRegistration(func(server *grpc.Server) {
			example.RegisterHelloServiceServer(server, hellosvc)
		}))
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := server.Run(context.Background()); err != nil {
		log.Fatal(err.Error())
	}
}
