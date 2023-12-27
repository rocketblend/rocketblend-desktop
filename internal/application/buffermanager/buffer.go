package buffermanager

type (
	Data interface{}

	BufferManager interface {
		AddData(data Data)
		GetNextData() (Data, bool)
	}

	bufferManager struct {
		buffer []Data
		ch     chan Data
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
		buffer: make([]Data, 0, options.MaxBufferSize),
		ch:     make(chan Data, options.MaxBufferSize),
	}

	go bm.manageBuffer()
	return bm
}

func (bm *bufferManager) manageBuffer() {
	for data := range bm.ch {
		if len(bm.buffer) == cap(bm.buffer) {
			bm.buffer = bm.buffer[1:]
		}
		bm.buffer = append(bm.buffer, data)
	}
}

func (bm *bufferManager) AddData(data Data) {
	select {
	case bm.ch <- data:
		// Data added successfully
	default:
		// Channel is full, so this will be handled in manageBuffer
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
