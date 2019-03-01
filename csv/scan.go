package csv

type Token byte

const (
	TokenNewline Token = iota
	TokenQuote
	TokenDelimiter
	TokenByte
	TokenEnd
	TokenError
)

func (t Token) String() string {
	var names = []string{"TokenNewline", "TokenQuote", "TokenDelimiter", "TokenByte", "TokenEnd", "TokenError"}
	if int(t) < len(names) {
		return names[t]
	}
	return "TokenUnknown"
}

func Byte2Token(b, delimiter byte) Token {
	switch b {
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
