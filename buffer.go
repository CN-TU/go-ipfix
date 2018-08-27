package ipfix

import (
	"io"
)

type scratchBuffer interface {
	append(int) ([]byte, error)
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

func (b *basicBuffer) append(num int) ([]byte, error) {
	blen := len(*b)
	if blen+num > cap(*b) {
		return nil, bufferFullError(blen + num - cap(*b))
	}
	*b = (*b)[:blen+num]
	return (*b)[blen : blen+num], nil
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
