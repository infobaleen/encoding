package csv

import (
	"math/bits"
)

// ParseEvents stores 0 or more events emitted by the parser. The numerical values order event types by scope,
// importance and order of appearance: NewCell < Byte < EndRecord < End < Error. This makes discarding lesser events
// easy, while waiting for a particular event type, without risk of missing more important events that need handling.
type ParseEvents byte

const (
	EventNewCell ParseEvents = 1 << iota
	EventByte
	EventEndRecord
	EventEnd
	EventError
	EventNone = 0
)

func (e ParseEvents) Next() ParseEvents {
	return ParseEvents(1 << uint(bits.TrailingZeros8(byte(e))))
}

func (e ParseEvents) Contains(c ParseEvents) bool {
	return e&c == c
}

func (e *ParseEvents) ConsumeNext() ParseEvents {
	var r = e.Next()
	e.Clear(r)
	return r
}

func (e *ParseEvents) Clear(c ParseEvents) {
	*e &^= c
}

func (e *ParseEvents) ClearAll() {
	*e = 0
}

// ClearUntil clears all events with lower importance than the specified one.
func (e *ParseEvents) ClearUntil(c ParseEvents) {
	var t = uint(bits.TrailingZeros8(byte(c)))
	*e = (*e >> t) << t
}

func (e ParseEvents) String() string {
	var r = "["
	for e != EventNone {
		if len(r) > 1 {
			r += ", "
		}
		switch e.ConsumeNext() {
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
	}
	return r + "]"
}
