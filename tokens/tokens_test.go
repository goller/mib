package tokens

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

const rfc1215 = `
RFC-1215 DEFINITIONS ::= BEGIN

-- This  module is a empty module.

TRAP-TYPE MACRO ::=
BEGIN
	TYPE NOTATION ::= "ENTERPRISE" value
					  (enterprise OBJECT IDENTIFIER)
					  VarPart
					  DescrPart
					  ReferPart
	VALUE NOTATION ::= value (VALUE INTEGER)
	VarPart ::=
			   "VARIABLES" "{" VarTypes "}"
			   | empty
	VarTypes ::=
			   VarType | VarTypes "," VarType
	VarType ::=
			   value (vartype ObjectName)
	DescrPart ::=
			   "DESCRIPTION" value (description DisplayString)
			   | empty
	ReferPart ::=
			   "REFERENCE" value (reference DisplayString)
			   | empty
END

END
`

func Test_lex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []tokenType
	}{
		{
			name:  "full mib",
			input: rfc1215,
			want: []tokenType{
				tokenLabel,
				tokenDefinitions,
				tokenEquals,
				tokenBegin,
				tokenTrapType,
				tokenMacro,
				tokenEquals,
				tokenBegin,
				tokenLabel, // Type
				tokenLabel,
				tokenEquals,
				tokenQuotestring,
				tokenLabel,
				tokenLeftParen,
				tokenEnterprise,
				tokenObject,
				tokenIdentifier,
				tokenRightParen,
				tokenLabel,
				tokenLabel,
				tokenLabel,
				tokenLabel, // Value
				tokenLabel,
				tokenEquals,
				tokenLabel,
				tokenLeftParen,
				tokenLabel,
				tokenInteger,
				tokenRightParen,
				tokenLabel, // VarPart
				tokenEquals,
				tokenQuotestring,
				tokenQuotestring,
				tokenLabel,
				tokenQuotestring,
				tokenBar,
				tokenLabel,
				tokenLabel, // VarTypes
				tokenEquals,
				tokenLabel,
				tokenBar,
				tokenLabel,
				tokenQuotestring,
				tokenLabel,
				tokenLabel, // VarType
				tokenEquals,
				tokenLabel,
				tokenLeftParen,
				tokenLabel,
				tokenObjname,
				tokenRightParen,
				tokenLabel, // DescrPart
				tokenEquals,
				tokenQuotestring,
				tokenLabel,
				tokenLeftParen,
				tokenDescription,
				tokenLabel,
				tokenRightParen,
				tokenBar,
				tokenLabel,
				tokenLabel, // ReferPart
				tokenEquals,
				tokenQuotestring,
				tokenLabel,
				tokenLeftParen,
				tokenReference,
				tokenLabel,
				tokenRightParen,
				tokenBar,
				tokenLabel,
				tokenEnd,
				tokenEnd,
				tokenEOF,
			},
		},
		{
			name: "comment newline eof",
			input: `
				-- comment
				`,
			want: []tokenType{tokenEOF},
		},
		{
			name: "comment eof",
			input: `
				-- comment`,
			want: []tokenType{tokenEOF},
		},
		{
			name: "inline comment eof",
			input: `
				-- comment -- howdy`,
			want: []tokenType{tokenLabel, tokenEOF},
		},
		{
			name:  "range",
			input: `..`,
			want:  []tokenType{tokenRange, tokenEOF},
		},
		{
			name:  "period prefixed label",
			input: `.label`,
			want:  []tokenType{tokenLabel, tokenLabel, tokenEOF},
		},
		{
			name:  "equals",
			input: `::=`,
			want:  []tokenType{tokenEquals, tokenEOF},
		},
		{
			name:  "colon prefixed label",
			input: `:label`,
			want:  []tokenType{tokenLabel, tokenLabel, tokenEOF},
		},
		{
			name:  "double colon prefixed label",
			input: `::label`,
			want:  []tokenType{tokenLabel, tokenLabel, tokenEOF},
		},
		{
			name:  "hex number empty h",
			input: `''h`,
			want:  []tokenType{tokenHex, tokenEOF},
		},
		{
			name:  "empty number literal (c parser is only tokenEOF)",
			input: `''`,
			want:  []tokenType{tokenEOF},
		},
		{
			name:  "hex number h",
			input: `'fedcba9876543210'h`,
			want:  []tokenType{tokenHex, tokenEOF},
		},
		{
			name:  "hex number H",
			input: `'0123456789abcdef'H`,
			want:  []tokenType{tokenHex, tokenEOF},
		},
		{
			name:  "hex number empty H",
			input: `''H`,
			want:  []tokenType{tokenHex, tokenEOF},
		},
		{
			name:  "binary number empty b",
			input: `''b`,
			want:  []tokenType{tokenBinary, tokenEOF},
		},
		{
			name:  "binary number b",
			input: `'1010'b`,
			want:  []tokenType{tokenBinary, tokenEOF},
		},
		{
			name:  "binary number B",
			input: `'1010'B`,
			want:  []tokenType{tokenBinary, tokenEOF},
		},
		{
			name:  "binary number empty B",
			input: `''B`,
			want:  []tokenType{tokenBinary, tokenEOF},
		},
		{
			name:  "binary number eof",
			input: `'`,
			want:  []tokenType{tokenError},
		},
		{
			name:  "number label no digits",
			input: `'label'`,
			want:  []tokenType{tokenEOF},
		},
		{
			name:  "number label with digits",
			input: `'01'`,
			want:  []tokenType{tokenEOF},
		},
		{
			name:  "unknown number type",
			input: `'01'u`,
			want:  []tokenType{tokenLabel, tokenEOF},
		},
		{
			name:  "string",
			input: `"string"`,
			want:  []tokenType{tokenQuotestring, tokenEOF},
		},
		{
			name:  "string with newlines",
			input: "\"string\r\nstring\"",
			want:  []tokenType{tokenQuotestring, tokenEOF},
		},
		{
			name:  "non-terminating string",
			input: `"string`,
			want:  []tokenType{tokenError},
		},
		{
			name:  "left paren",
			input: "(",
			want:  []tokenType{tokenLeftParen, tokenEOF},
		},
		{
			name:  "right paren",
			input: ")",
			want:  []tokenType{tokenRightParen, tokenEOF},
		},
		{
			name:  "left bracket",
			input: "{",
			want:  []tokenType{tokenLeftBracket, tokenEOF},
		},
		{
			name:  "right bracket",
			input: "}",
			want:  []tokenType{tokenRightBracket, tokenEOF},
		},
		{
			name:  "left square bracket",
			input: "[",
			want:  []tokenType{tokenLeftSquareBracket, tokenEOF},
		},
		{
			name:  "right square bracket",
			input: "]",
			want:  []tokenType{tokenRightSquareBracket, tokenEOF},
		},
		{
			name:  "semicolon",
			input: ";",
			want:  []tokenType{tokenSemicolon, tokenEOF},
		},
		{
			name:  "comma",
			input: ",",
			want:  []tokenType{tokenComma, tokenEOF},
		},
		{
			name:  "non-printable",
			input: "😀",
			want:  []tokenType{tokenError},
		},
		{
			name:  "imports keyword",
			input: "imports",
			want:  []tokenType{tokenImports, tokenEOF},
		},
		{
			name:  "non label chars for some reason are labels",
			input: "$",
			want:  []tokenType{tokenLabel, tokenEOF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := lex(tt.name, tt.input)
			tokens := []token{}
			for {
				token := lexer.nextToken()
				tokens = append(tokens, token)
				if token.typ == tokenEOF || token.typ == tokenError {
					break
				}
			}

			if len(tokens) != len(tt.want) {
				t.Logf("Got %v", tokens)
				t.Logf("Want %v", tt.want)
				t.Fatalf("unexpected difference number of tokens. got %d want %d", len(tokens), len(tt.want))
			}

			for i := range tt.want {
				if got, want := tokens[i].typ, tt.want[i]; got != want {
					t.Errorf("unexpected token type: %d %s want %d", tokens[i].typ, tokens[i], want)
				}
			}
		})
	}
}
func Test_SNMP(t *testing.T) {
	files, err := ioutil.ReadDir("/usr/share/snmp/mibs")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		name := filepath.Join("/usr/share/snmp/mibs", file.Name())
		b, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		t.Run(name, func(t *testing.T) {
			lexer := lex(name, string(b))
			for {
				token := lexer.nextToken()
				if token.typ == tokenError {
					t.Errorf("unexpected error %s", token)
					break
				}
				if token.typ == tokenEOF {
					break
				}
			}
		})
	}
}

func Test_token_String(t *testing.T) {
	type fields struct {
		typ  tokenType
		pos  Pos
		val  string
		line int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "eof",
			fields: fields{
				typ: tokenEOF,
			},
			want: "EOF",
		},
		{
			name: "error",
			fields: fields{
				typ: tokenError,
				val: "error",
			},
			want: "error",
		},
		{
			name: "keyword",
			fields: fields{
				typ: tokenOrganization,
				val: "organization",
			},
			want: "<organization>",
		},
		{
			name: "left paren",
			fields: fields{
				typ: tokenLeftBracket,
				val: "(",
			},
			want: `"("`,
		},
		{
			name: "label",
			fields: fields{
				typ: tokenLabel,
				val: "123456789abcd",
			},
			want: `"123456789a"...`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tkn := token{
				typ:  tt.fields.typ,
				pos:  tt.fields.pos,
				val:  tt.fields.val,
				line: tt.fields.line,
			}
			if got := tkn.String(); got != tt.want {
				t.Errorf("token.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_lex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lexer := lex("rfc1215", rfc1215)
		for {
			token := lexer.nextToken()
			if token.typ == tokenEOF || token.typ == tokenError {
				break
			}
		}
	}
}

func Benchmark_Dir(b *testing.B) {
	files, _ := ioutil.ReadDir("/usr/share/snmp/mibs")

	mibs := []string{}
	for _, file := range files {
		name := filepath.Join("/usr/share/snmp/mibs", file.Name())
		buf, _ := ioutil.ReadFile(name)
		mibs = append(mibs, string(buf))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range mibs {
			lexer := lex("a", mibs[j])
			for {
				token := lexer.nextToken()
				if token.typ == tokenError {
					break
				}
				if token.typ == tokenEOF {
					break
				}
			}
		}
	}
}
