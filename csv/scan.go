package csv

type ScanToken byte

const (
	TokenNewline ScanToken = iota
	TokenQuote
	TokenDelimiter
	TokenByte
	TokenEnd
)

func (t ScanToken) String() string {
	var names = []string{"TokenNewline", "TokenQuote", "TokenDelimiter", "TokenByte", "TokenEnd"}
	if int(t) < len(names) {
		return names[t]
	}
	return "TokenUnknown"
}

func Scan(next, delimiter byte) ScanToken {
	switch next {
	case '\n', '\r':
		return TokenNewline
	case '"':
		return TokenQuote
	case delimiter:
		return TokenDelimiter
	default:
		return TokenByte
	}
}
