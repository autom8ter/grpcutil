package main

import (
	"context"
	grpc_util "github.com/autom8ter/grpc-util"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
)

func main() {
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
		grpc_util.WithProxyServiceRegistration(func(ctx context.Context, mux *runtime.ServeMux, port string, opts ...grpc.DialOption) {

		}),
		grpc_util.WithGRPCServiceRegistration(func(server *grpc.Server) {

		}))
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := server.Run(context.Background()); err != nil {
		log.Fatal(err.Error())
	}
}
