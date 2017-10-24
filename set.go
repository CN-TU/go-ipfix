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

func MakeSet(buffer SerializeBuffer) Set {
	return Set{
		buffer: buffer,
	} // 0 id is illegal in ipfix
}

func (s *Set) startSet(record Record) (err error) {
	length := record.Length() + 4
	if s.buffer.BytesFree() < length {
		return BufferFullError(length)
	}
	s.length = length
	s.id = record.Id()
	b := s.buffer.Append(4)
	_ = b[3]
	binary.BigEndian.PutUint16(b[0:2], uint16(s.id))
	s.lengthBytes = b[2:4]
	record.SerializeTo(s.buffer)
	return
}

func (s *Set) AppendRecord(record Record) error {
	if s.id == 0 {
		if err := s.startSet(record); err != nil {
			return err
		}
		return nil
	}
	if record.Id() != s.id {
		return &RecordMismatchError{record.Id(), s.id}
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
	s.id = 0 // 0 id is illegal in ipfix
	s.length = 0
	return s.length, nil
}
