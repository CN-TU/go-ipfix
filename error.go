package ipfix

import "fmt"

type IPFIXError interface {
	error
	BufferFull() bool
	RecordTypeMismatch() bool
}

type BufferFullError int

func (e BufferFullError) Error() string            { return "Buffer full - would need " + string(e) }
func (e BufferFullError) BufferFull() bool         { return true }
func (e BufferFullError) RecordTypeMismatch() bool { return false }

type RecordMismatchError struct {
	a, b int16
}

func (e *RecordMismatchError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Record type mismatch: %d != %d!", e.a, e.b)
}
func (e *RecordMismatchError) BufferFull() bool         { return false }
func (e *RecordMismatchError) RecordTypeMismatch() bool { return true }

type IllegalSetError struct {
	Err string
}

func (e *IllegalSetError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Err
}
func (e *IllegalSetError) BufferFull() bool         { return false }
func (e *IllegalSetError) RecordTypeMismatch() bool { return false }

type IllegalEncodingError string

func (e IllegalEncodingError) Error() string            { return string(e) }
func (e IllegalEncodingError) BufferFull() bool         { return false }
func (e IllegalEncodingError) RecordTypeMismatch() bool { return false }
