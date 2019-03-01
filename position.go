package encoding

type Position struct {
	Line, Column uint64
}

func (p *Position) Advance(b byte) {
	if b == '\n' {
		p.Line++
		p.Column = 0
	} else {
		p.Column++
	}
}
