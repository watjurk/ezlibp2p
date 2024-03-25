package protoio

import "io"

type byteReaderT struct {
	reader io.Reader
	buffer []byte
}

func newByteReaderT(reader io.Reader) *byteReaderT {
	return &byteReaderT{
		reader: reader,
		buffer: make([]byte, 1),
	}
}

func (br *byteReaderT) ReadByte() (byte, error) {
	_, err := io.ReadFull(br.reader, br.buffer)
	if err != nil {
		return 0, err
	}

	return br.buffer[0], nil
}
