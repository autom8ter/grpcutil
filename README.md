# grpcutil
--
    import "github.com/autom8ter/grpcutil"


## Usage

#### func  AllHeaderMatcherFunc

```go
func AllHeaderMatcherFunc() runtime.HeaderMatcherFunc
```
AllHeaderMatcherFunc forwards all headers to the grpc server

#### func  MappedHeaderMatcherFunc

```go
func MappedHeaderMatcherFunc(decider map[string]bool) runtime.HeaderMatcherFunc
```
MappedHeaderMatcherFunc is a map from header key - a boolean indicating whether
to forward the header to the grpc server

#### type Config

```go
type Config struct {
}
```

Config configures a grpc/json server on the same port

#### type ConfigOption

```go
type ConfigOption func(c *Config)
```

ConfigOption is a first class function used for grpc/proxy initialization

#### func  WithClientProxyDialOptions

```go
func WithClientProxyDialOptions(opts ...grpc.DialOption) ConfigOption
```
WithClientProxyDialOptions registers client dial options between gateway and
grpc server

#### func  WithGRPCServiceRegistration

```go
func WithGRPCServiceRegistration(fn func(*grpc.Server)) ConfigOption
```
WithGRPCServiceRegistration registers service implementations against the input
grpc Server (generated by grpc plugin)

#### func  WithPort

```go
func WithPort(port int) ConfigOption
```
WithPort is the port to serve on

#### func  WithProxyOptions

```go
func WithProxyOptions(opts ...runtime.ServeMuxOption) ConfigOption
```
WithProxyOptions registers header matchers(ex MappedHeaderMatcher)

#### func  WithProxyPrefix

```go
func WithProxyPrefix(prefix string) ConfigOption
```
WithProxyPrefix servers the grpc gateway proxy after this prefix

#### func  WithProxyServiceRegistration

```go
func WithProxyServiceRegistration(fn func(ctx context.Context, mux *runtime.ServeMux, port string, opts ...grpc.DialOption)) ConfigOption
```
WithProxyServiceRegistration registers json proxies against the gateway runtime
servemux(generated by grpc gateway plugin)

#### func  WithServerOptions

```go
func WithServerOptions(opts ...grpc.ServerOption) ConfigOption
```
WithServerOptions registers grpc server options like unary/stream interceptors

#### type Server

```go
type Server struct {
}
```

Server is a grpc/http server that can run on the same port

#### func  NewServer

```go
func NewServer(ctx context.Context, opts ...ConfigOption) (*Server, error)
```
NewServer creates a new Server given the configuration options

#### func (*Server) GetGRPCServer

```go
func (s *Server) GetGRPCServer() *grpc.Server
```
GetGRPCServer gets the servers grpc server

#### func (*Server) GetHTTPMux

```go
func (s *Server) GetHTTPMux() *http.ServeMux
```
GetHTTPMux gets the servers http mux, you can add additional handlers before
starting server

#### func (*Server) Run

```go
func (s *Server) Run(ctx context.Context) error
```
Run starts the http and grpc server on the same port
