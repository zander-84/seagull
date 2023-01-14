package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/zander-84/seagull/internal/endpoint"
	"github.com/zander-84/seagull/internal/host"
	"github.com/zander-84/seagull/transport"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

// ServerHandler with server handler.
func ServerHandler(h http.Handler) ServerOption {
	return func(s *Server) {
		s.Server.Handler = h
	}
}

// ServerTLSConfig with server tls config.
func ServerTLSConfig(c *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConf = c
	}
}

// ServerReadTimeout with read timeout.
func ServerReadTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// ServerWriteTimeout with write timeout.
func ServerWriteTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

// ServerIdleTimeout with read timeout.
func ServerIdleTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.idleTimeout = timeout
	}
}

// Listener with server lis
func Listener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

// Server is an HTTP server wrapper.
type Server struct {
	*http.Server
	err          error
	name         string
	lis          net.Listener
	endpoint     *url.URL
	network      string
	runIp        string
	registerIp   string
	port         int
	tlsConf      *tls.Config
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

// NewServer creates  HTTP server by options.
func NewServer(name string, runIp string, registerIp string, port int, opts ...ServerOption) *Server {
	srv := &Server{
		network:      "tcp",
		name:         name,
		runIp:        runIp,
		registerIp:   registerIp,
		port:         port,
		readTimeout:  10 * time.Second,
		writeTimeout: 60 * time.Second,
		idleTimeout:  10 * time.Second,
	}

	h := http.NewServeMux()
	srv.Server = &http.Server{Handler: h}

	for _, o := range opts {
		o(srv)
	}
	return srv
}
func (s *Server) address() string {
	return fmt.Sprintf("%s:%d", s.runIp, s.port)
}
func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address())
		if err != nil {
			s.err = err
			return err
		}
		s.lis = lis
	}
	if s.endpoint == nil {
		addr, err := host.Extract(fmt.Sprintf("%s:%d", s.registerIp, s.port), s.lis)
		if err != nil {
			s.err = err
			return err
		}
		s.endpoint = endpoint.NewEndpoint(s.name, addr)
	}
	return s.err
}

// Endpoint return a real address to registry endpoint.
// examples:
//
//	https://127.0.0.1:8000
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

// Start the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return err
	}
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	logStr := fmt.Sprintf("[HTTP] server listening on: %s ", s.lis.Addr().String())
	if s.endpoint != nil {
		logStr += fmt.Sprintf("endpiont on: %v ", s.endpoint)
	}
	log.Println(logStr)

	var err error
	if s.tlsConf != nil {
		err = s.Server.ServeTLS(s.lis, "", "")
	} else {
		err = s.Server.Serve(s.lis)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err == nil {
		log.Printf("[HTTP]  GracefulStop On: %s\n", s.lis.Addr().String())
	} else {
		log.Printf("[HTTP]  errServerClosed: %e", err)
	}
	return err
}
