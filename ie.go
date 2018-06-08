package ipfix

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

// BasicListID is the IE ID of basicLists as defined by RFC6313
const BasicListID = 291

// SubType represents additional information needed by RFC6313 datatypes
type SubType interface{}

type InformationElement struct {
	Name    string
	Pen     uint32
	ID      uint16
	Type    Type
	Length  uint16
	SubType SubType
}

func NewInformationElement(name string, pen uint32, id uint16, t Type, length uint16) InformationElement {
	if t != IllegalType && length == 0 {
		length = DefaultSize[t]
	}
	return InformationElement{name, pen, id, t, length, nil}
}

func NewBasicList(name string, subelement InformationElement, number uint16) InformationElement {
	// RFC6313: semantic + template of element + number of elements * size of element
	length := 1 + uint16(subelement.TemplateSize()) + number*subelement.Length
	if number == 0 || subelement.Length == VariableLength || number == VariableLength {
		length = VariableLength
	}
	return InformationElement{name, 0, BasicListID, BasicListType, length, subelement}
}

func NewInformationElementBasicList(name string, elementpen uint32, elementid uint16, elementt Type, elementlength uint16, number uint16) InformationElement {
	if elementt != IllegalType && elementlength == 0 {
		elementlength = DefaultSize[elementt]
	}
	return NewBasicList(name, InformationElement{name, elementpen, elementid, elementt, elementlength, nil}, number)
}

func (ie InformationElement) String() string {
	if ie.Pen == 0 {
		return ie.Name
	}
	// Output information element spec according to RFC7013 Section 10.1
	if (ie.Length == 0 && ie.Type != IllegalType) || ie.Length == DefaultSize[ie.Type] {
		return fmt.Sprintf("%s(%d/%d)<%s>", ie.Name, ie.Pen, ie.ID, ie.Type)
	}
	return fmt.Sprintf("%s(%d/%d)<%s>[%d]", ie.Name, ie.Pen, ie.ID, ie.Type, ie.Length)
}

func (ie InformationElement) TemplateSize() int {
	if ie.Pen == 0 {
		return 4
	}
	return 8
}

func (ie InformationElement) SerializeTo(buffer SerializeBuffer) int {
	ident := ie.ID
	if ie.Pen == 0 {
		b := buffer.Append(4)
		binary.BigEndian.PutUint16(b[2:], uint16(ie.Length))
		binary.BigEndian.PutUint16(b[0:], uint16(ident))
		return 4
	}
	ident |= 0x8000
	b := buffer.Append(8)
	binary.BigEndian.PutUint32(b[4:], uint32(ie.Pen))
	binary.BigEndian.PutUint16(b[2:], uint16(ie.Length))
	binary.BigEndian.PutUint16(b[0:], uint16(ident))
	return 8
}

func (ie InformationElement) ListElement() (InformationElement, bool) {
	if ie.Type != BasicListType {
		return InformationElement{}, false
	}
	return ie.SubType.(InformationElement), true
}

func (ie InformationElement) SerializeDataTo(buffer SerializeBuffer, value interface{}) {
	switch ie.Type {
	case BasicListType:
		subie, _ := ie.ListElement()
		// Header according to RFC6313
		written := 1

		var lengthbuffer []byte

		if ie.Length == VariableLength {
			// RFC6313 recommends 3 byte encoding of length field
			b := buffer.Append(3)
			_ = b[2]
			b[0] = 0xff
			lengthbuffer = b[1:3]
		}

		// first semantic
		b := buffer.Append(1)
		b[0] = byte(UndefinedSemantic)
		// followed by template header
		written += subie.SerializeTo(buffer)
		// followed by all the values
		values := reflect.ValueOf(value)
		for values.Kind() == reflect.Ptr {
			values = values.Elem()
		}
		l := values.Len()
		for i := 0; i < l; i++ {
			written += subie.Type.SerializeDataTo(buffer, values.Index(i).Interface(), int(subie.Length))
		}
		if ie.Length == VariableLength {
			binary.BigEndian.PutUint16(lengthbuffer, uint16(written))
		} else {
			if written != int(ie.Length) {
				panic("Number of values doesn't fit ie length")
			}
		}
	default:
		ie.Type.SerializeDataTo(buffer, value, int(ie.Length))
	}
}
