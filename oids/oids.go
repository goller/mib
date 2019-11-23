package oids

import (
	"github.com/goller/mib/tokens"
)

/*
 * Parse a sequence of object subidentifiers for the given name.
 * The "label OBJECT IDENTIFIER ::=" portion has already been parsed.
 *
 * The majority of cases take this form :
 * label OBJECT IDENTIFIER ::= { parent 2 }
 * where a parent label and a child subidentifier number are specified.
 *
 * Variations on the theme include cases where a number appears with
 * the parent, or intermediate subidentifiers are specified by label,
 * by number, or both.
 *
 * Here are some representative samples :
 * internet        OBJECT IDENTIFIER ::= { iso org(3) dod(6) 1 }
 * mgmt            OBJECT IDENTIFIER ::= { internet 2 }
 * rptrInfoHealth  OBJECT IDENTIFIER ::= { snmpDot3RptrMgt 0 4 }
 *
 * Here is a very rare form :
 * iso             OBJECT IDENTIFIER ::= { 1 }
 *
 * Returns NULL on error.  When this happens, memory may be leaked.
 */

/// OIDS returns oids
func Search(lexer *tokens.Lexer) []tokens.Token {
	searchFor := []tokens.TokenType{
		tokens.TokenLabel,
		tokens.TokenObject,
		tokens.TokenIdentifier,
		tokens.TokenEquals,
		tokens.TokenLeftBracket,
	}

	res := make([]tokens.Token, len(searchFor))

	want := 0 // start with the zeroth index
	for tk := lexer.NextToken(); tk.Typ != tokens.TokenEOF && tk.Typ != tokens.TokenError; tk = lexer.NextToken() {
		if tk.Typ != searchFor[want] {
			continue
		}
		res[want] = tk
		want++
		if want == len(searchFor) {
			break
		}
	}

	if want != len(searchFor) {
		return nil
	}

	until := tokens.TokenRightBracket

	return res
}
