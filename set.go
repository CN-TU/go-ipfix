package ipfix

import (
	"encoding/binary"
)

type Record interface {
	Length() int
	SerializeTo(SerializeBuffer)
	Id() int16
}

const (
	IDTemplateSet        int16 = 2
	IDOptionsTemplateSet int16 = 3
)

type Set struct {
	buffer      SerializeBuffer
	lengthBytes []byte
	length      int
	id          int16
}

func MakeSet(record Record, buffer SerializeBuffer) (ret *Set, err error) {
	len := record.Length() + 4
	if buffer.BytesFree() < len {
		return nil, BufferFullError(len)
	}
	ret = &Set{
		buffer: buffer,
		id:     record.Id(),
		length: len,
	}
	b := buffer.Append(4)
	_ = b[3]
	binary.BigEndian.PutUint16(b[0:2], uint16(ret.id))
	ret.lengthBytes = b[2:4]
	record.SerializeTo(buffer)
	return
}

func (s *Set) Id() int16 {
	return s.id
}

func (s *Set) AppendRecord(record Record) error {
	if record.Id() != s.Id() {
		return &RecordMismatchError{record.Id(), s.Id()}
	}
	len := record.Length()
	if s.buffer.BytesFree() < len {
		return BufferFullError(len)
	}
	s.length += len
	record.SerializeTo(s.buffer)
	return nil
}

func (s *Set) Finalize() (int, error) {
	if s.length == 0 {
		return 0, &IllegalSetError{"Set does not contain data!"}
	}
	binary.BigEndian.PutUint16(s.lengthBytes, uint16(s.length))
	return s.length, nil
}
