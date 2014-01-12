
package charenc

import (
	"io"
	"errors"
)

type RuneReader struct {
	reader  io.Reader
	decoder RuneDecoder
	buf []byte
	pos, cnt int
	erract  int
	err error
}

func NewRuneReader(reader io.Reader, decoder RuneDecoder, erract int) *RuneReader {
	res := new(RuneReader)

	res.decoder = decoder
	res.reader  = reader
	res.pos = 0
	res.cnt = 0
	res.buf = make([]byte, 256)
	res.err = nil
	res.erract = erract

	return res
}

func (self *RuneReader) ReadRunes(p []rune) (int, error) {
	// Read from input if we don't have enought bytes:
	if self.err == nil && self.cnt - self.pos < 8 {
		if self.pos > 0 {
			copy(self.buf, self.buf[self.pos: self.cnt])
		}
		self.cnt -= self.pos
		self.pos = 0

		n, e := self.reader.Read(self.buf[self.cnt:])
		if n == 0 {
			self.err = e
		}

		self.cnt += n
	}

	if self.cnt == 0 { // This if is only for optimization: actually we can remove this lines
		return 0, self.err
	}

	// Write output:
	var pos int = 0
	for pos = 0; pos < len(p); {
		if !self.decoder.FullRune(self.buf[self.pos:]) {
			return pos, self.err // self.err may be EOF and it is Ok for us
		}

		r, cnt := self.decoder.DecodeRune(self.buf[self.pos:])
		if cnt < 0 {
			if self.erract == ReplaceErrors {
				p[pos] = '?'
				pos++
				self.pos++
			} else if self.erract == IgnoreErrors {
				r = 0
				self.pos++
			} else {
				self.err = errors.New("Unicode decoder failed")
				return pos, self.err
			}
		} else {
			p[pos] = r
			pos++
			self.pos += cnt
		}
	}

	return pos, nil
}

// Reader takes io.Reader in first encoding and provides io.Reader in second one
type Reader struct {
	reader  io.Reader
	decoder RuneDecoder
	encoder RuneEncoder
	buf []byte
	pos, cnt int
	erract  int
	err error
}

func NewReader(reader io.Reader, decoder RuneDecoder, encoder RuneEncoder, erract int) *Reader {
	res := new(Reader)

	res.decoder = decoder
	res.encoder = encoder
	res.reader  = reader
	res.pos = 0
	res.cnt = 0
	res.buf = make([]byte, 256)
	res.err = nil
	res.erract = erract

	return res
}

func GetReader(reader io.Reader, from_charset, to_charset string, erract int) *Reader {
	decoder := NewRuneDecoder(from_charset)
	encoder := NewRuneEncoder(to_charset)

	if decoder == nil || encoder == nil {
		return nil
	}

	return NewReader(reader, decoder, encoder, erract)
}

func (self *Reader) Read(p []byte) (int, error) {
	var charbuf []byte = make([]byte, 8)
	// Read from input if we don't have enought bytes:
	if self.err == nil && self.cnt - self.pos < len(charbuf) {
		if self.pos > 0 {
			copy(self.buf, self.buf[self.pos: self.cnt])
		}
		self.cnt -= self.pos
		self.pos = 0

		n, e := self.reader.Read(self.buf[self.cnt:])
		if n == 0 {
			self.err = e
		}

		self.cnt += n
	}

	if self.cnt == 0 { // This if is only for optimization: actually we can remove this lines
		return 0, self.err
	}

	// Write output:
	var pos int = 0
	for pos = 0; pos < len(p); {
		if !self.decoder.FullRune(self.buf[self.pos:]) {
			return pos, self.err // self.err may be EOF and it is Ok for us
		}

		r, cnt := self.decoder.DecodeRune(self.buf[self.pos:])
		if cnt < 0 {
			if self.erract == ReplaceErrors {
				r = '?'
				cnt = 1
			} else if self.erract == IgnoreErrors {
				r = 0
			} else {
				self.err = errors.New("Unicode decoder failed")
				return pos, self.err
			}
		}

		var ocnt int = 0
		if cnt >= 0 {
			ocnt = self.encoder.EncodeRune(charbuf, r)
			if ocnt < 0 {
				self.err = errors.New("Unicode encoder failed");
				return pos, self.err
			}
			if ocnt + pos >= len(p) { // Buffer is full
				return pos, nil
			}
			// Copy charbuf to buf:
			copy(p[pos:], charbuf[:ocnt])
		}
		if cnt < 0 {
			cnt = 1
		}

		pos += ocnt
		self.pos += cnt
		if self.pos >= self.cnt {
			self.pos = 0
			self.cnt = 0
			return pos, self.err // At this point input buffer empty so we can return error (EOF for example)
		}
	}

	return pos, nil
}

type RuneWriter struct {
	writer io.Writer
	encoder RuneEncoder
	err error
	erract  int
}

func NewRuneWriter(writer io.Writer, encoder RuneEncoder, erract int) *RuneWriter {
	r := new(RuneWriter)
	r.writer = writer
	r.encoder = encoder
	r.err = nil
	r.erract = erract

	return r
}

func (self RuneWriter) WriteRunes(p []rune) (int, error) {
	buf := make([]byte, 256)
	pos := 0
	var e error
	l := 0
	for i := range(p) {
		if pos + 8 > len(buf) {
			l, e = self.writer.Write(buf[:pos])
			if l < pos {
				return i, e
			}
			pos = 0
		}

		l = self.encoder.EncodeRune(buf[pos:], p[i])
		if l < 0 {
			if self.erract == ReplaceErrors {
				buf[pos] = '?'
				pos++
			} else if self.erract == IgnoreErrors {
			} else {
				return i, errors.New("Can not encode character")
			}
		} else {
			pos += l
		}
	}

	if pos > 0 {
		l, e = self.writer.Write(buf[:pos])
		if l < pos {
			return len(p), e // It is not true but we can not determine how many characters had written
		}
	}

	return len(p), nil
}

type Writer struct {
	writer *RuneWriter
	buf []byte
	pos   int
}

func NewWriter(writer io.Writer, encoder RuneEncoder, decoder RuneDecoder) *Writer {
	res := new(Writer)
	res.buf = nil
	res.writer = writer
	res.encoder = encoder
	res.decoder = decoder
	res.err = nil

	return res
}

func (self *Writer) Write(p []byte) (int, error) {

	return 0, io.EOF
}




