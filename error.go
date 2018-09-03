package ipfix

import (
	"fmt"
	"reflect"
)

// RecordTooBigError indicates that the template or data set was too big for the mtu
type RecordTooBigError struct {
	size, mtu int
}

func (e RecordTooBigError) Error() string {
	return fmt.Sprintf("ipfix: Template/Data set too big. Would need %d but maximum mtu is %d", e.size, e.mtu)
}

// TemplateMismatchError indicates that the number of required information elements did not match the number of passed data values.
type TemplateMismatchError struct {
	used, required int
}

func (e TemplateMismatchError) Error() string {
	return fmt.Sprintf("ipfix: Template does not match given data. Wanted %d elements, got %d elements", e.required, e.used)
}

// BasicListMismatchError indicates that the number of required information elements did not match the number of passed data values.
type BasicListMismatchError struct {
	used, required int
}

func (e BasicListMismatchError) Error() string {
	return fmt.Sprintf("ipfix: BasicList does not match given data. Wanted %d elements, got %d elements", e.required, e.used)
}

// UnknownTemplateError indicates that the given template id is unknown
type UnknownTemplateError int

func (e UnknownTemplateError) Error() string {
	return fmt.Sprintf("ipfix: Template id %d unknown", int(e))
}

// IllegalTypeError indicates that the given type is not known
type IllegalTypeError Type

func (e IllegalTypeError) Error() string {
	return fmt.Sprintf("ipfix: Illegal type %d", int(e))
}

// ConversionError indicates that the type of the given value can't be converted to the given ipfix type
type ConversionError struct {
	want Type
	have interface{}
}

func (e ConversionError) Error() string {
	return fmt.Sprintf("ipfix: Can't convert %s to %s", reflect.TypeOf(e.have), e.want)
}

// SizeError indicates that the given size is illegal for the given type
type SizeError struct {
	t     Type
	given int
}

func (e SizeError) Error() string {
	return fmt.Sprintf("ipfix: Illegal size %d for ipfix type %s", e.given, e.t)
}

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
