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
		want  []TokenType
	}{
		{
			name:  "full mib",
			input: rfc1215,
			want: []TokenType{
				Label,
				Definitions,
				Equals,
				Begin,
				TrapType,
				Macro,
				Equals,
				Begin,
				Label, // Type
				Label,
				Equals,
				Quotestring,
				Label,
				LeftParen,
				Enterprise,
				Object,
				Identifier,
				RightParen,
				Label,
				Label,
				Label,
				Label, // Value
				Label,
				Equals,
				Label,
				LeftParen,
				Label,
				Integer,
				RightParen,
				Label, // VarPart
				Equals,
				Quotestring,
				Quotestring,
				Label,
				Quotestring,
				Bar,
				Label,
				Label, // VarTypes
				Equals,
				Label,
				Bar,
				Label,
				Quotestring,
				Label,
				Label, // VarType
				Equals,
				Label,
				LeftParen,
				Label,
				Objname,
				RightParen,
				Label, // DescrPart
				Equals,
				Quotestring,
				Label,
				LeftParen,
				Description,
				Label,
				RightParen,
				Bar,
				Label,
				Label, // ReferPart
				Equals,
				Quotestring,
				Label,
				LeftParen,
				Reference,
				Label,
				RightParen,
				Bar,
				Label,
				End,
				End,
				EOF,
			},
		},
		{
			name: "comment newline eof",
			input: `
					-- comment
					`,
			want: []TokenType{EOF},
		},
		{
			name:  "single letter",
			input: `a`,
			want:  []TokenType{Label, EOF},
		},
		{
			name:  "empty mib",
			input: ``,
			want:  []TokenType{EOF},
		},

		{
			name: "comment eof",
			input: `
					-- comment`,
			want: []TokenType{EOF},
		},
		{
			name: "inline comment eof",
			input: `
					-- comment -- howdy`,
			want: []TokenType{Label, EOF},
		},
		{
			name:  "range",
			input: `..`,
			want:  []TokenType{Range, EOF},
		},
		{
			name:  "period prefixed label",
			input: `.label`,
			want:  []TokenType{Label, Label, EOF},
		},
		{
			name:  "equals",
			input: `::=`,
			want:  []TokenType{Equals, EOF},
		},
		{
			name:  "colon prefixed label",
			input: `:label`,
			want:  []TokenType{Label, Label, EOF},
		},
		{
			name:  "double colon prefixed label",
			input: `::label`,
			want:  []TokenType{Label, Label, EOF},
		},
		{
			name:  "hex number empty h",
			input: `''h`,
			want:  []TokenType{Hex, EOF},
		},
		{
			name:  "empty number literal (c parser is only EOF)",
			input: `''`,
			want:  []TokenType{EOF},
		},
		{
			name:  "hex number h",
			input: `'fedcba9876543210'h`,
			want:  []TokenType{Hex, EOF},
		},
		{
			name:  "hex number H",
			input: `'0123456789abcdef'H`,
			want:  []TokenType{Hex, EOF},
		},
		{
			name:  "hex number empty H",
			input: `''H`,
			want:  []TokenType{Hex, EOF},
		},
		{
			name:  "binary number empty b",
			input: `''b`,
			want:  []TokenType{Binary, EOF},
		},
		{
			name:  "binary number b",
			input: `'1010'b`,
			want:  []TokenType{Binary, EOF},
		},
		{
			name:  "binary number B",
			input: `'1010'B`,
			want:  []TokenType{Binary, EOF},
		},
		{
			name:  "binary number empty B",
			input: `''B`,
			want:  []TokenType{Binary, EOF},
		},
		{
			name:  "binary number eof",
			input: `'`,
			want:  []TokenType{Error},
		},
		{
			name:  "number label no digits",
			input: `'label'`,
			want:  []TokenType{EOF},
		},
		{
			name:  "number label with digits",
			input: `'01'`,
			want:  []TokenType{EOF},
		},
		{
			name:  "unknown number type",
			input: `'01'u`,
			want:  []TokenType{Label, EOF},
		},
		{
			name:  "string",
			input: `"string"`,
			want:  []TokenType{Quotestring, EOF},
		},
		{
			name:  "string with newlines",
			input: "\"string\r\nstring\"",
			want:  []TokenType{Quotestring, EOF},
		},
		{
			name:  "non-terminating string",
			input: `"string`,
			want:  []TokenType{Error},
		},
		{
			name:  "left paren",
			input: "(",
			want:  []TokenType{LeftParen, EOF},
		},
		{
			name:  "right paren",
			input: ")",
			want:  []TokenType{RightParen, EOF},
		},
		{
			name:  "left bracket",
			input: "{",
			want:  []TokenType{LeftBracket, EOF},
		},
		{
			name:  "right bracket",
			input: "}",
			want:  []TokenType{RightBracket, EOF},
		},
		{
			name:  "left square bracket",
			input: "[",
			want:  []TokenType{LeftSquareBracket, EOF},
		},
		{
			name:  "right square bracket",
			input: "]",
			want:  []TokenType{RightSquareBracket, EOF},
		},
		{
			name:  "semicolon",
			input: ";",
			want:  []TokenType{Semicolon, EOF},
		},
		{
			name:  "comma",
			input: ",",
			want:  []TokenType{Comma, EOF},
		},
		{
			name:  "non-printable",
			input: "😀",
			want:  []TokenType{Error},
		},
		{
			name:  "imports keyword",
			input: "imports",
			want:  []TokenType{Imports, EOF},
		},
		{
			name:  "non label chars for some reason are labels",
			input: "$",
			want:  []TokenType{Label, EOF},
		},
		{
			name:  "number",
			input: "2",
			want:  []TokenType{Number, EOF},
		},
		{
			name:  "object identifier",
			input: "{iso 2}",
			want:  []TokenType{LeftBracket, Label, Number, RightBracket, EOF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := []Token{}
			for {
				token := lexer.NextToken()
				tokens = append(tokens, token)
				if token.Typ == EOF || token.Typ == Error {
					break
				}
			}

			if len(tokens) != len(tt.want) {
				t.Logf("Got %v", tokens)
				t.Logf("Want %v", tt.want)
				t.Fatalf("unexpected difference number of tokens. got %d want %d", len(tokens), len(tt.want))
			}

			for i := range tt.want {
				if got, want := tokens[i].Typ, tt.want[i]; got != want {
					t.Errorf("unexpected token type: %d %s want %d", tokens[i].Typ, tokens[i], want)
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
			lexer := NewLexer(string(b))
			for {
				token := lexer.NextToken()
				if token.Typ == Error {
					t.Errorf("unexpected error %s", token)
					break
				}
				if token.Typ == EOF {
					break
				}
			}
		})
	}
}

func Test_BigMIB(t *testing.T) {
	b, err := ioutil.ReadFile("TIMETRA-SUBSCRIBER-MGMT-MIB")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("TIMETRA-SUBSCRIBER-MGMT-MIB", func(t *testing.T) {
		lexer := NewLexer(string(b))
		for {
			token := lexer.NextToken()
			if token.Typ == Error {
				t.Errorf("unexpected error %s", token)
				break
			}
			if token.Typ == EOF {
				break
			}
		}
	})
}

func Test_token_String(t *testing.T) {
	type fields struct {
		typ TokenType
		val string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "eof",
			fields: fields{
				typ: EOF,
			},
			want: "EOF",
		},
		{
			name: "error",
			fields: fields{
				typ: Error,
				val: "error",
			},
			want: "error",
		},
		{
			name: "keyword",
			fields: fields{
				typ: Organization,
				val: "organization",
			},
			want: "<organization>",
		},
		{
			name: "left paren",
			fields: fields{
				typ: LeftBracket,
				val: "(",
			},
			want: `"("`,
		},
		{
			name: "label",
			fields: fields{
				typ: Label,
				val: "123456789abcd",
			},
			want: `"123456789a"...`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tkn := Token{
				Typ: tt.fields.typ,
				Val: tt.fields.val,
			}
			if got := tkn.String(); got != tt.want {
				t.Errorf("token.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_lex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(rfc1215)
		for {
			token := lexer.NextToken()
			if token.Typ == EOF || token.Typ == Error {
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
			lexer := NewLexer(mibs[j])
			for {
				token := lexer.NextToken()
				if token.Typ == Error {
					break
				}
				if token.Typ == EOF {
					break
				}
			}
		}
	}
}

func Benchmark_Tons(b *testing.B) {
	files, _ := ioutil.ReadDir("mib")

	mibs := []string{}
	for _, file := range files {
		name := filepath.Join("mib", file.Name())
		buf, _ := ioutil.ReadFile(name)
		mibs = append(mibs, string(buf))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range mibs {
			lexer := NewLexer(mibs[j])
			for {
				token := lexer.NextToken()
				if token.Typ == Error {
					break
				}
				if token.Typ == EOF {
					break
				}
			}
		}
	}
}

func Benchmark_BigMIB(b *testing.B) {
	buf, _ := ioutil.ReadFile("TIMETRA-SUBSCRIBER-MGMT-MIB")
	mib := string(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(mib)
		for {
			token := lexer.NextToken()
			if token.Typ == Error {
				break
			}
			if token.Typ == EOF {
				break
			}
		}
	}
}
