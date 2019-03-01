package csv

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/infobaleen/encoding"
)

func decodeByte(s State, b, delimiter byte) (State, Events) {
	var t = Byte2Token(b, delimiter)
	var e = s.Advance(t)
	return s, e
}

// DecodeReader parses the reader contents until a byte is found that triggers events. If EOF is encountered and parsed
// without error, StateDone, EventEnd and no error are returned. Otherwise StateError, EventError and io.EOF are
// returned.
// Other read errors are returned without affecting the parse state or any events.
func decodeReader(s State, r io.Reader, delimiter byte, p *encoding.Position) (State, Events, byte, error) {
	var e Events
	var b byte

	// Avoid allocation. We should probably open a compiler issue for this.
	var sr = reflect.SliceHeader{uintptr(unsafe.Pointer(&b)), 1, 1}
	var sb = *(*[]byte)(unsafe.Pointer(&sr))

	for e == EventNone {
		var _, err = r.Read(sb)
		if err == io.EOF {
			e = s.Advance(TokenEnd)
			if e != EventError {
				continue
			}
		}
		if err != nil {
			return s, e, 0, err
		}
		p.Advance(b)
		s, e = decodeByte(s, b, delimiter)
	}
	return s, e, b, nil
}

type Decoder struct {
	Delimiter byte

	Reader    io.Reader
	ReadError error
	Position  encoding.Position

	Parser State
	Events Events
	Next   byte
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{Reader: r, Delimiter: ',', Parser: StateInitial, Events: EventEndRecord}
}

func (d *Decoder) advance() error {
	var err = d.error()
	if err == nil && d.Events == EventNone {
		d.Parser, d.Events, d.Next, d.ReadError = decodeReader(d.Parser, d.Reader, d.Delimiter, &d.Position)
		err = d.error()
	}
	return err
}

func (d *Decoder) NextRecord() error {
	for d.Events.Next() != EventEndRecord {
		d.Events.Clear(EventNewCell | EventByte)
		if err := d.advance(); err != nil {
			return err
		}
	}
	d.Events.Clear(EventEndRecord)
	if err := d.advance(); err != nil {
		return err
	}
	if d.Events.Next() == EventNewCell {
		return nil
	}
	return io.EOF
}

func (d *Decoder) NextCell() error {
	for {
		if err := d.advance(); err != nil {
			return err
		}
		if d.Events.Contains(EventNewCell) {
			d.Events.Clear(EventNewCell)
			return nil
		}
		d.Events.Clear(EventByte)
		if d.Events != EventNone {
			return io.EOF
		}
	}
}

func (d *Decoder) NextByte() (byte, error) {
	if d.ReadError == nil && d.Events == EventNone {
		d.Parser, d.Events, d.Next, d.ReadError = decodeReader(d.Parser, d.Reader, d.Delimiter, &d.Position)
	}
	if err := d.error(); err != nil {
		return 0, err
	}
	if d.Events.Next() != EventByte {
		return 0, io.EOF
	}
	d.Events.Clear(EventByte)
	return d.Next, nil
}

func (d *Decoder) Read(p []byte) (int, error) {
	for i := range p {
		var b, err = d.NextByte()
		if err != nil {
			return i, err
		}
		p[i] = b
	}
	return len(p), nil
}

func (d *Decoder) error() error {
	if d.Events == EventError {
		return DecoderError{nil, d.Position, d.Next}
	} else if d.ReadError != nil {
		return DecoderError{d.ReadError, d.Position, 0}
	}
	return nil
}

type DecoderError struct {
	ReaderError error
	Position    encoding.Position
	Last        byte
}

func (e DecoderError) Cause() error {
	if e.ReaderError == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return e.ReaderError
}

func (e DecoderError) Error() string {
	var s = fmt.Sprintf("line %d, column %d", e.Position.Line+1, e.Position.Column)
	if c := e.Cause(); c != nil {
		return fmt.Sprintf("%s: %v", s, c)
	}
	return fmt.Sprintf("%s: unexpected byte %q", s, e.Last)
}
