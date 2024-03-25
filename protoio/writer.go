package protoio

import (
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/multiformats/go-varint"
)

// Writer is a type able to write protobuf messages.
type Writer struct {
	writer       io.Writer
	varintBuffer []byte
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer:       writer,
		varintBuffer: make([]byte, varint.MaxLenUvarint63),
	}
}

func (uw *Writer) WriteMsg(msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	length := uint64(len(data))
	n := varint.PutUvarint(uw.varintBuffer, length)
	_, err = uw.writer.Write(uw.varintBuffer[:n])
	if err != nil {
		return err
	}
	_, err = uw.writer.Write(data)
	return err
}
