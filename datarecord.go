package ipfix

import "io"

type DataRecordElement interface {
	Length() int
	SerializeTo(buffer SerializeBuffer)
}

type DataRecord struct {
	template int16
	elements []DataRecordElement
}

func (d DataRecord) Id() int16 {
	return d.template
}

func (d DataRecord) Length() (ret int) {
	for _, element := range d.elements {
		ret += element.Length()
	}
	return
}

func (d DataRecord) SerializeTo(buffer SerializeBuffer) {
	for _, element := range d.elements {
		element.SerializeTo(buffer)
	}
}

type BufferedDataRecord struct {
	DataRecord
	buffer []byte
}

func MakeBufferedDataRecord(num int) BufferedDataRecord {
	return BufferedDataRecord{buffer: make([]byte, 0, num)}
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

func (b *BufferedDataRecord) BytesFree() int {
	panic("Not implemented")
}

func (b *BufferedDataRecord) Finalize(io.Writer) (err error) {
	panic("Not implemented")
}
