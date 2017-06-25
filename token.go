package parser

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	EOF   Token = iota
	WS          // \t \n \s
	COLON       // :

	// Literals
	IDENT // main

	// Keywords
	TODAY
	YESTERDAY
	MEETINGS
	BLOCKERS
	LP
	JIRA
)

// isKeyword is true if the Token `t` is a keyword.
func isKeyword(t Token) bool {
	return t == TODAY ||
		t == YESTERDAY ||
		t == MEETINGS ||
		t == BLOCKERS ||
		t == LP ||
		t == JIRA
}
