package ipfix

import (
	"encoding/binary"
)

type template struct {
	identifier int16
	elements   []InformationElement
}

func (t template) id() int16 {
	return templateSetID
}

func (t template) length() (ret int) {
	ret = 4
	for _, element := range t.elements {
		ret += element.templateSize()
	}
	return
}

func (t template) serializeTo(buffer scratchBuffer) error {
	b, err := buffer.append(4)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint16(b[2:], uint16(len(t.elements)))
	binary.BigEndian.PutUint16(b[0:], uint16(t.identifier))
	for _, element := range t.elements {
		_, err := element.serializeTo(buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t template) assignDataRecord(record *recordBuffer, values ...interface{}) error {
	if len(values) != len(t.elements) {
		return TemplateMismatchError{len(values), len(t.elements)}
	}
	// TODO: reset record if serialize fails
	record.template = t.identifier
	values = values[:len(t.elements)]
	for i, element := range t.elements {
		element.serializeDataTo(record, values[i])
	}
	return nil
}
