package mib

import (
	"fmt"

	"github.com/goller/mib/tokens"
)

// NewExpectedTokenError informs that a token was missing where expected.
func NewExpectedTokenError(t tokens.Token) error {
	return fmt.Errorf("missing %s", t)
}

// NewEOFError informs that an end of file was found when expecting
// a token.
func NewEOFError(expected tokens.Token) error {
	return fmt.Errorf("unexpected end of file; missing %s", expected)
}
