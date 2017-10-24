package ipfix

import (
	"encoding/binary"
	"io"
	"log"
	"time"
)

type MessageStream struct {
	w             io.Writer
	buffer        SerializeBuffer
	length        []byte
	time          []byte
	record        Record
	sequence      uint32
	observationID uint32
	templates     []*Template
	currentSet    *Set
	dirty         bool
}

func MakeMessageStream(w io.Writer, mtu uint16, observationID uint32) (ret *MessageStream) {
	if mtu == 0 {
		mtu = 65535
	}
	ret = &MessageStream{
		w:             w,
		buffer:        MakeSerializeBuffer(int(mtu)),
		observationID: observationID,
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

func (m *MessageStream) createSet(record Record, now time.Time) (err error) {
	var newSet *Set
	if !m.dirty {
		m.startMessage()
	}
	if newSet, err = MakeSet(record, m.buffer); err != nil {
		if ipfixerr, ok := err.(IPFIXError); ok && ipfixerr.BufferFull() {
			// Buffer full -> finalize and try again
			if err = m.Finalize(now); err != nil {
				return
			}
			m.startMessage()
			if newSet, err = MakeSet(record, m.buffer); err != nil {
				// ok if this happens, we lost
				return
			}
		}
	}
	m.currentSet = newSet
	return
}

func (m *MessageStream) sendRecord(record Record, now time.Time) (err error) {
	if m.currentSet == nil {
		if err = m.createSet(record, now); err != nil {
			return
		}
		m.sequence++
		return nil
	}
	if err = m.currentSet.AppendRecord(record); err != nil {
		if ipfixerr, ok := err.(IPFIXError); ok {
			switch {
			case ipfixerr.BufferFull():
				m.currentSet.Finalize()
				m.Finalize(now)
				m.createSet(record, now)
			case ipfixerr.RecordTypeMismatch():
				m.currentSet.Finalize()
				m.createSet(record, now)
			default:
				return
			}
		}
	}
	m.sequence++
	return nil
}

func (m *MessageStream) AddTemplate(now time.Time, elements ...InformationElement) (id int, err error) {
	id = len(m.templates) + 256
	template := Template{int16(id), elements}
	if err = m.sendRecord(template, now); err == nil {
		m.templates = append(m.templates, &template)
	}
	return
}

func (m *MessageStream) SendData(now time.Time, template int, data ...interface{}) (err error) {
	id := template - 256
	if id < 0 || id >= len(m.templates) {
		log.Panicf("Unknown template id %d\n", template)
	}
	t := m.templates[id]
	if t == nil {
		log.Panicf("Unknown template id %d\n", template)
	}
	return m.sendRecord(t.MakeDataRecord(data...), now)
}

func (m *MessageStream) Finalize(now time.Time) (err error) {
	if !m.dirty {
		return nil
	}
	m.currentSet.Finalize()
	binary.BigEndian.PutUint16(m.length, uint16(m.buffer.Length()))
	binary.BigEndian.PutUint32(m.time, uint32(now.Unix()))
	if err = m.buffer.Finalize(m.w); err == nil {
		m.dirty = false
	}
	return
}
