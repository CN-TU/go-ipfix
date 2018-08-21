package ipfix

import (
	"io"
)

type scratchBuffer interface {
	append(int) []byte
	bytesFree() int
	finalize(io.Writer) (err error)
	length() int
}

type basicBuffer []byte

func makeBasicBuffer(size int) scratchBuffer {
	ret := basicBuffer(make([]byte, 0, size))
	return &ret
}

func (b *basicBuffer) bytesFree() int {
	return cap(*b) - len(*b)
}

func (b *basicBuffer) append(num int) []byte {
	blen := len(*b)
	*b = (*b)[:blen+num]
	return (*b)[blen : blen+num]
}

func (b *basicBuffer) finalize(w io.Writer) (err error) {
	_, err = w.Write(*b)
	if err == nil {
		*b = (*b)[:0]
	}
	return
}

func (b *basicBuffer) length() int {
	return len(*b)
}
