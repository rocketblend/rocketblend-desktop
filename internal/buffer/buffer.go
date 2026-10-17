package buffer

import "sync"

type (
	Data interface {
	}

	BufferManager interface {
		AddData(data Data)
		Channel() <-chan Data
		Close()
	}

	bufferManager struct {
		mu     sync.RWMutex
		ch     chan Data
		closed bool
	}

	Options struct {
		MaxBufferSize int
	}

	Option func(*Options)
)

func WithMaxBufferSize(size int) Option {
	return func(o *Options) {
		o.MaxBufferSize = size
	}
}

func New(opts ...Option) BufferManager {
	options := &Options{
		MaxBufferSize: 100, // Default buffer size
	}

	for _, opt := range opts {
		opt(options)
	}

	return &bufferManager{
		ch: make(chan Data, options.MaxBufferSize),
	}
}

func (bm *bufferManager) Close() {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if bm.closed {
		return
	}

	bm.closed = true
	close(bm.ch)
}

// AddData drops the oldest data when the buffer is full.
func (bm *bufferManager) AddData(data Data) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if bm.closed {
		return
	}

	for {
		select {
		case bm.ch <- data:
			return
		default:
		}

		select {
		case <-bm.ch:
		default:
		}
	}
}

func (bm *bufferManager) Channel() <-chan Data {
	return bm.ch
}
