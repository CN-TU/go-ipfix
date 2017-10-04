package ipfix

type DataRecordElement interface {
	Length() int
	SerializeTo(buffer SerializeBuffer)
}

type DataRecord struct {
	template int16
	elements []DataRecordElement
}

func (d DataRecord) Id() int16 {
	return d.template
}

func (d DataRecord) Length() (ret int) {
	for _, element := range d.elements {
		ret += element.Length()
	}
	return
}

func (d DataRecord) SerializeTo(buffer SerializeBuffer) {
	for _, element := range d.elements {
		element.SerializeTo(buffer)
	}
}
