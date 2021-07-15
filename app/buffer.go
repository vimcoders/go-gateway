package app

const (
	DefaultBufferSize = 512
)

type Buffer struct {
	buf []byte
}

func (b *Buffer) Take(n int) []byte {
	if n < len(b.buf) {
		return b.buf[:n]
	}

	return make([]byte, n)
}

func NewBuffer() *Buffer {
	return NewBufferSize(DefaultBufferSize)
}

func NewBufferSize(n int) *Buffer {
	return &Buffer{
		buf: make([]byte, n),
	}
}
