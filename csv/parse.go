package csv

type ParseState func(ScanToken) (ParseState, ParseEvents)

func (s *ParseState) Advance(t ScanToken) ParseEvents {
	var events = EventError
	if *s != nil {
		*s, events = (*s)(t)
	}
	return events
}

func StateInitial(t ScanToken) (ParseState, ParseEvents) {
	switch t {
	case TokenQuote:
		return StateQuotedCell, EventNewCell
	case TokenByte:
		return StateUnquotedCell, EventByte + EventNewCell
	case TokenNewline:
		return StateInitial, EventNone
	case TokenEnd:
		return StateDone, EventEnd
	case TokenDelimiter:
		return StateDelimiter, EventNewCell
	default:
		return nil, EventError
	}
}

func StateUnquotedCell(t ScanToken) (ParseState, ParseEvents) {
	if t == TokenByte {
		return StateUnquotedCell, EventByte
	}
	return StateEndCell(t)
}

func StateQuotedCell(t ScanToken) (ParseState, ParseEvents) {
	switch t {
	case TokenQuote:
		return StateQuotedCellQuote, EventNone
	case TokenEnd:
		return nil, EventError
	}
	return StateQuotedCell, EventByte
}

func StateQuotedCellQuote(t ScanToken) (ParseState, ParseEvents) {
	if t == TokenQuote {
		return StateQuotedCell, EventByte
	}
	return StateEndCell(t)
}

func StateEndCell(t ScanToken) (ParseState, ParseEvents) {
	switch t {
	case TokenDelimiter:
		return StateDelimiter, EventNone
	case TokenNewline:
		return StateInitial, EventEndRecord
	case TokenEnd:
		return StateDone, EventEndRecord + EventEnd
	default:
		return nil, EventError
	}
}

func StateDelimiter(t ScanToken) (ParseState, ParseEvents) {
	switch t {
	case TokenQuote:
		return StateQuotedCell, EventNewCell
	case TokenByte:
		return StateUnquotedCell, EventByte + EventNewCell
	case TokenNewline:
		return StateInitial, EventNewCell + EventEndRecord
	case TokenEnd:
		return StateDone, EventNewCell + EventEndRecord + EventEnd
	case TokenDelimiter:
		return StateDelimiter, EventNewCell
	default:
		return nil, EventError
	}
}

func StateDone(_ ScanToken) (ParseState, ParseEvents) {
	return nil, EventError
}
