package protoio

import (
	"bufio"
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/multiformats/go-varint"
)

type uvarintReader struct {
	r       *bufio.Reader
	buf     []byte
	maxSize int
	closer  io.Closer
}

func NewReader(r io.Reader, maxSize int) ReadCloser {
	var closer io.Closer
	if c, ok := r.(io.Closer); ok {
		closer = c
	}
	return &uvarintReader{bufio.NewReader(r), nil, maxSize, closer}
}

func (ur *uvarintReader) ReadMsg(msg proto.Message) error {
	length64, err := varint.ReadUvarint(ur.r)
	if err != nil {
		return err
	}
	length := int(length64)
	if length < 0 || length > ur.maxSize {
		return io.ErrShortBuffer
	}
	if len(ur.buf) < length {
		ur.buf = make([]byte, length)
	}
	buf := ur.buf[:length]
	if _, err := io.ReadFull(ur.r, buf); err != nil {
		return err
	}
	return proto.Unmarshal(buf, msg)
}

func (ur *uvarintReader) Close() error {
	if ur.closer != nil {
		return ur.closer.Close()
	}
	return nil
}
