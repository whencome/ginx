package server

import (
    "context"
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
)

// define gin run mode constant
const (
    ModeDebug   = "debug"
    ModeRelease = "release"
)

// Options http server run options
type Options struct {
    Port     int    `json:"port" yaml:"port" toml:"port"`                // 服务端口
    Mode     string `json:"mode" yaml:"mode" toml:"mode"`                // 运行模式, debug or release
    Tls      bool   `json:"tls" yaml:"tls" toml:"tls"`                   // 是否启用HTTPS
    CertFile string `json:"cert_file" yaml:"cert_file" toml:"cert_file"` // 证书文件
    KeyFile  string `json:"key_file" yaml:"key_file" toml:"key_file"`    // 密钥文件
}

// HookFunc http server init & stop hooks
type HookFunc func(r *gin.Engine) error

// DefaultOptions create a default options
func DefaultOptions() *Options {
    return &Options{
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
    options *Options
    // hooks of init server
    preInitFunc  HookFunc
    postInitFunc HookFunc
    // hooks of stop server
    preStopFunc  HookFunc
    postStopFunc HookFunc
}

// New create a http server
func New(options *Options) *HTTPServer {
    if options == nil {
        options = DefaultOptions()
    }
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

func (s *HTTPServer) PreInit(f HookFunc) {
    s.preInitFunc = f
}

func (s *HTTPServer) PostInit(f HookFunc) {
    s.postInitFunc = f
}

func (s *HTTPServer) PreStop(f HookFunc) {
    s.preStopFunc = f
}

func (s *HTTPServer) PostStop(f HookFunc) {
    s.postStopFunc = f
}

func (s *HTTPServer) execHook(f HookFunc) error {
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
        return errors.New("http server not runnable, it's probably has already started")
    }
    // set gin run mode
    if s.options.Mode != ModeDebug {
        gin.SetMode(gin.ReleaseMode)
    }

    // 初始化server
    if e := s.execHook(s.preInitFunc); e != nil {
        return e
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

    // 启动http服务
    s.running = true
    go func() {
        s.Wait()
    }()
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

    // 启动http服务
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
        log.Printf("http server started on %s", s.svr.Addr)
        return true, nil
    }
}

// Stop stop the server
func (s *HTTPServer) Stop() {
    // exec pre stop hook
    err := s.execHook(s.preStopFunc)
    if err != nil {
        log.Printf("prepare stop server failed: %s \n", err)
    }
    // shutdown the http server
    log.Println("start to shutdown http server")
    if err = s.svr.Shutdown(context.Background()); err != nil {
        log.Printf("shutdown server failed: %s \n", err)
        return
    }
    s.running = false
    // exec post stop hook
    err = s.execHook(s.postStopFunc)
    if err != nil {
        log.Printf("stop server error: %s \n", err)
    }
    log.Println("http server closed")
}

// Wait block and wait
func (s *HTTPServer) Wait() {
    sigChan := make(chan os.Signal)
    signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
    for {
        sig := <-sigChan
        switch sig {
        case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
            s.Stop()
            time.Sleep(1 * time.Second)
            return
        }
    }
}
