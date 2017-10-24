package ipfix

import (
	"encoding/binary"
	"fmt"
)

type Template struct {
	ID       int16
	Elements []InformationElement
}

func (t Template) Id() int16 {
	return IDTemplateSet
}

func (t Template) Length() (ret int) {
	ret = 4
	for _, element := range t.Elements {
		ret += element.TemplateSize()
	}
	return
}

func (t Template) SerializeTo(buffer SerializeBuffer) {
	b := buffer.Append(4)
	binary.BigEndian.PutUint16(b[2:], uint16(len(t.Elements)))
	binary.BigEndian.PutUint16(b[0:], uint16(t.ID))
	for _, element := range t.Elements {
		element.SerializeTo(buffer)
	}
}

func (t Template) MakeDataRecord(values ...interface{}) (ret DataRecord) {
	if len(values) != len(t.Elements) {
		panic(fmt.Sprintf("Supplied values (%d) differ from number of information elements (%d)!\n", len(values), len(t.Elements)))
	}
	ret.template = t.ID
	elements := make([]DataRecordElement, len(t.Elements))
	values = values[:len(t.Elements)]
	for i, element := range t.Elements {
		elements[i] = element.MakeDataRecord(values[i])
	}
	ret.elements = elements
	return
}
