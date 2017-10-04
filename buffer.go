package ipfix

import (
	"io"
)

type SerializeBuffer interface {
	Append(int) []byte
	BytesFree() int
	Finalize(io.Writer) (err error)
	Length() int
}

type serializeBuffer struct {
	buffer []byte
}

func MakeSerializeBuffer(size int) SerializeBuffer {
	return &serializeBuffer{make([]byte, 0, size)}
}

func (b *serializeBuffer) BytesFree() int {
	return cap(b.buffer) - len(b.buffer)
}

func (b *serializeBuffer) Append(num int) []byte {
	blen := len(b.buffer)
	b.buffer = b.buffer[:blen+num]
	return b.buffer[blen : blen+num]
}

func (b *serializeBuffer) Finalize(w io.Writer) (err error) {
	_, err = w.Write(b.buffer)
	if err == nil {
		b.buffer = b.buffer[:0]
	}
	return
}

func (b *serializeBuffer) Length() int {
	return len(b.buffer)
}
