package ipfix

import (
	"encoding/binary"
)

type record interface {
	length() int
	serializeTo(scratchBuffer) error
	id() int16
}

const (
	templateSetID        int16 = 2
	optionsTemplateSetID int16 = 3
)

type set struct {
	buffer      scratchBuffer
	lengthBytes []byte
	length      int
	id          int16
}

func makeSet(buffer scratchBuffer) set {
	return set{
		buffer: buffer,
	} // 0 id is illegal in ipfix
}

func (s *set) startSet(rec record) error {
	length := rec.length() + 4
	if s.buffer.bytesFree() < length {
		return bufferFullError(length)
	}
	s.length = length
	s.id = rec.id()
	b, err := s.buffer.append(4)
	if err != nil {
		return err
	}
	_ = b[3]
	binary.BigEndian.PutUint16(b[0:2], uint16(s.id))
	s.lengthBytes = b[2:4]
	return rec.serializeTo(s.buffer)
}

func (s *set) appendRecord(rec record) error {
	if s.id == 0 {
		if err := s.startSet(rec); err != nil {
			return err
		}
		return nil
	}
	if rec.id() != s.id {
		return &recordMismatchError{rec.id(), s.id}
	}
	len := rec.length()
	if s.buffer.bytesFree() < len {
		return bufferFullError(len)
	}
	s.length += len
	return rec.serializeTo(s.buffer)
}

func (s *set) finalize() (int, error) {
	if s.length == 0 {
		return 0, nil
	}
	binary.BigEndian.PutUint16(s.lengthBytes, uint16(s.length))
	s.id = 0 // 0 id is illegal in ipfix
	s.length = 0
	return s.length, nil
}
