package protoio

import (
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/multiformats/go-varint"
)

// Reader is a type able to read protobuf messages.
//
// Reader reads nothing more than the specified message,
// in other words: if there is some data right after the message that data is fully preserved.
//
// Reader requires specifying the maxReadSize so if a bad actor sends an
// arbitrary large message it would be discarded.
type Reader struct {
	byteReader io.ByteReader
	reader     io.Reader

	buf         []byte
	maxReadSize int
}

func NewReader(reader io.Reader, maxReadSize int) *Reader {
	byteReader := newByteReaderT(reader)
	return &Reader{
		byteReader:  byteReader,
		reader:      reader,
		buf:         nil,
		maxReadSize: maxReadSize,
	}
}

func (ur *Reader) ReadMsg(msg proto.Message) error {
	length64, err := varint.ReadUvarint(ur.byteReader)
	if err != nil {
		return err
	}

	length := int(length64)
	if length < 0 || length > ur.maxReadSize {
		return io.ErrShortBuffer
	}

	if len(ur.buf) < length {
		ur.buf = make([]byte, length)
	}

	buf := ur.buf[:length]
	if _, err := io.ReadFull(ur.reader, buf); err != nil {
		return err
	}

	return proto.Unmarshal(buf, msg)
}
