package protoio

import "io"

type Stream struct {
	*Reader
	*Writer
}

func NewStream(readWriter io.ReadWriter, maxReadSize int) *Stream {
	reader := NewReader(readWriter, maxReadSize)
	writer := NewWriter(readWriter)
	return NewStreamReaderWriter(reader, writer)
}

func NewStreamReaderWriter(reader *Reader, writer *Writer) *Stream {
	return &Stream{
		Reader: reader,
		Writer: writer,
	}
}
