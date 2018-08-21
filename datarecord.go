package ipfix

import (
	"io"
)

type recordBuffer struct {
	template int16
	buffer   []byte
}

func makeRecordBuffer(num int) recordBuffer {
	return recordBuffer{buffer: make([]byte, 0, num)}
}

func (b recordBuffer) id() int16 {
	return b.template
}

func (b *recordBuffer) append(num int) []byte {
	blen := len(b.buffer)
	newlen := blen + num
	if newlen > cap(b.buffer) {
		panic("Data record too small")
		// Never append to datarecord; This breaks slice pointers
	}
	b.buffer = b.buffer[:newlen]
	return b.buffer[blen:newlen]
}

func (b recordBuffer) length() int {
	return len(b.buffer)
}

func (b *recordBuffer) serializeTo(buffer scratchBuffer) {
	buf := buffer.append(len(b.buffer))
	copy(buf, b.buffer)
	b.buffer = b.buffer[:0]
}

func (b recordBuffer) bytesFree() int {
	panic("Not implemented")
}

func (b recordBuffer) finalize(io.Writer) (err error) {
	panic("Not implemented")
}
