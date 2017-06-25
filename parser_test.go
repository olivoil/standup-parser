package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"bitbucket.org/RocksauceStudios/standup-parser"
	"github.com/davecgh/go-spew/spew"
)

// Ensure the parser can parse strings into Standup ASTs.
func TestParser_ParseStandup(t *testing.T) {
	var tests = map[string]struct {
		s    string
		stmt *parser.Statement
		err  string
	}{
		"single field statement with yesterday": {
			s: `yesterday: ibm, slack`,
			stmt: &parser.Statement{
				Yesterday: parser.StringField{
					Key:   "yesterday",
					Val:   "ibm, slack",
					Valid: true,
				},
			},
		},

		"single field statement with today": {
			s: `
today:
- ibm: work on something
- slack: something else`,
			stmt: &parser.Statement{
				Today: parser.StringField{
					Key:   "today",
					Val:   "- ibm: work on something\n- slack: something else",
					Valid: true,
				},
			},
		},

		"single field statement without keyword": {
			s: `working on something`,
			stmt: &parser.Statement{
				Today: parser.StringField{
					Key:   "",
					Val:   `working on something`,
					Valid: true,
				},
			},
		},

		"typical statement from chrome app": {
			s: `
Friday: yourtrainer, halo, it's your birthday
Today:
  - halo: finish deployment?
  - yourtrainer: last issues
  - coomo: architecture planning
  - meetings: none
  - blockers: none
LP: up to date
Jira: not yet
`,
			stmt: &parser.Statement{
				Yesterday: parser.StringField{
					Key:   "Friday",
					Val:   `yourtrainer, halo, it's your birthday`,
					Valid: true,
				},
				Today: parser.StringField{
					Key:   "Today",
					Val:   "- halo: finish deployment?\n- yourtrainer: last issues\n- coomo: architecture planning",
					Valid: true,
				},
				Meetings: parser.StringField{
					Key:   "- meetings",
					Val:   "none",
					Valid: true,
				},
				Blockers: parser.StringField{
					Key:   "- blockers",
					Val:   "none",
					Valid: true,
				},
				LP: parser.BoolField{
					Key:   "LP",
					Val:   true,
					Lit:   "up to date",
					Valid: true,
				},
				Jira: parser.BoolField{
					Key:   "Jira",
					Val:   false,
					Lit:   "not yet",
					Valid: true,
				},
			},
		},

		"Alice's standup": {
			s: `
Friday: NewCo, Knod, Solitaire
Today:
  - Possibly NewCo, QA needs
  -client revisions
LP: up to date
`,
			stmt: &parser.Statement{
				Yesterday: parser.StringField{
					Key:   "Friday",
					Val:   `NewCo, Knod, Solitaire`,
					Valid: true,
				},
				Today: parser.StringField{
					Key:   "Today",
					Val:   "- Possibly NewCo, QA needs\n-client revisions",
					Valid: true,
				},
				LP: parser.BoolField{
					Key:   "LP",
					Val:   true,
					Lit:   "up to date",
					Valid: true,
				},
			},
		},

		"Chris's standup": {
			s: `
Friday: IBM, CooMo
Today: CooMo
time: current
`,
			stmt: &parser.Statement{
				Yesterday: parser.StringField{
					Key:   "Friday",
					Val:   `IBM, CooMo`,
					Valid: true,
				},
				Today: parser.StringField{
					Key:   "Today",
					Val:   "CooMo",
					Valid: true,
				},
				LP: parser.BoolField{
					Key:   "time",
					Val:   true,
					Lit:   "current",
					Valid: true,
				},
			},
		},

		"RJ's standup": {
			s: `
Friday: Mistbox, CFL
Today: NewCo Naming, Mistbox Slices/Redlines, ACN Enrollment Design?
Hours are up to date
`,
			stmt: &parser.Statement{
				Yesterday: parser.StringField{
					Key:   "Friday",
					Val:   `Mistbox, CFL`,
					Valid: true,
				},
				Today: parser.StringField{
					Key:   "Today",
					Val:   "NewCo Naming, Mistbox Slices/Redlines, ACN Enrollment Design?\nHours are up to date",
					Valid: true,
				},
			},
		},

		"John's standup": {
			s: `
Today: Meetings & Coomo
Friday: ACN
LP: updated
`,
			stmt: &parser.Statement{
				Yesterday: parser.StringField{
					Key:   "Friday",
					Val:   `ACN`,
					Valid: true,
				},
				Today: parser.StringField{
					Key:   "Today",
					Val:   "Meetings & Coomo",
					Valid: true,
				},
				LP: parser.BoolField{
					Key:   "LP",
					Val:   true,
					Lit:   "updated",
					Valid: true,
				},
			},
		},

		"Jason's standup": {
			s: `
Friday: meetings, IBM, Highball
Today:
  - Highball
  - Meetings all day
  - meetings: Huddle, UX w/ John, UX w/ Alice, IYB call, WIG, Leadership
  - blockers: none
LP: up to date
Jira: up to date
`,
			stmt: &parser.Statement{
				Yesterday: parser.StringField{
					Key:   "Friday",
					Val:   `meetings, IBM, Highball`,
					Valid: true,
				},
				Today: parser.StringField{
					Key:   "Today",
					Val:   "- Highball\n- Meetings all day",
					Valid: true,
				},
				Meetings: parser.StringField{
					Key:   "- meetings",
					Val:   "Huddle, UX w/ John, UX w/ Alice, IYB call, WIG, Leadership",
					Valid: true,
				},
				Blockers: parser.StringField{
					Key:   "- blockers",
					Val:   "none",
					Valid: true,
				},
				LP: parser.BoolField{
					Key:   "LP",
					Val:   true,
					Lit:   "up to date",
					Valid: true,
				},
				Jira: parser.BoolField{
					Key:   "Jira",
					Val:   true,
					Lit:   "up to date",
					Valid: true,
				},
			},
		},
	}

	for label, tt := range tests {
		stmt, err := parser.New(strings.NewReader(tt.s)).Parse()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf(
				"[%v] %q: error mismatch:\n  exp=%s\n  got=%s\n\n",
				label, tt.s, tt.err, err,
			)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf(
				"[%v] %q\n\nstmt mismatch:\n\nexp=%v\n\ngot=%v\n\n",
				label, tt.s, spew.Sdump(tt.stmt), spew.Sdump(stmt),
			)
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
