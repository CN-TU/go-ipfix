package ipfix

import (
	"encoding/binary"
	"math"
	"net"
	"time"
)

// Type is a datatype according to RFC7011
type Type int

// DateTimeSeconds represents time in units of seconds from 00:00 UTC, Januray 1, 1970 according to RFC5102.
type DateTimeSeconds uint64

// DateTimeMilliseconds represents time in units of milliseconds from 00:00 UTC, Januray 1, 1970 according to RFC5102.
type DateTimeMilliseconds uint64

// DateTimeMicroseconds represents time in units of microseconds from 00:00 UTC, Januray 1, 1970 according to RFC5102.
type DateTimeMicroseconds uint64

// DateTimeNanoseconds represents time in units of nanoseconds from 00:00 UTC, Januray 1, 1970 according to RFC5102.
type DateTimeNanoseconds uint64

const (
	// OctetArrayType as defined by RFC7011
	OctetArrayType Type = iota
	// Unsigned8Type as defined by RFC7011
	Unsigned8Type
	// Unsigned16Type as defined by RFC7011
	Unsigned16Type
	// Unsigned32Type as defined by RFC7011
	Unsigned32Type
	// Unsigned64Type as defined by RFC7011
	Unsigned64Type
	// Signed8Type as defined by RFC7011
	Signed8Type
	// Signed16Type as defined by RFC7011
	Signed16Type
	// Signed32Type as defined by RFC7011
	Signed32Type
	// Signed64Type as defined by RFC7011
	Signed64Type
	// Float32Type as defined by RFC7011
	Float32Type
	// Float64Type as defined by RFC7011
	Float64Type
	// BooleanType as defined by RFC7011
	BooleanType
	// MacAddressType as defined by RFC7011
	MacAddressType
	// StringType as defined by RFC7011
	StringType
	// DateTimeSecondsType as defined by RFC7011
	DateTimeSecondsType
	// DateTimeMillisecondsType as defined by RFC7011
	DateTimeMillisecondsType
	// DateTimeMicrosecondsType as defined by RFC7011
	DateTimeMicrosecondsType
	// DateTimeNanosecondsType as defined by RFC7011
	DateTimeNanosecondsType
	// Ipv4AddressType as defined by RFC7011
	Ipv4AddressType
	// Ipv6AddressType as defined by RFC7011
	Ipv6AddressType
	// BasicListType as defined by RFC7011
	BasicListType
	// IllegalType is an undefined type
	IllegalType = -1
)

// VariableLength is the variable length specifier as defined by RFC7011
const VariableLength uint16 = 65535

// DefaultSize can be used to look up the default size of an ipfix type
var DefaultSize = [...]uint16{
	VariableLength,
	1,
	2,
	4,
	8,
	1,
	2,
	4,
	8,
	4,
	8,
	1,
	6,
	VariableLength,
	4,
	8,
	8,
	8,
	4,
	16,
	VariableLength,
}

// NameToType converts the given textual representation of a type to the ipfix type. Returns IllegalType if the type is not recognised.
func NameToType(x []byte) Type {
	switch string(x) {
	case "octetArray":
		return OctetArrayType
	case "unsigned8":
		return Unsigned8Type
	case "unsigned16":
		return Unsigned16Type
	case "unsigned32":
		return Unsigned32Type
	case "unsigned64":
		return Unsigned64Type
	case "signed8":
		return Signed8Type
	case "signed16":
		return Signed16Type
	case "signed32":
		return Signed32Type
	case "signed64":
		return Signed64Type
	case "float32":
		return Float32Type
	case "float64":
		return Float64Type
	case "boolean":
		return BooleanType
	case "macAddress":
		return MacAddressType
	case "string":
		return StringType
	case "dateTimeSeconds":
		return DateTimeSecondsType
	case "dateTimeMilliseconds":
		return DateTimeMillisecondsType
	case "dateTimeMicroseconds":
		return DateTimeMicrosecondsType
	case "dateTimeNanoseconds":
		return DateTimeNanosecondsType
	case "ipv4Address":
		return Ipv4AddressType
	case "ipv6Address":
		return Ipv6AddressType
	}
	return IllegalType
}

func (t Type) String() string {
	switch t {
	case OctetArrayType:
		return "octetArray"
	case Unsigned8Type:
		return "unsigned8"
	case Unsigned16Type:
		return "unsigned16"
	case Unsigned32Type:
		return "unsigned32"
	case Unsigned64Type:
		return "unsigned64"
	case Signed8Type:
		return "signed8"
	case Signed16Type:
		return "signed16"
	case Signed32Type:
		return "signed32"
	case Signed64Type:
		return "signed64"
	case Float32Type:
		return "float32"
	case Float64Type:
		return "float64"
	case BooleanType:
		return "boolean"
	case MacAddressType:
		return "macAddress"
	case StringType:
		return "string"
	case DateTimeSecondsType:
		return "dateTimeSeconds"
	case DateTimeMillisecondsType:
		return "dateTimeMilliseconds"
	case DateTimeMicrosecondsType:
		return "dateTimeMicroseconds"
	case DateTimeNanosecondsType:
		return "dateTimeNanoseconds"
	case Ipv4AddressType:
		return "ipv4Address"
	case Ipv6AddressType:
		return "ipv6Address"
	case BasicListType:
		return "basicList"
	case IllegalType:
		return "<bad>"
	}
	return "unknownType"
}

//Seconds between NTP and Unix epoch
const ntp2Unix uint32 = 0x83AA7E80

func (t Type) serializeDataTo(buffer scratchBuffer, value interface{}, length int) (int, error) {
	switch t {
	case OctetArrayType, Ipv4AddressType, Ipv6AddressType, MacAddressType, StringType:
		return serializeOctetArrayTo(buffer, t, value, length)
	case Unsigned8Type, Unsigned16Type, Unsigned32Type, Unsigned64Type, Signed8Type, Signed16Type, Signed32Type, Signed64Type, BooleanType:
		return serializeIntegerTo(buffer, t, value, length)
	case Float32Type, Float64Type:
		return serializeFloatTo(buffer, t, value, length)
	case DateTimeSecondsType, DateTimeMillisecondsType, DateTimeMicrosecondsType, DateTimeNanosecondsType:
		return serializeDateTimeTo(buffer, t, value, length)
	}
	return 0, IllegalTypeError(t)
}

func serializeOctetArrayTo(buffer scratchBuffer, t Type, value interface{}, length int) (int, error) {
	var val []byte
	switch v := value.(type) {
	case string:
		val = []byte(v)
	case []byte:
		val = v
	case net.IP:
		val = []byte(v)
	case net.HardwareAddr:
		val = []byte(v)
	case nil:
		// val is already nil
	default:
		return 0, ConversionError{t, value}
	}
	if length == 0 {
		length = int(DefaultSize[t])
	}
	if length == int(VariableLength) {
		length = len(val)
		written := length
		var assign []byte
		if length == 0 {
			b, err := buffer.append(length + 1)
			if err != nil {
				return 0, err
			}
			b[0] = 0
			return 1, nil
		} else if length < 255 {
			written++
			b, err := buffer.append(length + 1)
			if err != nil {
				return 0, err
			}
			_ = b[1]
			b[0] = uint8(length)
			assign = b[1:]
		} else {
			written += 3
			b, err := buffer.append(length + 3)
			if err != nil {
				return 0, err
			}
			_ = b[2]
			b[0] = 0xff
			binary.BigEndian.PutUint16(b[1:3], uint16(length))
			assign = b[3:]
		}
		copy(assign, val)
		return written, nil
	}
	if len(val) == length {
		b, err := buffer.append(length)
		if err != nil {
			return 0, err
		}
		copy(b, val)
		return length, nil
	}
	if t == Ipv4AddressType || t == Ipv6AddressType || t == MacAddressType {
		if val == nil {
			tmp, err := buffer.append(length)
			if err != nil {
				return 0, err
			}
			for i := range tmp {
				tmp[i] = 0
			}
			return length, nil
		}
		return 0, ConversionError{t, value}
	}
	var clear []byte
	assign, err := buffer.append(length)
	if err != nil {
		return 0, err
	}
	copy(assign, val)
	clear = assign[len(val):]
	for i := range clear {
		clear[i] = 0
	}
	return length, nil
}

func serializeIntegerTo(buffer scratchBuffer, t Type, value interface{}, length int) (int, error) {
	var val uint64
	switch v := value.(type) {
	case float64:
		val = uint64(v)
	case float32:
		val = uint64(v)
	case int64:
		val = uint64(v)
	case int32:
		val = uint64(v)
	case int16:
		val = uint64(v)
	case int8:
		val = uint64(v)
	case int:
		val = uint64(v)
	case uint64:
		val = v
	case uint32:
		val = uint64(v)
	case uint16:
		val = uint64(v)
	case uint8:
		val = uint64(v)
	case uint:
		val = uint64(v)
	case nil:
		// val already 0
	case bool:
		if v {
			val = 1
		} else {
			val = 2
		}
	default:
		return 0, ConversionError{t, value}
	}
	if length == 0 {
		length = int(DefaultSize[t])
	}
	switch length {
	case 1:
		b, err := buffer.append(1)
		if err != nil {
			return 0, err
		}
		b[0] = byte(val)
		return 1, nil
	case 2:
		b, err := buffer.append(2)
		if err != nil {
			return 0, err
		}
		binary.BigEndian.PutUint16(b, uint16(val))
		return 2, nil
	case 3:
		b, err := buffer.append(3)
		if err != nil {
			return 0, err
		}
		_ = b[2]
		b[0] = byte(val >> 16)
		b[1] = byte(val >> 8)
		b[2] = byte(val)
		return 3, nil
	case 4:
		b, err := buffer.append(4)
		if err != nil {
			return 0, err
		}
		binary.BigEndian.PutUint32(b, uint32(val))
		return 4, nil
	case 5:
		b, err := buffer.append(5)
		if err != nil {
			return 0, err
		}
		_ = b[4]
		b[0] = byte(val >> 32)
		b[1] = byte(val >> 24)
		b[2] = byte(val >> 16)
		b[3] = byte(val >> 8)
		b[4] = byte(val)
		return 5, nil
	case 6:
		b, err := buffer.append(6)
		if err != nil {
			return 0, err
		}
		_ = b[5]
		b[0] = byte(val >> 40)
		b[1] = byte(val >> 32)
		b[2] = byte(val >> 24)
		b[3] = byte(val >> 16)
		b[4] = byte(val >> 8)
		b[5] = byte(val)
		return 6, nil
	case 7:
		b, err := buffer.append(7)
		if err != nil {
			return 0, err
		}
		_ = b[6]
		b[0] = byte(val >> 48)
		b[1] = byte(val >> 40)
		b[2] = byte(val >> 32)
		b[3] = byte(val >> 24)
		b[4] = byte(val >> 16)
		b[5] = byte(val >> 8)
		b[6] = byte(val)
		return 7, nil
	case 8:
		b, err := buffer.append(8)
		if err != nil {
			return 0, err
		}
		binary.BigEndian.PutUint64(b, val)
		return 8, nil
	default:
		return 0, SizeError{t, length}
	}
}

func serializeFloatTo(buffer scratchBuffer, t Type, value interface{}, length int) (int, error) {
	if t == Float32Type {
		var val float32
		switch v := value.(type) {
		case float64:
			val = float32(v)
		case float32:
			val = v
		case int64:
			val = float32(v)
		case int32:
			val = float32(v)
		case int16:
			val = float32(v)
		case int8:
			val = float32(v)
		case int:
			val = float32(v)
		case uint64:
			val = float32(v)
		case uint32:
			val = float32(v)
		case uint16:
			val = float32(v)
		case uint8:
			val = float32(v)
		case uint:
			val = float32(v)
		case nil:
			// val already 0
		case bool:
			if v {
				val = 1
			} else {
				val = 2
			}
		default:
			return 0, ConversionError{t, value}
		}
		b, err := buffer.append(4)
		if err != nil {
			return 0, err
		}
		bits := math.Float32bits(val)
		binary.BigEndian.PutUint32(b, bits)
		return 4, nil
	}
	var val float64
	switch v := value.(type) {
	case float64:
		val = v
	case float32:
		val = float64(v)
	case int64:
		val = float64(v)
	case int32:
		val = float64(v)
	case int16:
		val = float64(v)
	case int8:
		val = float64(v)
	case int:
		val = float64(v)
	case uint64:
		val = float64(v)
	case uint32:
		val = float64(v)
	case uint16:
		val = float64(v)
	case uint8:
		val = float64(v)
	case uint:
		val = float64(v)
	case nil:
		// val already 0
	case bool:
		if v {
			val = 1
		} else {
			val = 2
		}
	default:
		return 0, ConversionError{t, value}
	}
	switch length {
	case 4:
		b, err := buffer.append(4)
		if err != nil {
			return 0, err
		}
		bits := math.Float32bits(float32(val))
		binary.BigEndian.PutUint32(b, bits)
	case 8:
		b, err := buffer.append(8)
		if err != nil {
			return 0, err
		}
		bits := math.Float64bits(val)
		binary.BigEndian.PutUint64(b, bits)
	default:
		return 0, SizeError{t, length}
	}
	return length, nil
}

func serializeDateTimeTo(buffer scratchBuffer, t Type, value interface{}, length int) (int, error) {
	var seconds, nanoseconds uint64
	switch v := value.(type) {
	case time.Time:
		seconds = uint64(v.Unix())
		nanoseconds = uint64(v.Nanosecond())
	case DateTimeMilliseconds:
		seconds = uint64(v) / 1e3
		nanoseconds = (uint64(v) % 1e3) * 1e6
	case DateTimeMicroseconds:
		seconds = uint64(v) / 1e6
		nanoseconds = (uint64(v) % 1e6) * 1e3
	case DateTimeNanoseconds:
		seconds = uint64(v) / 1e9
		nanoseconds = uint64(v) % 1e9
	case uint64:
		seconds = v / 1e9
		nanoseconds = v % 1e9
	case int64:
		seconds = uint64(v) / 1e9
		nanoseconds = uint64(v) % 1e9
	case float64:
		seconds = uint64(v) / 1e9
		nanoseconds = uint64(v) % 1e9
	case nil:
		// val already 0
	default:
		return 0, ConversionError{t, value}
	}
	switch t {
	case DateTimeSecondsType:
		b, err := buffer.append(4)
		if err != nil {
			return 0, err
		}
		binary.BigEndian.PutUint32(b, uint32(seconds))
		return 4, nil
	case DateTimeMillisecondsType:
		b, err := buffer.append(8)
		if err != nil {
			return 0, err
		}
		binary.BigEndian.PutUint64(b, uint64(seconds*1e3+nanoseconds/1e6))
		return 8, nil
	case DateTimeMicrosecondsType:
		//NTP epoch as 32bit seconds + 32bit fraction (~244ps)
		//-> get time in Unixpoch seconds, add ntp epoch to unix epoch offset
		//-> shift nanoseconds 32 bit to the left divide by nano (1e9)
		//according to RFC7011 the last 11 bits should be zero (0xFFFFF800) to get micro seconds (~.447 microsecond resolution)
		b, err := buffer.append(8)
		if err != nil {
			return 0, err
		}
		_ = b[7]
		binary.BigEndian.PutUint32(b[:4], uint32(seconds)+ntp2Unix)
		binary.BigEndian.PutUint32(b[4:8], uint32((nanoseconds<<32)/1e9)&0xFFFFF800)
		return 8, nil
	case DateTimeNanosecondsType:
		b, err := buffer.append(8)
		if err != nil {
			return 0, err
		}
		_ = b[7]
		binary.BigEndian.PutUint32(b[:4], uint32(seconds)+ntp2Unix)
		binary.BigEndian.PutUint32(b[4:8], uint32((nanoseconds<<32)/1e9))
		return 8, nil
	}
	return 0, ConversionError{t, value}
}
