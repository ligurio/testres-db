// Package subunit provides a reader of the SubUnit protocol.
// https://github.com/testing-cabal/subunit

package subunit

import (
	"io"
	"log"
	"reflect"
	_ "strconv"
	_ "unsafe"
)

type SubUnitReport struct {
}

type (
	Bytes1 [1]byte
	Bytes2 [2]byte
	Bytes4 [4]byte
	Bytes8 [8]byte
)

/*
Feature bits:

Bit 11	mask 0x0800	Test id present.
Bit 10	mask 0x0400	Routing code present.
Bit 9	mask 0x0200	Timestamp present.
Bit 8	mask 0x0100	Test is 'runnable'.
Bit 7	mask 0x0080	Tags are present.
Bit 6	mask 0x0040	File content is present.
Bit 5	mask 0x0020	File MIME type is present.
Bit 4	mask 0x0010	EOF marker.
Bit 3	mask 0x0008	Must be zero in version 2.
*/
const (
	SUBUNIT_SIGNATURE = 0xB3
	SUBUNIT_VERSION   = 0x02
	PACKET_MAX_LENGTH = 4194303

	FLAG_TEST_ID      = 0x0800
	FLAG_ROUTE_CODE   = 0x0400
	FLAG_TIMESTAMP    = 0x0200
	FLAG_RUNNABLE     = 0x0100
	FLAG_TAGS         = 0x0080
	FLAG_MIME_TYPE    = 0x0020
	FLAG_EOF          = 0x0010
	FLAG_FILE_CONTENT = 0x0040
)

/*
000 - undefined / no test
001 - Enumeration / existence
002 - In progress
003 - Success
004 - Unexpected Success
005 - Skipped
006 - Failed
007 - Expected failure
*/
const (
	UNDEFINED = iota
	ENUMERATION
	INPROGRESS
	SUCCESS
	UNEXPECTED_SUCCESS
	SKIPPED
	FAILED
	EXPECTED_SUCCESS
)

func readByte(reader io.Reader) (v uint8, err error) {
	var data Bytes1
	_, e := reader.Read(data[0:])
	if e != nil {
		return 0, e
	}
	return data[0], nil
}

func readUint16(reader io.Reader) (v uint16, n int, err error) {
	var data Bytes2
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (uint16(data[0]) << 8) | uint16(data[1]), n, nil
}

func readUint32(reader io.Reader) (v uint32, n int, err error) {
	var data Bytes4
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (uint32(data[0]) << 24) | (uint32(data[1]) << 16) | (uint32(data[2]) << 8) | uint32(data[3]), n, nil
}

func readUint64(reader io.Reader) (v uint64, n int, err error) {
	var data Bytes8
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (uint64(data[0]) << 56) | (uint64(data[1]) << 48) | (uint64(data[2]) << 40) | (uint64(data[3]) << 32) | (uint64(data[4]) << 24) | (uint64(data[5]) << 16) | (uint64(data[6]) << 8) | uint64(data[7]), n, nil
}

func readInt16(reader io.Reader) (v int16, n int, err error) {
	var data Bytes2
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (int16(data[0]) << 8) | int16(data[1]), n, nil
}

func readInt32(reader io.Reader) (v int32, n int, err error) {
	var data Bytes4
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (int32(data[0]) << 24) | (int32(data[1]) << 16) | (int32(data[2]) << 8) | int32(data[3]), n, nil
}

func readInt64(reader io.Reader) (v int64, n int, err error) {
	var data Bytes8
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (int64(data[0]) << 56) | (int64(data[1]) << 48) | (int64(data[2]) << 40) | (int64(data[3]) << 32) | (int64(data[4]) << 24) | (int64(data[5]) << 16) | (int64(data[6]) << 8) | int64(data[7]), n, nil
}

// Get the four lowest bits
func lownibble(u8 uint8) uint {
	return uint(u8 & 0xf)
}

// Get the five lowest bits
func lowfive(u8 uint8) uint {
	return uint(u8 & 0x1f)
}

/*
In short the structure of a packet is:

PACKET := SIGNATURE FLAGS PACKET_LENGTH TIMESTAMP? TESTID? TAGS? MIME?
FILECONTENT? ROUTING_CODE? CRC32

In more detail...

Packets are identified by a single byte signature - 0xB3, which is never legal
in a UTF-8 stream as the first byte of a character. 0xB3 starts with the first
bit set and the second not, which is the UTF-8 signature for a continuation
byte. 0xB3 was chosen as 0x73 ('s' in ASCII') with the top two bits replaced by
the 1 and 0 for a continuation byte.

If subunit packets are being embedded in a non-UTF-8 text stream, where 0x73 is
a legal character, consider either recoding the text to UTF-8, or using
subunit's 'file' packets to embed the text stream in subunit, rather than the
other way around.

Following the signature byte comes a 16-bit flags field, which includes a 4-bit
version field - if the version is not 0x2 then the packet cannot be read. It is
recommended to signal an error at this point (e.g. by emitting a synthetic
error packet and returning to the top level loop to look for new packets, or
exiting with an error). If recovery is desired, treat the packet signature as
an opaque byte and scan for a new synchronisation point. NB: Subunit V1 and V2
packets may legitimately included 0xB3 internally, as they are an 8-bit safe
container format, so recovery from this situation may involve an arbitrary
number of false positives until an actual packet is encountered : and even then
it may still be false, failing after passing the version check due to
coincidence.

Flags are stored in network byte order too.
*/

func readSubUnitPacket(reader io.Reader) error {
	return nil
}

func unpack(reader io.Reader, reflected bool) (v reflect.Value, n int, err error) {
	var nbytesread int

	c, e := readByte(reader)
	if e != nil {
		return reflect.Value{}, 0, e
	}
	nbytesread++

	if c == SUBUNIT_SIGNATURE {
		log.Println("SubUnit")
	}

	flags, n, err := readInt16(reader)
	if e != nil {
		return reflect.Value{}, 0, e
	}
	log.Println(n, flags)

	// readSubUnitPacket(reader)

	return reflect.Value{}, 0, e
}

/*
func unpack(reader io.Reader, reflected bool) (v reflect.Value, n int, err error) {
	var retval reflect.Value
	var nbytesread int

	c, e := readByte(reader)
	if e != nil {
		return reflect.Value{}, 0, e
	}
	nbytesread++
	if c < FIXMAP || c >= NEGFIXNUM {
		retval = reflect.ValueOf(int8(c))
	} else if c >= FIXMAP && c <= FIXMAPMAX {
		if reflected {
			retval, n, e = unpackMapReflected(reader, lownibble(c))
		} else {
			retval, n, e = unpackMap(reader, lownibble(c))
		}
		nbytesread += n
		if e != nil {
			return reflect.Value{}, nbytesread, e
		}
		nbytesread += n
	} else if c >= FIXARRAY && c <= FIXARRAYMAX {
		if reflected {
			retval, n, e = unpackArrayReflected(reader, lownibble(c))
		} else {
			retval, n, e = unpackArray(reader, lownibble(c))
		}
		nbytesread += n
		if e != nil {
			return reflect.Value{}, nbytesread, e
		}
		nbytesread += n
	} else if c >= FIXRAW && c <= FIXRAWMAX {
		data := make([]byte, lowfive(c))
		n, e := reader.Read(data)
		nbytesread += n
		if e != nil {
			return reflect.Value{}, nbytesread, e
		}
		retval = reflect.ValueOf(data)
	} else {
		switch c {
		case NIL:
			retval = reflect.ValueOf(nil)
		case FALSE:
			retval = reflect.ValueOf(false)
		case TRUE:
			retval = reflect.ValueOf(true)
		case FLOAT:
			data, n, e := readUint32(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(*(*float32)(unsafe.Pointer(&data)))
		case DOUBLE:
			data, n, e := readUint64(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(*(*float64)(unsafe.Pointer(&data)))
		case UINT8:
			data, e := readByte(reader)
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(uint8(data))
			nbytesread++
		case UINT16:
			data, n, e := readUint16(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(data)
		case UINT32:
			data, n, e := readUint32(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(data)
		case UINT64:
			data, n, e := readUint64(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(data)
		case INT8:
			data, e := readByte(reader)
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(int8(data))
			nbytesread++
		case INT16:
			data, n, e := readInt16(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(data)
		case INT32:
			data, n, e := readInt32(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(data)
		case INT64:
			data, n, e := readInt64(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(data)
		case RAW16:
			nbytestoread, n, e := readUint16(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			data := make([]byte, nbytestoread)
			n, e = reader.Read(data)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(data)
		case RAW32:
			nbytestoread, n, e := readUint32(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			data := make(Bytes, nbytestoread)
			n, e = reader.Read(data)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			retval = reflect.ValueOf(data)
		case ARRAY16:
			nelemstoread, n, e := readUint16(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			if reflected {
				retval, n, e = unpackArrayReflected(reader, uint(nelemstoread))
			} else {
				retval, n, e = unpackArray(reader, uint(nelemstoread))
			}
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
		case ARRAY32:
			nelemstoread, n, e := readUint32(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			if reflected {
				retval, n, e = unpackArrayReflected(reader, uint(nelemstoread))
			} else {
				retval, n, e = unpackArray(reader, uint(nelemstoread))
			}
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
		case MAP16:
			nelemstoread, n, e := readUint16(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			if reflected {
				retval, n, e = unpackMapReflected(reader, uint(nelemstoread))
			} else {
				retval, n, e = unpackMap(reader, uint(nelemstoread))
			}
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
		case MAP32:
			nelemstoread, n, e := readUint32(reader)
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
			if reflected {
				retval, n, e = unpackMapReflected(reader, uint(nelemstoread))
			} else {
				retval, n, e = unpackMap(reader, uint(nelemstoread))
			}
			nbytesread += n
			if e != nil {
				return reflect.Value{}, nbytesread, e
			}
		default:
			panic("unsupported code: " + strconv.Itoa(int(c)))
		}
	}
	return retval, nbytesread, nil
}

// Reads a value from the reader, unpack and returns it.
func Unpack(reader io.Reader) (v reflect.Value, n int, err error) {
	return unpack(reader, false)
}

// Reads unpack a value from the reader, unpack and returns it.  When the
// value is an array or map, leaves the elements wrapped by corresponding
// wrapper objects defined in reflect package.
func UnpackReflected(reader io.Reader) (v reflect.Value, n int, err error) {
	return unpack(reader, true)
}
*/

func NewParser(r io.Reader) (*SubUnitReport, error) {

	return nil, nil
}
