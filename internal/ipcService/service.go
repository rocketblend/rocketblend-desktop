package ipcService

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
)

const DefaultPort = 36363

type (
	Args []string

	InterProcessService struct {
		logger logger.Logger
		port   int

		rpcServer *rpc.Server
		listener  net.Listener

		ch chan Args

		mu sync.Mutex
	}

	Options struct {
		Logger      logger.Logger
		Port        int
		ArgsChannel chan Args
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithArgsChannel(ch chan Args) Option {
	return func(o *Options) {
		o.ArgsChannel = ch
	}
}

func New(opts ...Option) (*InterProcessService, error) {
	options := &Options{
		Logger: logger.NoOp(),
		Port:   DefaultPort,
	}

	for _, o := range opts {
		o(options)
	}

	if options.ArgsChannel == nil {
		return nil, fmt.Errorf("args channel is nil")
	}

	options.Logger.Info("Starting interprocess server", map[string]interface{}{"port": options.Port})

	return &InterProcessService{
		logger: options.Logger,
		port:   options.Port,
		ch:     options.ArgsChannel,
	}, nil
}

func (s *InterProcessService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rpcServer = rpc.NewServer()
	if err := s.rpcServer.Register(s.HandleArgs); err != nil {
		return err
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.logger.Error("Failed to listen", map[string]interface{}{"error": err.Error()})
		return err
	}
	s.listener = ln

	// Accept connections and handle RPCs in a separate goroutine.
	go func() {
		s.rpcServer.Accept(ln)
	}()

	return nil
}

func (s *InterProcessService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}

	if s.rpcServer != nil {
		s.rpcServer = nil
	}

	return nil
}

// HandleArgs is an RPC method that handles incoming arguments.
func (s *InterProcessService) HandleArgs(args *Args, reply *bool) error {
	s.logger.Info("Received arguments", map[string]interface{}{"args": *args})

	s.ch <- *args

	*reply = true
	return nil
}
