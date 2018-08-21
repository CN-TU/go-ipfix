package ipfix

import "fmt"

type ipfixError interface {
	error
	bufferFull() bool
	recordTypeMismatch() bool
}

type bufferFullError int

func (e bufferFullError) Error() string            { return "Buffer full - would need " + fmt.Sprint(int(e)) }
func (e bufferFullError) bufferFull() bool         { return true }
func (e bufferFullError) recordTypeMismatch() bool { return false }

type recordMismatchError struct {
	a, b int16
}

func (e *recordMismatchError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Record type mismatch: %d != %d!", e.a, e.b)
}
func (e *recordMismatchError) bufferFull() bool         { return false }
func (e *recordMismatchError) recordTypeMismatch() bool { return true }

type illegalSetError struct {
	Err string
}

func (e *illegalSetError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Err
}
func (e *illegalSetError) bufferFull() bool         { return false }
func (e *illegalSetError) recordTypeMismatch() bool { return false }
