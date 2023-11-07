package protoio

import (
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/multiformats/go-varint"
)

type uvarintWriter struct {
	w            io.Writer
	varintBuffer []byte
	buffer       []byte
}

func NewWriter(w io.Writer) WriteCloser {
	return &uvarintWriter{w, make([]byte, varint.MaxLenUvarint63), nil}
}

func (uw *uvarintWriter) WriteMsg(msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	length := uint64(len(data))
	n := varint.PutUvarint(uw.varintBuffer, length)
	_, err = uw.w.Write(uw.varintBuffer[:n])
	if err != nil {
		return err
	}
	_, err = uw.w.Write(data)
	return err
}

func (uw *uvarintWriter) Close() error {
	if closer, ok := uw.w.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
