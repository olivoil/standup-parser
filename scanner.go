package parser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case ':':
		return COLON, string(ch)
	default:
		s.unread()
		return s.scanIdent()
	}
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if isLineBreak(ch) {
			s.unread()
			break
		} else if ch == ':' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.TrimSpace(strings.Trim(strings.ToUpper(buf.String()), "_*-+>")) {

	case "TODAY":
		return TODAY, buf.String()

	case "YESTERDAY":
		return YESTERDAY, buf.String()
	case "WEEKEND":
		return YESTERDAY, buf.String()
	case "WEEK-END":
		return YESTERDAY, buf.String()
	case "FRIDAY":
		return YESTERDAY, buf.String()
	case "FRIDAY/WEEKEND":
		return YESTERDAY, buf.String()

	case "MEETING":
		return MEETINGS, buf.String()
	case "MEETINGS":
		return MEETINGS, buf.String()

	case "BLOCKER":
		return BLOCKERS, buf.String()
	case "BLOCKERS":
		return BLOCKERS, buf.String()

	case "TIME":
		return LP, buf.String()
	case "HOURS":
		return LP, buf.String()
	case "LP":
		return LP, buf.String()

	case "JIRA":
		return JIRA, buf.String()
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch) || ch == ' ' || ch == '\t' || ch == '\u2002' || isLineBreak(ch)
}

// isLineBreak returns true if the rune is a space, tab, or newline.
func isLineBreak(ch rune) bool { return ch == '\n' }

// isAlphanumeric returns true if the rune is a letter or a number.
func isAlphanumeric(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
