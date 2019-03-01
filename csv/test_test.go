package csv

import (
	"fmt"
	"github.com/infobaleen/encoding"
	"io"
	"reflect"
	"strings"
	"testing"
)

func run(d *Decoder, in string, out [][]string, err error) string {
	if d == nil {
		d = NewDecoder(nil)
	}
	d.Reader = strings.NewReader(in)
	var testOut, testErr = d.StringsAll()
	if !reflect.DeepEqual(testErr, err) {
		return fmt.Sprintf("Error mismatch:\ngot  %#v\nwant %#v", testErr, err)
	} else if !reflect.DeepEqual(testOut, out) {
		return fmt.Sprintf("Output mismatch:\ngot  %#v\nwant %#v", testOut, out)
	}
	return ""
}

func TestBasic(t *testing.T) {
	var r = run(nil, "a,s,d,f", [][]string{{"a", "s", "d", "f"}}, nil)
	if r != "" {
		t.Error(r)
	}
}

func TestEmptyFields(t *testing.T) {
	var r = run(nil, "a,\n,s\n,,", [][]string{{"a", ""}, {"", "s"}, {"", "", ""}}, nil)
	if r != "" {
		t.Error(r)
	}
}

func TestNewline(t *testing.T) {
	var r = run(nil, "a,s\nd,f\r\ng,h\n", [][]string{{"a", "s"}, {"d", "f"}, {"g", "h"}}, nil)
	if r != "" {
		t.Error(r)
	}
}

func TestMultiNewline(t *testing.T) {
	var r = run(nil, "a,s\n\nd,f", [][]string{{"a", "s"}, {"d", "f"}}, nil)
	if r != "" {
		t.Error(r)
	}
}

func TestQuoted(t *testing.T) {
	var r = run(nil, "\"a\",s\n\"d\nf\",\"g\"\"\"", [][]string{{"a", "s"}, {"d\nf", "g\""}}, nil)
	if r != "" {
		t.Error(r)
	}
}

func TestTabDelimiter(t *testing.T) {
	var d = NewDecoder(nil)
	d.Delimiter = '\t'
	var r = run(d, "a\ts\nd\tf", [][]string{{"a", "s"}, {"d", "f"}}, nil)
	if r != "" {
		t.Error(r)
	}
}

func TestQuoteInUnquotedCell(t *testing.T) {
	var r = run(nil, "a,s\",d,f", nil, DecoderError{nil, encoding.Position{0, 4}, '"'})
	if r != "" {
		t.Error(r)
	}
}

func TestNextSkips(t *testing.T) {
	for _, skipCells := range []bool{true, false} {
		var d = NewDecoder(strings.NewReader("a\n\n\n,\n\ns,d,\n,f,g,h\n"))
		var err error
		var recordCount = 0
		for err = d.NextRecord(); err == nil; err = d.NextRecord() {
			recordCount++
			if !skipCells {
				var cellCount = 0
				for err = d.NextCell(); err == nil; err = d.NextCell() {
					cellCount++
				}
				if err != io.EOF {
					t.Errorf("wrong error\ngot %#v\nwant io.EOF", err)
				}
				if cellCount != recordCount {
					t.Errorf("wrong cell count\ngot  %d\nwant %d", cellCount, recordCount)
				}
			}
		}
		if err != io.EOF {
			t.Errorf("wrong error\ngot %#v\nwant io.EOF", err)
		}
		if want := 4; recordCount != want {
			t.Errorf("wrong record count\ngot  %d\nwant %d", recordCount, want)
		}
	}
}

type repeatReader struct {
	s       string
	o, i, n int
}

func (r *repeatReader) Read(p []byte) (int, error) {
	var n int
	for r.i < r.n && n < len(p) {
		var m = copy(p[n:], r.s[r.o:])
		r.o += m
		n += m
		if r.o == len(r.s) {
			r.i++
			r.o = 0
		}
	}
	if n == 0 && r.i == r.n {
		return 0, io.EOF
	}
	return n, nil
}

func bench(b *testing.B, data string) {
	var r = repeatReader{s: data, n: b.N}
	var d = NewDecoder(&r)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	var err error
	for err == nil {
		_, err = d.NextByte()
		if err == io.EOF {
			err = d.NextCell()
			if err == io.EOF {
				err = d.NextRecord()
			}
		}
	}
	b.StopTimer()
	if err != io.EOF {
		b.Fatal(err)
	}
}

func BenchmarkSmall(b *testing.B) {
	bench(b, "a,s,d,f\n")
}

func BenchmarkLarge(b *testing.B) {
	var data = strings.Repeat(strings.Repeat("a", 100)+",", 10)
	data = data[:len(data)-1] + "\n"
	bench(b, data)
}

func BenchmarkOverhead(b *testing.B) {
	var data = strings.Repeat("a", 100)
	var r = repeatReader{s: data, n: b.N}
	var buf = make([]byte, 1)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	var err error
	for err == nil {
		_, err = r.Read(buf)
	}
	b.StopTimer()
}
