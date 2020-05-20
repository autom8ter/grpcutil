//go:generate godocdown -o README.md

package grpcutil

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

//Config configures a grpc/json server on the same port
type Config struct {
	port                int //port to serve on
	serviceRegistration func(*grpc.Server)
	proxyRegistration   func(ctx context.Context, mux *runtime.ServeMux, host string, opts ...grpc.DialOption) //register json proxies against the servemux(generated by grpc gateway plugin)
	proxyOptions        []runtime.ServeMuxOption                                                               //register header matchers(ex MappedHeaderMatcher)
	proxyClientOpts     []grpc.DialOption                                                                      //register client options between gateway and grpc server
	serverOptions       []grpc.ServerOption                                                                    //register interceptors and other server options
	muxPrefix           string
}

//ConfigOption is a first class function used for grpc/proxy initialization
type ConfigOption func(c *Config)

//WithPort is the port to serve on
func WithPort(port int) ConfigOption {
	return func(c *Config) {
		c.port = port
	}
}

//WithProxyPrefix servers the grpc gateway proxy after this prefix
func WithProxyPrefix(prefix string) ConfigOption {
	return func(c *Config) {
		c.muxPrefix = prefix
	}
}

//WithGRPCServiceRegistration registers service implementations against the input grpc Server (generated by grpc plugin)
func WithGRPCServiceRegistration(fn func(*grpc.Server)) ConfigOption {
	return func(c *Config) {
		c.serviceRegistration = fn
	}
}

//WithProxyServiceRegistration registers json proxies against the gateway runtime servemux(generated by grpc gateway plugin)
func WithProxyServiceRegistration(fn func(ctx context.Context, mux *runtime.ServeMux, port string, opts ...grpc.DialOption)) ConfigOption {
	return func(c *Config) {
		c.proxyRegistration = fn
	}
}

//WithProxyOptions registers header matchers(ex MappedHeaderMatcher)
func WithProxyOptions(opts ...runtime.ServeMuxOption) ConfigOption {
	return func(c *Config) {
		c.proxyOptions = opts
	}
}

//WithClientProxyDialOptions registers client dial options between gateway and grpc server
func WithClientProxyDialOptions(opts ...grpc.DialOption) ConfigOption {
	return func(c *Config) {
		c.proxyClientOpts = opts
	}
}

//WithServerOptions registers grpc server options like unary/stream interceptors
func WithServerOptions(opts ...grpc.ServerOption) ConfigOption {
	return func(c *Config) {
		c.serverOptions = opts
	}
}

//MappedHeaderMatcherFunc is a map from header key - a boolean indicating whether to forward the header to the grpc server
func MappedHeaderMatcherFunc(decider map[string]bool) runtime.HeaderMatcherFunc {
	return func(s2 string) (s string, b bool) {
		allow := decider[s2]
		return s2, allow
	}
}

//AllHeaderMatcherFunc forwards all headers to the grpc server
func AllHeaderMatcherFunc() runtime.HeaderMatcherFunc {
	return func(s2 string) (s string, b bool) {
		return s2, true
	}
}

//Server is a grpc/http server that can run on the same port
type Server struct {
	gServer *grpc.Server
	httpMux *http.ServeMux
	port    int
}

//NewServer creates a new Server given the configuration options
func NewServer(ctx context.Context, opts ...ConfigOption) (*Server, error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	gServer := grpc.NewServer(cfg.serverOptions...)
	cfg.serviceRegistration(gServer)
	proxy := runtime.NewServeMux(cfg.proxyOptions...)
	cfg.proxyRegistration(ctx, proxy, fmt.Sprintf("localhost:%v", cfg.port), cfg.proxyClientOpts...)
	mux := http.NewServeMux()
	if cfg.muxPrefix == "" {
		cfg.muxPrefix = "/"
	}
	mux.Handle(cfg.muxPrefix, proxy)
	mux.HandleFunc("/service-info", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(gServer.GetServiceInfo()); err != nil {
			http.Error(w, "failed to encode grpc service info", http.StatusInternalServerError)
		}
	})
	return &Server{
		gServer: gServer,
		httpMux: mux,
		port:    cfg.port,
	}, nil
}

//GetGRPCServer gets the servers grpc server
func (s *Server) GetGRPCServer() *grpc.Server {
	return s.gServer
}

//GetHTTPMux gets the servers http mux, you can add additional handlers before starting server
func (s *Server) GetHTTPMux() *http.ServeMux {
	return s.httpMux
}

//Run starts the http and grpc server on the same port
func (s *Server) Run(ctx context.Context) error {
	hServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", s.port),
		Handler: s.httpMux,
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		return err
	}
	mux := cmux.New(lis)
	gMux := mux.Match(cmux.HTTP2())
	hMux := mux.Match(cmux.Any())

	fmt.Printf("starting grpc/json server on port %v\n", s.port)
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return s.gServer.Serve(gMux)
	})
	group.Go(func() error {
		return hServer.Serve(hMux)
	})
	group.Go(func() error {
		return mux.Serve()
	})
	return group.Wait()
}
