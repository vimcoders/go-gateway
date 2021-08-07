package apk

import (
	"bufio"
	"io"
	"net"
	"time"
)

type Writer struct {
	net.Conn
	*bufio.Writer
	t time.Duration
}

func (w *Writer) Write(p []byte) (n int, err error) {
	defer func() {
		if err := w.Writer.Flush(); err != nil {
			logger.Error("writer %v", err)
		}
	}()

	length := len(p)

	var header [2]byte

	header[0] = uint8(length >> 8)
	header[1] = uint8(length)

	if err := w.SetWriteDeadline(time.Now().Add(w.t)); err != nil {
		return 0, err
	}

	if n, err := w.Writer.Write(header[:]); err != nil {
		return n, err
	}

	return w.Writer.Write(p)
}

func NewWriter(c net.Conn) io.Writer {
	return &Writer{c, bufio.NewWriter(c), time.Second * 15}
}

type Reader struct {
	net.Conn
	*bufio.Reader
	t time.Duration
}

func (r *Reader) Read() (p []byte, err error) {
	header := make([]byte, 2)

	if err := r.SetReadDeadline(time.Now().Add(r.t)); err != nil {
		return nil, err
	}

	if _, err := r.Reader.Read(header); err != nil {
		return nil, err
	}

	length := uint16(uint16(header[0])<<8 | uint16(header[1]))

	if err := r.SetReadDeadline(time.Now().Add(r.t)); err != nil {
		return nil, err
	}

	return r.Reader.Peek(int(length))
}

func NewReader(c net.Conn) *Reader {
	return &Reader{c, bufio.NewReaderSize(c, 512), time.Second * 5}
}
