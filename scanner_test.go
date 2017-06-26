package parser_test

import (
	"strings"
	"testing"

	"github.com/olivoil/standup-parser"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok parser.Token
		lit string
	}{
		// Special tokens (EOF, WS, COLON)
		{s: ``, tok: parser.EOF},
		{s: ` `, tok: parser.WS, lit: " "},
		{s: "\t", tok: parser.WS, lit: "\t"},
		{s: "\n", tok: parser.WS, lit: "\n"},
		{s: ":", tok: parser.COLON, lit: ":"},

		// Identifiers
		{s: `foo`, tok: parser.IDENT, lit: `foo`},
		{s: `Zx12_3U_-`, tok: parser.IDENT, lit: `Zx12_3U_-`},
		{s: `yourtrainer, energi`, tok: parser.IDENT, lit: `yourtrainer, energi`},
		{s: `project: something\nproject: something else`, tok: parser.IDENT, lit: `project`},

		// Keywords
		{s: `TODAY`, tok: parser.TODAY, lit: "TODAY"},
		{s: `Yesterday`, tok: parser.YESTERDAY, lit: "Yesterday"},
		{s: `Friday`, tok: parser.YESTERDAY, lit: "Friday"},
		{s: `Friday/weekend`, tok: parser.YESTERDAY, lit: "Friday/weekend"},
		{s: `meetings`, tok: parser.MEETINGS, lit: "meetings"},
		{s: `meetings:`, tok: parser.MEETINGS, lit: "meetings"},
		{s: `- meetings: hello`, tok: parser.MEETINGS, lit: "- meetings"},
		{s: `blockers`, tok: parser.BLOCKERS, lit: "blockers"},
		{s: `LP`, tok: parser.LP, lit: "LP"},
		{s: `Jira`, tok: parser.JIRA, lit: "Jira"},
	}

	for i, tt := range tests {
		s := parser.NewScanner(strings.NewReader(tt.s))
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
