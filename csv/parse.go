package csv

type State byte

const (
	StateInitial State = iota
	StateUnquotedCell
	StateQuotedCell
	StateQuotedCellQuote
	StateDelimiter
	StateDone
	StateError
)

func (s *State) Advance(t Token) Events {
	var e Events
	*s, e = s.Next(t)
	return e
}

func (s State) Next(t Token) (State, Events) {
	if s < StateError && t < TokenError {
		var trans = stateTransitions[s][t]
		return trans.State, trans.Events
	}
	return StateError, EventError
}

type transition struct {
	State  State
	Events Events
}

var stateTransitions = func() [StateError][TokenError]transition {
	var d [TokenError]transition
	for i := range d {
		d[i] = transition{StateError, EventError}
	}

	var t = [StateError][TokenError]transition{}
	for i := range t {
		t[i] = d
	}

	t[StateInitial][TokenQuote] = transition{StateQuotedCell, EventNewCell}
	t[StateInitial][TokenByte] = transition{StateUnquotedCell, EventByte | EventNewCell}
	t[StateInitial][TokenNewline] = transition{StateInitial, EventNone}
	t[StateInitial][TokenEnd] = transition{StateDone, EventEnd}
	t[StateInitial][TokenDelimiter] = transition{StateDelimiter, EventNewCell}

	t[StateQuotedCell][TokenByte] = transition{StateQuotedCell, EventByte}
	t[StateQuotedCell][TokenNewline] = t[StateQuotedCell][TokenByte]
	t[StateQuotedCell][TokenDelimiter] = t[StateQuotedCell][TokenByte]
	t[StateQuotedCell][TokenQuote] = transition{StateQuotedCellQuote, EventNone}
	t[StateQuotedCell][TokenEnd] = transition{StateError, EventError}

	var end = d
	end[TokenDelimiter] = transition{StateDelimiter, EventNone}
	end[TokenNewline] = transition{StateInitial, EventEndRecord}
	end[TokenEnd] = transition{StateDone, EventEndRecord | EventEnd}

	t[StateUnquotedCell] = end
	t[StateUnquotedCell][TokenByte] = transition{StateUnquotedCell, EventByte}

	t[StateQuotedCellQuote] = end
	t[StateQuotedCellQuote][TokenQuote] = transition{StateQuotedCell, EventByte}

	t[StateDelimiter][TokenQuote] = transition{StateQuotedCell, EventNewCell}
	t[StateDelimiter][TokenByte] = transition{StateUnquotedCell, EventByte | EventNewCell}
	t[StateDelimiter][TokenNewline] = transition{StateInitial, EventNewCell | EventEndRecord}
	t[StateDelimiter][TokenEnd] = transition{StateDone, EventNewCell | EventEndRecord | EventEnd}
	t[StateDelimiter][TokenDelimiter] = transition{StateDelimiter, EventNewCell}
	return t
}()
