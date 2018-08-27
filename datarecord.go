package ipfix

type recordBuffer struct {
	basicBuffer
	template int16
}

func makeRecordBuffer(num int) recordBuffer {
	return recordBuffer{basicBuffer: make([]byte, 0, num)}
}

func (b recordBuffer) id() int16 {
	return b.template
}

func (b *recordBuffer) serializeTo(basicBuffer scratchBuffer) error {
	buf, err := basicBuffer.append(len(b.basicBuffer))
	if err != nil {
		return err
	}
	copy(buf, b.basicBuffer)
	b.basicBuffer = b.basicBuffer[:0]
	return nil
}
