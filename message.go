package ipfix

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

// MessageStream represents an ipfix message stream.
type MessageStream struct {
	w                 io.Writer
	buffer            scratchBuffer
	length            []byte
	time              []byte
	record            record
	sequence          uint32
	observationID     uint32
	templates         []*template
	currentSet        set
	currentDataRecord recordBuffer
	mtu               int
	dirty             bool
}

// MakeMessageStream initializes a new message stream, which writes to the given writer and uses the given mtu size.
// The observationID is used as the observation id in the ipfix messages.
func MakeMessageStream(w io.Writer, mtu uint16, observationID uint32) (ret *MessageStream, err error) {
	if mtu == 0 {
		mtu = 65535
	} else if mtu < 28 {
		return nil, errors.New("mtu must be at least 28")
	}
	buffer := makeBasicBuffer(int(mtu))
	ret = &MessageStream{
		w:                 w,
		buffer:            buffer,
		observationID:     observationID,
		currentSet:        makeSet(buffer),
		currentDataRecord: makeRecordBuffer(int(mtu)),
		mtu:               int(mtu),
	}
	return
}

func (m *MessageStream) startMessage() error {
	b, err := m.buffer.append(16)
	if err != nil {
		return err
	}
	_ = b[15]
	b[0] = 0
	b[1] = 0x0a
	m.length = b[2:4]
	m.time = b[4:8]
	m.dirty = true
	binary.BigEndian.PutUint32(b[8:12], uint32(m.sequence))
	binary.BigEndian.PutUint32(b[12:16], uint32(m.observationID))
	return nil
}

func (m *MessageStream) sendRecord(rec record, now interface{}) (err error) {
	if !m.dirty {
		m.startMessage()
	}
RETRY:
	err = m.currentSet.appendRecord(rec)
	if err == nil {
		if rec.id() >= 256 {
			m.sequence++
		}
		return
	}
	//	fmt.Println(err)
	if ipfixerr, ok := err.(ipfixError); ok {
		switch {
		case ipfixerr.bufferFull():
			if m.buffer.length() == 16 {
				return RecordTooBigError{16 + rec.length(), m.mtu}
			}
			m.Flush(now)
			m.startMessage()
			goto RETRY
		case ipfixerr.recordTypeMismatch():
			m.currentSet.finalize()
			goto RETRY
		}
	}
	return
}

// AddTemplate adds the given InformationElement as a new template. now must be the current or exported
// time either as a time.Time value or as one of the provieded ipfix time types. A template id is
// returned that can be used with SendData. In case of error an error value is provided.
func (m *MessageStream) AddTemplate(now interface{}, elements ...InformationElement) (id int, err error) {
	id = len(m.templates) + 256
	newTemplate := template{int16(id), elements}
	if err = m.sendRecord(newTemplate, now); err == nil {
		m.templates = append(m.templates, &newTemplate)
	}
	return
}

// SendTemplate resend an existing template by id.
// It returns an error if the template can not be found or the send failed.
func (m *MessageStream) SendTemplate(now interface{}, id int) (err error) {
	if id >= len(m.templates) {
		return UnknownTemplateError(id)
	}
	newTemplate := m.templates[id]
	err = m.sendRecord(newTemplate, now)
	return
}

// SendData sends the given values for the given template id (Can be allocated with AddTemplate).
// now must be the current or exported time either as a time.Time value or as one of the provieded ipfix time types.
// Template InformationElements and given data types must match. Numeric types are converted automatically.
// In case of error an error is returned.
func (m *MessageStream) SendData(now interface{}, template int, data ...interface{}) (err error) {
	id := template - 256
	if id < 0 || id >= len(m.templates) {
		return UnknownTemplateError(template)
	}
	t := m.templates[id]
	if t == nil {
		return UnknownTemplateError(template)
	}
	err = t.assignDataRecord(&m.currentDataRecord, data...)
	if err != nil {
		return
	}
	return m.sendRecord(&m.currentDataRecord, now)
}

// Flush must be called before the underlying writer is closed. This function finishes and flushes
// eventual not yet finalized messages. This does not flush the underlying buffer!
func (m *MessageStream) Flush(now interface{}) (err error) {
	if !m.dirty {
		return nil
	}
	m.currentSet.finalize()
	binary.BigEndian.PutUint16(m.length, uint16(m.buffer.length()))
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
	if err = m.buffer.finalize(m.w); err == nil {
		m.dirty = false
	}
	return
}
