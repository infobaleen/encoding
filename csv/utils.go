package csv

import (
	"bytes"
	"io"
)

func (d *Decoder) StringsAll() ([][]string, error) {
	var out [][]string
	var bytBuffer bytes.Buffer
	for {
		var record, err = d.stringsRecord(&bytBuffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		out = append(out, record)
	}
	return out, nil
}

func (d *Decoder) StringsRecord() ([]string, error) {
	var bytBuffer bytes.Buffer
	return d.stringsRecord(&bytBuffer)
}

func (d *Decoder) stringsRecord(bytBuffer *bytes.Buffer) ([]string, error) {
	var err = d.NextRecord()
	if err != nil {
		return nil, err
	}

	var out []string
	for {
		err = d.NextCell()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		bytBuffer.Reset()
		_, err = bytBuffer.ReadFrom(d)
		if err != nil {
			return nil, err
		}
		out = append(out, bytBuffer.String())
	}
	return out, nil
}
