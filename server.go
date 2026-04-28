package ginx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whencome/ginx/log"
)

// define gin run mode constant
const (
	ModeTest    = "test"
	ModeDebug   = "debug"
	ModeRelease = "release"
)

// ServerOptions http server run options
type ServerOptions struct {
	Port     int    `json:"port" yaml:"port" toml:"port"`                // server port
	Mode     string `json:"mode" yaml:"mode" toml:"mode"`                // run mode, debug or release
	Tls      bool   `json:"tls" yaml:"tls" toml:"tls"`                   // enable HTTPS
	CertFile string `json:"cert_file" yaml:"cert_file" toml:"cert_file"` // certificate file
	KeyFile  string `json:"key_file" yaml:"key_file" toml:"key_file"`    // key file
}

// ServerHookFunc http server init & stop hooks
type ServerHookFunc func(r *gin.Engine) error

// DefaultServerOptions create default options
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		Port: 8080,
		Mode: ModeRelease,
		Tls:  false,
	}
}

// HTTPServer define a simple http server
type HTTPServer struct {
	running bool
	engine  *gin.Engine
	svr     *http.Server
	options *ServerOptions
	// hooks of init server
	postInitFunc ServerHookFunc
	// hooks of stop server
	preStopFunc  ServerHookFunc
	postStopFunc ServerHookFunc
}

// NewServer create a http server
func NewServer(options *ServerOptions) *HTTPServer {
	if options == nil {
		options = DefaultServerOptions()
	}
	// set gin run mode
	switch options.Mode {
	case ModeTest:
		gin.SetMode(gin.TestMode)
	case ModeDebug:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	// create http server
	s := &HTTPServer{
		running: false,
		engine:  gin.Default(),
		svr:     nil,
		options: options,
	}
	s.svr = &http.Server{
		Addr:    fmt.Sprintf(":%d", options.Port),
		Handler: s.engine,
	}
	return s
}

func (s *HTTPServer) GinEngine() *gin.Engine {
	return s.engine
}

func (s *HTTPServer) PostInit(f ServerHookFunc) {
	s.postInitFunc = f
}

func (s *HTTPServer) PreStop(f ServerHookFunc) {
	s.preStopFunc = f
}

func (s *HTTPServer) PostStop(f ServerHookFunc) {
	s.postStopFunc = f
}

func (s *HTTPServer) execHook(f ServerHookFunc) error {
	if f == nil {
		return nil
	}
	return f(s.engine)
}

// Runnable check whether server is runnable
func (s *HTTPServer) Runnable() bool {
	return !s.running
}

func (s *HTTPServer) prepare() error {
	if !s.Runnable() {
		return errors.New("http server not runnable, it has probably already started")
	}
	if e := s.execHook(s.postInitFunc); e != nil {
		return e
	}
	return nil
}

// Run start http server in block mode
func (s *HTTPServer) Run() error {
	// prepare server
	e := s.prepare()
	if e != nil {
		return e
	}

	// start http service
	s.running = true
	if s.options.Tls {
		return s.svr.ListenAndServeTLS(s.options.CertFile, s.options.KeyFile)
	}
	return s.svr.ListenAndServe()
}

// Start start http server in non-blocked mode
func (s *HTTPServer) Start() (bool, error) {
	// prepare server
	e := s.prepare()
	if e != nil {
		return false, e
	}

	// start http service in non-blocking mode
	s.running = true
	startCh := make(chan error)
	go func() {
		if s.options.Tls {
			if err := s.svr.ListenAndServeTLS(s.options.CertFile, s.options.KeyFile); err != nil {
				startCh <- err
			}
		} else {
			if err := s.svr.ListenAndServe(); err != nil {
				startCh <- err
			}
		}
	}()

	select {
	case err := <-startCh:
		return false, err
	case <-time.After(time.Second * 3):
		log.Infof("http server started on %s", s.svr.Addr)
		return true, nil
	}
}

// Stop the server with graceful shutdown
func (s *HTTPServer) Stop() {
	// exec pre stop hook
	err := s.execHook(s.preStopFunc)
	if err != nil {
		log.Errorf("prepare stop server failed: %s", err)
	}
	// shutdown the http server with timeout
	log.Infof("start to shutdown http server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err = s.svr.Shutdown(ctx); err != nil {
		log.Errorf("shutdown server: %s", err)
		return
	}
	s.running = false
	// exec post stop hook
	err = s.execHook(s.postStopFunc)
	if err != nil {
		log.Errorf("stop server: %s", err)
	}
	log.Infof("http server closed")
}

// Wait block and wait for exit signal
func (s *HTTPServer) Wait() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigChan
	switch sig {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
		log.Infof("received exit signal: %v", sig)
		s.Stop()
		time.Sleep(1 * time.Second)
		return
	}
}
