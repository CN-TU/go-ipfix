package ipfix

import (
	"encoding/binary"
	"fmt"
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

func (t template) assignDataRecord(record *recordBuffer, values ...interface{}) {
	if len(values) != len(t.elements) {
		panic(fmt.Sprintf("Supplied values (%d) differ from number of information elements (%d)!\n", len(values), len(t.elements)))
	}
	record.template = t.identifier
	values = values[:len(t.elements)]
	for i, element := range t.elements {
		element.serializeDataTo(record, values[i])
	}
}
