package grpcutil_test

import (
	"context"
	"github.com/autom8ter/grpcutil"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/mwitkow/go-proto-validators"
	"google.golang.org/grpc"
	"testing"
)

func Test(t *testing.T) {
	_, err := grpcutil.NewServer(
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
		grpcutil.WithProxyServiceRegistration(func(ctx context.Context, mux *runtime.ServeMux, port string, opts ...grpc.DialOption) {

		}),
		grpcutil.WithGRPCServiceRegistration(func(server *grpc.Server) {

		}))
	if err != nil {
		t.Fatal(err.Error())
	}
}
