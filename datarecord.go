package ipfix

import "io"

type BufferedDataRecord struct {
	template int16
	buffer   []byte
}

func MakeBufferedDataRecord(num int) BufferedDataRecord {
	return BufferedDataRecord{buffer: make([]byte, 0, num)}
}

func (b BufferedDataRecord) Id() int16 {
	return b.template
}

func (b *BufferedDataRecord) Append(num int) []byte {
	blen := len(b.buffer)
	newlen := blen + num
	if newlen > cap(b.buffer) {
		b.buffer = append(b.buffer, make([]byte, newlen-cap(b.buffer))...)
	}
	b.buffer = b.buffer[:newlen]
	return b.buffer[blen:newlen]
}

func (b BufferedDataRecord) Length() int {
	return len(b.buffer)
}

func (b *BufferedDataRecord) SerializeTo(buffer SerializeBuffer) {
	buf := buffer.Append(len(b.buffer))
	copy(buf, b.buffer)
	b.buffer = b.buffer[:0]
}

func (b BufferedDataRecord) BytesFree() int {
	panic("Not implemented")
}

func (b BufferedDataRecord) Finalize(io.Writer) (err error) {
	panic("Not implemented")
}
