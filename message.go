package ipfix

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type MessageStream struct {
	w                 io.Writer
	buffer            SerializeBuffer
	length            []byte
	time              []byte
	record            Record
	sequence          uint32
	observationID     uint32
	templates         []*Template
	currentSet        Set
	currentDataRecord BufferedDataRecord
	dirty             bool
}

func MakeMessageStream(w io.Writer, mtu uint16, observationID uint32) (ret *MessageStream) {
	if mtu == 0 {
		mtu = 65535
	}
	buffer := MakeSerializeBuffer(int(mtu))
	ret = &MessageStream{
		w:                 w,
		buffer:            buffer,
		observationID:     observationID,
		currentSet:        MakeSet(buffer),
		currentDataRecord: MakeBufferedDataRecord(4096),
	}
	return
}

func (m *MessageStream) startMessage() {
	b := m.buffer.Append(16)
	_ = b[15]
	b[0] = 0
	b[1] = 0x0a
	m.length = b[2:4]
	m.time = b[4:8]
	m.dirty = true
	binary.BigEndian.PutUint32(b[8:12], uint32(m.sequence))
	binary.BigEndian.PutUint32(b[12:16], uint32(m.observationID))
}

func (m *MessageStream) sendRecord(record Record, now interface{}) (err error) {
	if !m.dirty {
		m.startMessage()
	}
	for {
		err = m.currentSet.AppendRecord(record)
		if err == nil {
			if record.Id() >= 256 {
				m.sequence++
			}
			return
		}
		if ipfixerr, ok := err.(IPFIXError); ok {
			switch {
			case ipfixerr.BufferFull():
				m.Finalize(now)
				m.startMessage()
			case ipfixerr.RecordTypeMismatch():
				m.currentSet.Finalize()
			default:
				return
			}
		} else {
			return
		}
	}
}

func (m *MessageStream) AddTemplate(now interface{}, elements ...InformationElement) (id int, err error) {
	id = len(m.templates) + 256
	template := Template{int16(id), elements}
	if err = m.sendRecord(template, now); err == nil {
		m.templates = append(m.templates, &template)
	}
	return
}

func (m *MessageStream) SendData(now interface{}, template int, data ...interface{}) (err error) {
	id := template - 256
	if id < 0 || id >= len(m.templates) {
		panic(fmt.Sprintf("Unknown template id %d\n", template))
	}
	t := m.templates[id]
	if t == nil {
		panic(fmt.Sprintf("Unknown template id %d\n", template))
	}
	t.AssignDataRecord(&m.currentDataRecord, data...)
	return m.sendRecord(&m.currentDataRecord, now)
}

func (m *MessageStream) Finalize(now interface{}) (err error) {
	if !m.dirty {
		return nil
	}
	m.currentSet.Finalize()
	binary.BigEndian.PutUint16(m.length, uint16(m.buffer.Length()))
	switch v := now.(type) {
	case time.Time:
		binary.BigEndian.PutUint32(m.time, uint32(v.Unix()))
	case DateTimeSeconds:
		binary.BigEndian.PutUint32(m.time, uint32(v))
	case DateTimeMilliseconds:
		binary.BigEndian.PutUint32(m.time, uint32(v/1e3))
	case DateTimeMicroseconds:
		binary.BigEndian.PutUint32(m.time, uint32(v/1e6))
	case DateTimeNanoseconds:
		binary.BigEndian.PutUint32(m.time, uint32(v/1e9))
	}
	if err = m.buffer.Finalize(m.w); err == nil {
		m.dirty = false
	}
	return
}
