package csv

import (
	"math/bits"
)

// Events stores 0 or more events emitted by the parser. The numerical values order event types by scope,
// importance and order of appearance: NewCell < Byte < EndRecord < End < Error. This makes discarding lesser events
// easy, while waiting for a particular event type, without risk of missing more important events that need handling.
type Events byte

const (
	EventNewCell Events = 1 << iota
	EventByte
	EventEndRecord
	EventEnd
	EventError
	EventNone = 0
)

func (e Events) Next() Events {
	return Events(1 << uint(bits.TrailingZeros8(byte(e))))
}

func (e Events) Contains(c Events) bool {
	return e&c == c
}

func (e *Events) Clear(c Events) {
	*e &^= c
}

func (e Events) String() string {
	var r = "["
	for e != EventNone {
		if len(r) > 1 {
			r += ", "
		}
		switch e.Next() {
		case EventNewCell:
			r += "EventNewCell"
		case EventByte:
			r += "EventByte"
		case EventEndRecord:
			r += "EventEndRecord"
		case EventEnd:
			r += "EventEnd"
		case EventError:
			r += "EventError"
		}
		e.Clear(e.Next())
	}
	return r + "]"
}
