package ipfix

import (
	"encoding/binary"
	"log"
	"math"
	"net"
	"time"
)

// Datatypes according to RFC7011

type Type int

const (
	OctetArray Type = iota
	Unsigned8
	Unsigned16
	Unsigned32
	Unsigned64
	Signed8
	Signed16
	Signed32
	Signed64
	Float32
	Float64
	Boolean
	MacAddress
	String
	DateTimeSeconds
	DateTimeMilliseconds
	DateTimeMicroseconds
	DateTimeNanoseconds
	Ipv4Address
	Ipv6Address
)

const VariableLength uint16 = 65535

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
}

func NameToType(x []byte) Type {
	switch string(x) {
	case "octetArray":
		return OctetArray
	case "unsigned8":
		return Unsigned8
	case "unsigned16":
		return Unsigned16
	case "unsigned32":
		return Unsigned32
	case "unsigned64":
		return Unsigned64
	case "signed8":
		return Signed8
	case "signed16":
		return Signed16
	case "signed32":
		return Signed32
	case "signed64":
		return Signed64
	case "float32":
		return Float32
	case "float64":
		return Float64
	case "boolean":
		return Boolean
	case "macAddress":
		return MacAddress
	case "string":
		return String
	case "dateTimeSeconds":
		return DateTimeSeconds
	case "dateTimeMilliseconds":
		return DateTimeMilliseconds
	case "dateTimeMicroseconds":
		return DateTimeMicroseconds
	case "dateTimeNanoseconds":
		return DateTimeNanoseconds
	case "ipv4Address":
		return Ipv4Address
	case "ipv6Address":
		return Ipv6Address
	}
	log.Panicf("Unknown type %s\n", x)
	return -1
}

func (t Type) String() string {
	switch t {
	case OctetArray:
		return "octetArray"
	case Unsigned8:
		return "unsigned8"
	case Unsigned16:
		return "unsigned16"
	case Unsigned32:
		return "unsigned32"
	case Unsigned64:
		return "unsigned64"
	case Signed8:
		return "signed8"
	case Signed16:
		return "signed16"
	case Signed32:
		return "signed32"
	case Signed64:
		return "signed64"
	case Float32:
		return "float32"
	case Float64:
		return "float64"
	case Boolean:
		return "boolean"
	case MacAddress:
		return "macAddress"
	case String:
		return "string"
	case DateTimeSeconds:
		return "dateTimeSeconds"
	case DateTimeMilliseconds:
		return "dateTimeMilliseconds"
	case DateTimeMicroseconds:
		return "dateTimeMicroseconds"
	case DateTimeNanoseconds:
		return "dateTimeNanoseconds"
	case Ipv4Address:
		return "ipv4Address"
	case Ipv6Address:
		return "ipv6Address"
	}
	return "unknownType"
}

//Seconds between NTP and Unix epoch
const NTPToUnix uint32 = 0x83AA7E80

func MakeDataRecordElement(t Type, value interface{}, length int) DataRecordElement {
	switch t {
	case OctetArray:
		return &OctetArrayDataRecord{
			BaseDataRecord{length},
			value.([]byte),
			length == 0,
		}
	case Unsigned8:
		return &Unsigned8DataRecord{
			value.(uint8),
		}
	case Unsigned16:
		if length == 0 {
			length = 2
		}
		return &Unsigned16DataRecord{
			BaseDataRecord{length},
			value.(uint16),
		}
	case Unsigned32:
		if length == 0 {
			length = 4
		}
		return &Unsigned32DataRecord{
			BaseDataRecord{length},
			value.(uint32),
		}
	case Unsigned64:
		if length == 0 {
			length = 8
		}
		return &Unsigned64DataRecord{
			BaseDataRecord{length},
			value.(uint64),
		}
	case Signed8:
		return &Signed8DataRecord{
			value.(int8),
		}
	case Signed16:
		if length == 0 {
			length = 2
		}
		return &Signed16DataRecord{
			BaseDataRecord{length},
			value.(int16),
		}
	case Signed32:
		if length == 0 {
			length = 4
		}
		return &Signed32DataRecord{
			BaseDataRecord{length},
			value.(int32),
		}
	case Signed64:
		if length == 0 {
			length = 8
		}
		return &Signed64DataRecord{
			BaseDataRecord{length},
			value.(int64),
		}
	case Float32:
		return &Float32DataRecord{
			value.(float32),
		}
	case Float64:
		if length == 0 {
			length = 8
		}
		return &Float64DataRecord{
			BaseDataRecord{length},
			value.(float64),
		}
	case Boolean:
		return &BooleanDataRecord{
			value.(bool),
		}
	case MacAddress:
		return &MacAddressDataRecord{
			value.(net.HardwareAddr),
		}
	case String:
		return &StringDataRecord{
			BaseDataRecord{length},
			value.(string),
			length == 0,
		}
	case DateTimeSeconds:
		return &DateTimeSecondsDataRecord{
			value.(time.Time),
		}
	case DateTimeMilliseconds:
		return &DateTimeMillisecondsDataRecord{
			value.(time.Time),
		}
	case DateTimeMicroseconds:
		return &DateTimeMicrosecondsDataRecord{
			value.(time.Time),
		}
	case DateTimeNanoseconds:
		return &DateTimeNanosecondsDataRecord{
			value.(time.Time),
		}
	case Ipv4Address:
		return &Ipv4AddressDataRecord{
			value.(net.IP),
		}
	case Ipv6Address:
		return &Ipv6AddressDataRecord{
			value.(net.IP),
		}
	}
	return nil
}

type BaseDataRecord struct {
	length int
}

func (bd *BaseDataRecord) Length() int { return bd.length }

type OctetArrayDataRecord struct {
	BaseDataRecord
	value  []byte
	varlen bool
}

func (d *OctetArrayDataRecord) SerializeTo(buffer SerializeBuffer) {
	var b, assign, clear []byte
	if d.varlen {
		len := len(d.value)
		if len < 255 {
			b = buffer.Append(len + 1)
			_ = b[1]
			b[0] = uint8(len)
			assign = b[1:]
		} else {
			b = buffer.Append(len + 3)
			_ = b[3]
			b[0] = 0xff
			binary.BigEndian.PutUint16(b[1:3], uint16(len))
			assign = b[3:]
		}
	} else {
		assign = buffer.Append(d.length)
		clear = assign[len(d.value):]
	}
	copy(assign, d.value)
	for i := range clear {
		clear[i] = 0
	}
}

func (d *OctetArrayDataRecord) Length() int {
	if d.varlen {
		len := len(d.value)
		if len < 255 {
			return len + 1
		}
		return len + 3
	}
	return d.length
}

func varEncodeInt(buffer SerializeBuffer, v uint64, length int) {
	switch length {
	case 1:
		b := buffer.Append(1)
		b[0] = byte(v)
	case 2:
		b := buffer.Append(2)
		binary.BigEndian.PutUint16(b, uint16(v))
	case 3:
		b := buffer.Append(3)
		_ = b[2]
		b[0] = byte(v >> 16)
		b[1] = byte(v >> 8)
		b[2] = byte(v)
	case 4:
		b := buffer.Append(2)
		binary.BigEndian.PutUint32(b, uint32(v))
	case 5:
		b := buffer.Append(5)
		_ = b[4]
		b[0] = byte(v >> 32)
		b[1] = byte(v >> 24)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 8)
		b[4] = byte(v)
	case 6:
		b := buffer.Append(6)
		_ = b[5]
		b[0] = byte(v >> 40)
		b[1] = byte(v >> 32)
		b[2] = byte(v >> 24)
		b[3] = byte(v >> 16)
		b[4] = byte(v >> 8)
		b[5] = byte(v)
	case 7:
		b := buffer.Append(7)
		_ = b[6]
		b[0] = byte(v >> 48)
		b[1] = byte(v >> 40)
		b[2] = byte(v >> 32)
		b[3] = byte(v >> 24)
		b[4] = byte(v >> 16)
		b[5] = byte(v >> 8)
		b[6] = byte(v)
	case 8:
		b := buffer.Append(8)
		binary.BigEndian.PutUint64(b, v)
	default:
		log.Panicln("Illegal encoding length for integers. Must be 1-8. Was ", length)
	}
}

type Unsigned8DataRecord struct {
	value uint8
}

func (d *Unsigned8DataRecord) SerializeTo(buffer SerializeBuffer) {
	b := buffer.Append(1)
	b[0] = byte(d.value)
}

func (d *Unsigned8DataRecord) Length() int {
	return 1
}

type Unsigned16DataRecord struct {
	BaseDataRecord
	value uint16
}

func (d *Unsigned16DataRecord) SerializeTo(buffer SerializeBuffer) {
	varEncodeInt(buffer, uint64(d.value), d.length)
}

type Unsigned32DataRecord struct {
	BaseDataRecord
	value uint32
}

func (d *Unsigned32DataRecord) SerializeTo(buffer SerializeBuffer) {
	varEncodeInt(buffer, uint64(d.value), d.length)
}

type Unsigned64DataRecord struct {
	BaseDataRecord
	value uint64
}

func (d *Unsigned64DataRecord) SerializeTo(buffer SerializeBuffer) {
	varEncodeInt(buffer, d.value, d.length)
}

type Signed8DataRecord struct {
	value int8
}

func (d *Signed8DataRecord) SerializeTo(buffer SerializeBuffer) {
	b := buffer.Append(1)
	b[0] = byte(d.value)
}

func (d *Signed8DataRecord) Length() int {
	return 1
}

type Signed16DataRecord struct {
	BaseDataRecord
	value int16
}

func (d *Signed16DataRecord) SerializeTo(buffer SerializeBuffer) {
	varEncodeInt(buffer, uint64(d.value), d.length)
}

type Signed32DataRecord struct {
	BaseDataRecord
	value int32
}

func (d *Signed32DataRecord) SerializeTo(buffer SerializeBuffer) {
	varEncodeInt(buffer, uint64(d.value), d.length)
}

type Signed64DataRecord struct {
	BaseDataRecord
	value int64
}

func (d *Signed64DataRecord) SerializeTo(buffer SerializeBuffer) {
	varEncodeInt(buffer, uint64(d.value), d.length)
}

type Float32DataRecord struct {
	value float32
}

func (d *Float32DataRecord) SerializeTo(buffer SerializeBuffer) {
	b := buffer.Append(4)
	bits := math.Float32bits(d.value)
	binary.LittleEndian.PutUint32(b, bits)
}

func (d *Float32DataRecord) Length() int {
	return 4
}

type Float64DataRecord struct {
	BaseDataRecord
	value float64
}

func (d *Float64DataRecord) SerializeTo(buffer SerializeBuffer) {
	switch d.length {
	case 4:
		b := buffer.Append(4)
		bits := math.Float32bits(float32(d.value))
		binary.LittleEndian.PutUint32(b, bits)
	case 8:
		b := buffer.Append(8)
		bits := math.Float64bits(d.value)
		binary.LittleEndian.PutUint64(b, bits)
	default:
		log.Panicln("Illegal encoding length for float64. Must be 4, 8. Was ", d.length)
	}
}

type BooleanDataRecord struct {
	value bool
}

func (d *BooleanDataRecord) SerializeTo(buffer SerializeBuffer) {
	b := buffer.Append(1)
	if d.value {
		b[0] = 1
	} else {
		b[0] = 2
	}
}

func (d *BooleanDataRecord) Length() int {
	return 1
}

type MacAddressDataRecord struct {
	value net.HardwareAddr
}

func (d *MacAddressDataRecord) SerializeTo(buffer SerializeBuffer) {
	b := buffer.Append(6)
	copy(b, d.value)
}

func (d *MacAddressDataRecord) Length() int {
	return 6
}

type StringDataRecord struct {
	BaseDataRecord
	value  string
	varlen bool
}

func (d *StringDataRecord) SerializeTo(buffer SerializeBuffer) {
	var b, assign, clear []byte
	if d.varlen {
		len := len(d.value)
		if len < 255 {
			b = buffer.Append(len + 1)
			_ = b[1]
			b[0] = uint8(len)
			assign = b[1:]
		} else {
			b = buffer.Append(len + 3)
			_ = b[3]
			b[0] = 0xff
			binary.BigEndian.PutUint16(b[1:3], uint16(len))
			assign = b[3:]
		}
	} else {
		assign = buffer.Append(d.length)
		clear = assign[len(d.value):]
	}
	copy(assign, d.value) // FIXME: fixup string if cutting short results in utf8 errors
	for i := range clear {
		clear[i] = 0
	}
}

func (d *StringDataRecord) Length() int {
	if d.varlen {
		len := len(d.value)
		if len < 255 {
			return len + 1
		}
		return len + 3
	}
	return d.length
}

type DateTimeSecondsDataRecord struct {
	value time.Time
}

func (d *DateTimeSecondsDataRecord) SerializeTo(buffer SerializeBuffer) {
	b := buffer.Append(4)
	binary.BigEndian.PutUint32(b, uint32(d.value.Unix()))
}

func (d *DateTimeSecondsDataRecord) Length() int {
	return 4
}

type DateTimeMillisecondsDataRecord struct {
	value time.Time
}

func (d *DateTimeMillisecondsDataRecord) SerializeTo(buffer SerializeBuffer) {
	b := buffer.Append(8)
	binary.BigEndian.PutUint64(b, uint64(d.value.UnixNano()/int64(time.Millisecond)))
}

func (d *DateTimeMillisecondsDataRecord) Length() int {
	return 8
}

type DateTimeMicrosecondsDataRecord struct {
	value time.Time
}

func (d *DateTimeMicrosecondsDataRecord) SerializeTo(buffer SerializeBuffer) {
	//NTP epoch as 32bit seconds + 32bit fraction (~244ps)
	//-> get time in Unixpoch seconds, add ntp epoch to unix epoch offset
	//-> shift nanoseconds 32 bit to the left divide by nano (1e9)
	//according to RFC7011 the last 11 bits should be zero (0xFFFFF800) to get micro seconds (~.447 microsecond resolution)
	b := buffer.Append(8)
	binary.BigEndian.PutUint32(b, uint32(d.value.Unix())+NTPToUnix)
	binary.BigEndian.PutUint32(b, uint32((uint64(d.value.Nanosecond())<<32)/1e9)&0xFFFFF800)
}

func (d *DateTimeMicrosecondsDataRecord) Length() int {
	return 8
}

type DateTimeNanosecondsDataRecord struct {
	value time.Time
}

func (d *DateTimeNanosecondsDataRecord) SerializeTo(buffer SerializeBuffer) {
	//same as microseconds, but without truncation
	b := buffer.Append(8)
	binary.BigEndian.PutUint32(b, uint32(d.value.Unix())+NTPToUnix)
	binary.BigEndian.PutUint32(b, uint32((uint64(d.value.Nanosecond())<<32)/1e9))
}

func (d *DateTimeNanosecondsDataRecord) Length() int {
	return 8
}

type Ipv4AddressDataRecord struct {
	value net.IP
}

func (d *Ipv4AddressDataRecord) SerializeTo(buffer SerializeBuffer) {
	b := buffer.Append(4)
	copy(b, d.value)
}

func (d *Ipv4AddressDataRecord) Length() int {
	return 4
}

type Ipv6AddressDataRecord struct {
	value net.IP
}

func (d *Ipv6AddressDataRecord) SerializeTo(buffer SerializeBuffer) {
	b := buffer.Append(16)
	copy(b, d.value)
}

func (d *Ipv6AddressDataRecord) Length() int {
	return 16
}
