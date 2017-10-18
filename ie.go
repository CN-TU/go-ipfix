package ipfix

import (
	"encoding/binary"
	"fmt"
)

type InformationElement struct {
	Name   string
	Pen    uint32
	ID     uint16
	Type   Type
	Length uint16
}

func NewInformationElement(name string, pen uint32, id uint16, t Type, length uint16) InformationElement {
	if length == 0 {
		length = DefaultSize[t]
	}
	return InformationElement{name, pen, id, t, length}
}

func (ie InformationElement) String() string {
	if ie.Pen == 0 {
		return ie.Name
	}
	if ie.Length == 0 || ie.Length == DefaultSize[ie.Type] {
		return fmt.Sprintf("%s(%d/%d)<%s>", ie.Name, ie.Pen, ie.ID, ie.Type)
	}
	return fmt.Sprintf("%s(%d/%d)<%s>[%d]", ie.Name, ie.Pen, ie.ID, ie.Type, ie.Length)
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
