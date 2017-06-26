package parser

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

// Statement represents a standup statement.
type Statement struct {
	Yesterday StringField `json:"yesterday"`
	Today     StringField `json:"today"`
	Meetings  StringField `json:"meetings"`
	Blockers  StringField `json:"blockers"`
	LP        BoolField `json:"lp"`
	Jira      BoolField `json:"jira"`
}

// StringField is a key/value pair that holds one or several string values
type StringField struct {
	Key   string `json:"key"`
	Val   string `json:"val"`
	Valid bool `json:"valid"`
}

// BoolField is a key/value pair that holds one boolean value
type BoolField struct {
	Key   string `json:"key"`
	Val   bool `json:"val"`
	Lit   string `json:"lit"`
	Valid bool `json:"valid"`
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// New returns a new instance of Parser.
func New(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses a Statement.
func (p *Parser) Parse() (*Statement, error) {
	stmt := &Statement{}

	// loop over all tokens
	for {
		// Read a keyword and its statement
		key, keyLit, _ := p.scanIgnoreWhitespace()
		if key == EOF {
			break
		}

		// if it does not start with a keyword, consider it's TODAY
		if !isKeyword(key) {
			p.unscan()
			key = TODAY
			keyLit = ""
		}

		// keyword is optionally followed by a colon. Ignore it.
		col, _, _ := p.scanIgnoreWhitespace()
		if col != COLON {
			p.unscan()
		}

		values := []string{}

		for {
			tok, lit, ws := p.scanIgnoreWhitespace()
			if isKeyword(tok) || tok == EOF {
				p.unscan()
				break
			}

			if tok == IDENT || tok == COLON {
				values = append(values, ws, lit)
			}
		}

		switch key {
		case TODAY:
			val := splitAndTrimSpace(values)
			stmt.Today = StringField{
				Key:   keyLit,
				Val:   val,
				Valid: val != "",
			}
		case YESTERDAY:
			val := splitAndTrimSpace(values)
			stmt.Yesterday = StringField{
				Key:   keyLit,
				Val:   val,
				Valid: val != "",
			}
		case MEETINGS:
			val := splitAndTrimSpace(values)
			stmt.Meetings = StringField{
				Key:   keyLit,
				Val:   val,
				Valid: val != "",
			}
		case BLOCKERS:
			val := splitAndTrimSpace(values)
			stmt.Blockers = StringField{
				Key:   keyLit,
				Val:   val,
				Valid: val != "",
			}
		case LP:
			lit := splitAndTrimSpace(values)
			val, err := isPositive(lit)

			stmt.LP = BoolField{
				Key:   keyLit,
				Val:   val,
				Lit:   lit,
				Valid: err == nil,
			}
		case JIRA:
			lit := splitAndTrimSpace(values)
			val, err := isPositive(lit)

			stmt.Jira = BoolField{
				Key:   keyLit,
				Val:   val,
				Lit:   lit,
				Valid: err == nil,
			}
		}
	}

	return stmt, nil
}

// isPositive is a naive attempt at determining
// if the string representation of a boolean value is true or false.
func isPositive(s string) (bool, error) {
	negative := regexp.MustCompile(`.*(no|off|updating|negative).*`)
	positive := regexp.MustCompile(`.*(done|yes|up\s+to\s+date|ok|1|affirmative|current|updated)`)

	n := negative.Match([]byte(s))
	p := positive.Match([]byte(s))

	if p && n {
		return true, errors.New("ambiguous")
	}
	if !p && !n {
		return true, errors.New("unclear")
	}

	return p && !n, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string, ws string) {
	tok, lit = p.scan()
	if tok == WS {
		ws = lit
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

func splitAndTrimSpace(values []string) string {
	val := strings.TrimSpace(strings.Join(values, ""))
	lines := strings.Split(val, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimSpace(l)
	}
	return strings.Join(lines, "\n")
}
