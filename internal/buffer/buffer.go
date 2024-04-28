package buffer

type (
	Data interface {
	}

	BufferManager interface {
		AddData(data Data)
		GetNextData() (Data, bool)
		Close()
	}

	bufferManager struct {
		buffer  []Data
		ch      chan Data
		closing chan struct{}
		closed  bool
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

	bm := &bufferManager{
		buffer:  make([]Data, 0, options.MaxBufferSize),
		ch:      make(chan Data, options.MaxBufferSize),
		closing: make(chan struct{}),
	}

	go bm.manageBuffer()
	return bm
}

func (bm *bufferManager) Close() {
	if bm.closed {
		return
	}
	bm.closed = true
	close(bm.closing)
	close(bm.ch)
}

func (bm *bufferManager) manageBuffer() {
	for {
		select {
		case data, ok := <-bm.ch:
			if !ok {
				return // Channel is closed, exit the goroutine
			}
			if len(bm.buffer) == cap(bm.buffer) {
				// Only remove the oldest element if there's more than one element
				if len(bm.buffer) > 1 {
					bm.buffer = bm.buffer[1:]
				} else {
					bm.buffer = bm.buffer[:0] // Reset buffer if it contains only one element
				}
			}
			bm.buffer = append(bm.buffer, data)
		case <-bm.closing:
			return // Close signal received, exit the goroutine
		}
	}
}

func (bm *bufferManager) AddData(data Data) {
	if bm.closed {
		return // Do not add data if BufferManager is closed
	}
	select {
	case bm.ch <- data:
		// Data added successfully
	default:
		// Channel is full, handle accordingly
	}
}

func (bm *bufferManager) GetNextData() (Data, bool) {
	if len(bm.buffer) == 0 {
		return nil, false
	}
	data := bm.buffer[0]
	bm.buffer = bm.buffer[1:]
	return data, true
}
