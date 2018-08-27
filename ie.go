package ipfix

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

// basicListID is the IE ID of basicLists as defined by RFC6313
const basicListID = 291

// subType represents additional information needed by RFC6313 datatypes
type subType interface{}

// InformationElement represents the description of an information element according to RFC7011
type InformationElement struct {
	// Name of the information element
	Name string
	// Pen is the enterprise number (0 for iana reserved)
	Pen uint32
	// ID is the information element ID
	ID uint16
	// Type is the associated data type
	Type Type
	// Length is the length of the field value
	Length  uint16
	subType subType
}

// NewInformationElement returns an information element for the given specification. If length is 0,
// the default length for this data type is chosen.
func NewInformationElement(name string, pen uint32, id uint16, t Type, length uint16) InformationElement {
	if t != IllegalType && length == 0 {
		length = DefaultSize[t]
	}
	return InformationElement{name, pen, id, t, length, nil}
}

// NewBasicList returns an InformationElement holding the basic list according to RFC6313. If number is 0,
// a variable length list is returned.
func NewBasicList(name string, subelement InformationElement, number uint16) InformationElement {
	// RFC6313: semantic + template of element + number of elements * size of element
	length := 1 + uint16(subelement.templateSize()) + number*subelement.Length
	if number == 0 || subelement.Length == VariableLength || number == VariableLength {
		length = VariableLength
	}
	return InformationElement{name, 0, basicListID, BasicListType, length, subelement}
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

func (ie InformationElement) templateSize() int {
	if ie.Pen == 0 {
		return 4
	}
	return 8
}

func (ie InformationElement) serializeTo(buffer scratchBuffer) (int, error) {
	ident := ie.ID
	if ie.Pen == 0 {
		b, err := buffer.append(4)
		if err != nil {
			return 0, err
		}
		binary.BigEndian.PutUint16(b[2:], uint16(ie.Length))
		binary.BigEndian.PutUint16(b[0:], uint16(ident))
		return 4, nil
	}
	ident |= 0x8000
	b, err := buffer.append(8)
	if err != nil {
		return 0, err
	}
	binary.BigEndian.PutUint32(b[4:], uint32(ie.Pen))
	binary.BigEndian.PutUint16(b[2:], uint16(ie.Length))
	binary.BigEndian.PutUint16(b[0:], uint16(ident))
	return 8, nil
}

// ListElement returns the InformationElement of a list item and true if this InformationElement is a list.
// Otherwise an empty InformationElement and false is returned.
func (ie InformationElement) ListElement() (InformationElement, bool) {
	if ie.Type != BasicListType {
		return InformationElement{}, false
	}
	return ie.subType.(InformationElement), true
}

func (ie InformationElement) serializeDataTo(buffer scratchBuffer, value interface{}) error {
	switch ie.Type {
	case BasicListType:
		subie, _ := ie.ListElement()
		// Header according to RFC6313
		written := 1

		var lengthbuffer []byte

		if ie.Length == VariableLength {
			// RFC6313 recommends 3 byte encoding of length field
			b, err := buffer.append(3)
			if err != nil {
				return err
			}
			_ = b[2]
			b[0] = 0xff
			lengthbuffer = b[1:3]
		}

		// first semantic
		b, err := buffer.append(1)
		if err != nil {
			return err
		}
		b[0] = byte(UndefinedSemantic)
		// followed by template header
		headersize, err := subie.serializeTo(buffer)
		if err != nil {
			return err
		}
		written += headersize
		// followed by all the values
		if value != nil {
			values := reflect.ValueOf(value)
			for values.Kind() == reflect.Ptr {
				values = values.Elem()
			}
			l := values.Len()
			for i := 0; i < l; i++ {
				subiesize, err := subie.Type.serializeDataTo(buffer, values.Index(i).Interface(), int(subie.Length))
				if err != nil {
					return err
				}
				written += subiesize
			}
		}
		if ie.Length == VariableLength {
			binary.BigEndian.PutUint16(lengthbuffer, uint16(written))
		} else {
			if written != int(ie.Length) {
				panic("Number of values doesn't fit ie length")
			}
		}
	default:
		ie.Type.serializeDataTo(buffer, value, int(ie.Length))
	}
	return nil
}
