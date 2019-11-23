package oids

import (
	"testing"

	"github.com/goller/mib/tokens"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name  string
		lexer *tokens.Lexer
	}{
		{
			name: "find something",
			lexer: func() *tokens.Lexer {
				return tokens.NewLexer("internet        OBJECT IDENTIFIER ::= { iso org(3) dod(6) 1 }")
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Search(tt.lexer)
			t.Errorf("res %v", res)
		})
	}
}
