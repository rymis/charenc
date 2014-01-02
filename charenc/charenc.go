// Package charenc provides structures and functions to manipulate text encoded with a lot of encodings.
// This package implemented with clean Go and doesn't use iconv.
// Supported encodings includes IBM CP8??, Windows CP12??, MAC, KOI, UTF8/UTF16/UCS2/UCS4 encodings.
package charenc

import (
	"unicode/utf8"
	"strings"
	"errors"
	"bytes"
)

// RuneError is rune representing error value when returned. This contant has imported from utf8 package
const RuneError rune = utf8.RuneError

// RuneDecoder provides DecodeRune and FullRune functions.
// DecodeRune method can be used to decode one rune from string.
// p is original text.
// Return values: resulting rune and len of this character in original text.
// On error returns RuneError, 1
// FullRune method can be used to check if bytes starts from valid rune in this encoding
type RuneDecoder interface {
	DecodeRune(p []byte) (rune, int)
	FullRune(p []byte) bool
}

// RuneEncoder provides interface to encode one rune into string using EncodeRune method.
// p is buffer for rune encoding
// r is rune you want to encode
// Return value: number of bytes written.
// On error method will return -1
type RuneEncoder interface {
	/// Encode one rune into the string. Return bytes written.
	EncodeRune(p []byte, r rune) int
}

// CharacterEncoding interface joins RuneEncoder and RuneDecoder into one interface for encoding and decoding runes from byte arrays
type CharacterEncoding interface {
	RuneDecoder
	RuneEncoder
}

// IgnoreErrors indicates that errors must be ignored
const IgnoreErrors int = 0x0001
// ReplaceErrors indicated that invalid characters will be replaced by '?'
const ReplaceErrors int = 0x0002

/* Unicode encoders/decoders: */
/* List of Unicodes supported by GNU iconv:
UTF-8 
UCS-2, UCS-2BE, UCS-2LE 
UCS-4, UCS-4BE, UCS-4LE 
UTF-16, UTF-16BE, UTF-16LE 
UTF-32, UTF-32BE, UTF-32LE 
*/

// UTF8 is UTF-8 decoder/encoder context
type enc_UTF8 struct { }

func (self enc_UTF8) DecodeRune(p []byte) (rune, int) {
	return utf8.DecodeRune(p)
}

func (self enc_UTF8) FullRune(p []byte) bool {
	return utf8.FullRune(p)
}

func (self enc_UTF8) EncodeRune(p []byte, r rune) int {
	buf := make([]byte, 8)
	rv := utf8.EncodeRune(buf, r)
	if rv > len(p) {
		return -1
	}
	for i := range(buf) {
		p[i] = buf[i]
	}
	return rv
}

func get_UTF8() CharacterEncoding {
	return enc_UTF8{}
}

// Universal reader/writer for 16/32bit integers in Low/Big endian:
func decode_rune(p []byte, be bool, size int) (rune, int) {
	var val uint32 = 0
	if len(p) < size {
		return 0, 0 // End of string
	}

	for i := 0; i < size; i++ {
		if be {
			val = val << 8 + uint32(p[i])
		} else {
			val = val << 8 + uint32(p[size - i - 1])
		}
	}

	return rune(val), size
}

func encode_rune(p []byte, r rune, be bool, size int) int {
	if len(p) < size {
		return -1
	}

	val := uint32(r)
	for i := 0; i < size; i++ {
		if be {
			p[size - i - 1] = byte(val >> uint(i * 8))
		} else {
			p[i] = byte(val >> uint(i * 8))
		}
	}

	return size
}

type enc_UCS2LE struct { }

func (self enc_UCS2LE) DecodeRune(p []byte) (rune, int) {
	return decode_rune(p, false, 2)
}

func (self enc_UCS2LE) FullRune(p []byte) bool {
	return len(p) >= 2
}

func (self enc_UCS2LE) EncodeRune(p []byte, r rune) int {
	return encode_rune(p, r, false, 2)
}

func get_UCS2LE() CharacterEncoding {
	return enc_UCS2LE{}
}

type enc_UCS2BE struct { }

func (self enc_UCS2BE) DecodeRune(p []byte) (rune, int) {
	return decode_rune(p, true, 2)
}

func (self enc_UCS2BE) FullRune(p []byte) bool {
	return len(p) >= 2
}

func (self enc_UCS2BE) EncodeRune(p []byte, r rune) int {
	return encode_rune(p, r, true, 2)
}

func get_UCS2BE() CharacterEncoding {
	return enc_UCS2BE{}
}

type enc_UCS2 struct {
	endian int // Endian can be 0 - not detected (LE), 1 - LE, 2 - BE
}

func (self enc_UCS2) DecodeRune(p []byte) (rune, int) {
	if self.endian == 0 {
		if len(p) > 2 {
			if p[0] == 0xFE && p[1] == 0xFF {
				self.endian = 2
			} else if p[0] == 0xFF && p[1] == 0xFE {
				self.endian = 1
			}
		}

		if (self.endian > 0) {
			r, l := decode_rune(p[2:], self.endian == 2, 2)
			return r, l + 2
		}

		self.endian = 2 // This is default
	}

	return decode_rune(p[2:], self.endian == 2, 2)
}

func (self enc_UCS2) FullRune(p []byte) bool {
	return len(p) >= 2
}

func (self enc_UCS2) EncodeRune(p []byte, r rune) int {
	return encode_rune(p, r, true, 2)
}

func get_UCS2() CharacterEncoding {
	return enc_UCS2{0}
}

type enc_UCS4LE struct { }

func (self enc_UCS4LE) DecodeRune(p []byte) (rune, int) {
	return decode_rune(p, false, 4)
}

func (self enc_UCS4LE) FullRune(p []byte) bool {
	return len(p) >= 4
}

func (self enc_UCS4LE) EncodeRune(p []byte, r rune) int {
	return encode_rune(p, r, false, 4)
}

func get_UCS4LE() CharacterEncoding {
	return enc_UCS4LE{}
}

type enc_UCS4BE struct { }

func (self enc_UCS4BE) DecodeRune(p []byte) (rune, int) {
	return decode_rune(p, true, 4)
}

func (self enc_UCS4BE) FullRune(p []byte) bool {
	return len(p) >= 4
}

func (self enc_UCS4BE) EncodeRune(p []byte, r rune) int {
	return encode_rune(p, r, true, 4)
}

func get_UCS4BE() CharacterEncoding {
	return enc_UCS4BE{}
}

type enc_UCS4 struct {
	endian int // 0 - start, 1 - LE, 2 - BE
}

func (self enc_UCS4) DecodeRune(p []byte) (rune, int) {
	if self.endian == 0 {
		if len(p) > 4 {
			if p[0] == 0 && p[1] == 0 && p[2] == 0xFE && p[3] == 0xFF {
				self.endian = 2
			} else if (p[0] == 0xFF && p[1] == 0xFE && p[2] == 0x00 && p[3] == 0x00) {
				self.endian = 1
			}

			if self.endian > 0 {
				r, l := decode_rune(p[4:], self.endian == 2, 4)
				return r, l + 4
			}

			self.endian = 2 // This is default
		}
	}

	return decode_rune(p, self.endian == 2, 4)
}

func (self enc_UCS4) FullRune(p []byte) bool {
	return len(p) >= 4
}

func (self enc_UCS4) EncodeRune(p []byte, r rune) int {
	return encode_rune(p, r, true, 4)
}

func get_UCS4() CharacterEncoding {
	return enc_UCS4{0}
}

type init_unicode func () CharacterEncoding

var unicode = map[string]init_unicode{
	"UTF-8": get_UTF8,
	"UTF8": get_UTF8,
	"UCS2": get_UCS2,
	"UCS2LE": get_UCS2LE,
	"UCS2BE": get_UCS2BE,
	"UCS4": get_UCS4,
	"UCS4LE": get_UCS4LE,
	"UCS4BE": get_UCS4BE,
}

type bit8 struct {
	id int
}

func (self bit8) DecodeRune(p []byte) (rune, int) {
	if len(p) < 1 {
		return 0, 0
	}

	r := ByteToRune(self.id, p[0])
	if r == 0 {
		return RuneError, 1
	}

	return r, 1
}

func (self bit8) FullRune(p []byte) bool {
	return len(p) >= 1
}

func (self bit8) EncodeRune(p []byte, r rune) int {
	if len(p) < 1 {
		return -1
	}

	b := RuneToByte(self.id, r)
	if b == 0 {
		return -1
	}

	p[0] = b
	return 1
}

func NewRuneDecoder(encoding string) RuneDecoder {
	// 1. Try to create unicode decoder, then 8-bit decoder wrap them on success to error checker
	enc := strings.ToUpper(encoding)
	init, e := unicode[enc]
	if e {
		f := init()
		if f == nil {
			return nil
		}
		return f
	}

	// Not found. Try to find 8-bit encoding
	id := Open8bit(encoding)
	if id >= 0 {
		return bit8{id}
	}

	return nil
}

func NewRuneEncoder(encoding string) RuneEncoder {
	// 1. Try to create unicode decoder, then 8-bit decoder wrap them on success to error checker
	enc := strings.ToUpper(encoding)
	init, e := unicode[enc]
	if e {
		f := init()
		if f == nil {
			return nil
		}
		return f
	}

	// Not found. Try to find 8-bit encoding
	id := Open8bit(encoding)
	if id >= 0 {
		return bit8{id}
	}

	return nil
}

/// Decode bytes to array of runes using specified characters encoding
func DecodeBytes(ctx RuneDecoder, s []byte) ([]rune, error) {
	res := make([]rune, 0)

	for pos := 0; pos < len(s); {
		r, l := ctx.DecodeRune(s[pos:])
		if l < 0 {
			return nil, errors.New("can not decode rune at position " + string(pos))
		}

		res = append(res, r)
	}

	return res, nil
}

/// Encode runes to specified encoding
func EncodeRunes(ctx RuneEncoder, r []rune) ([]byte, error) {
	b := make([]byte, len(r))
	tmpbuf := make([]byte, 8) // I belive there will not be an charset with symbol more then 8 bytes
	buf := bytes.NewBuffer(b)

	for i := range(r) {
		l := ctx.EncodeRune(tmpbuf, r[i])
		if l < 0 {
			return nil, errors.New("can not encode run at position " + string(i))
		}

		buf.Write(tmpbuf[0:l])
	}

	return buf.Bytes(), nil
}

func StringToRunes(ctx RuneDecoder, s string) ([]rune, error) {
	return DecodeBytes(ctx, ([]byte)(s))
}

func RunesToString(ctx RuneEncoder, r []rune) (string, error) {
	res, err := EncodeRunes(ctx, r)

	return string(res), err
}
