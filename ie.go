package ipfix

import (
	"encoding/binary"
)

type InformationElement struct {
	Name   string
	Pen    uint32
	ID     uint16
	Type   Type
	Length uint16
}

func (ie *InformationElement) TemplateSize() int {
	if ie.Pen == 0 {
		return 4
	}
	return 8
}

func (ie *InformationElement) SerializeTo(buffer SerializeBuffer) {
	ident := ie.ID
	if ie.Pen == 0 {
		b := buffer.Append(4)
		binary.BigEndian.PutUint16(b[2:], uint16(ie.Length))
		binary.BigEndian.PutUint16(b[0:], uint16(ident))
		return
	}
	ident |= 0x8000
	b := buffer.Append(8)
	binary.BigEndian.PutUint32(b[4:], uint32(ie.Pen))
	binary.BigEndian.PutUint16(b[2:], uint16(ie.Length))
	binary.BigEndian.PutUint16(b[0:], uint16(ident))
}

func (ie *InformationElement) MakeDataRecord(value interface{}) DataRecordElement {
	return MakeDataRecordElement(ie.Type, value, int(ie.Length))
}
