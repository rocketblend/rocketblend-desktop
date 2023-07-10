package ipcService

import "github.com/flowshot-io/x/pkg/logger"

type Handler struct {
	logger logger.Logger
	ch     chan Args
}

// HandleArgs is an RPC method that handles incoming arguments.
func (h *Handler) HandleArgs(args *Args, reply *bool) error {
	h.logger.Info("Received arguments", map[string]interface{}{"args": *args})

	h.ch <- *args

	*reply = true
	return nil
}
