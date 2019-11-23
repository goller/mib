package tokens

import (
	"fmt"
	"strings"
)

// TokenType is one of the specific MIB token types.
type TokenType uint

// Token is a string in a MIB file with an identified meaning.
type Token struct {
	Typ TokenType
	Val string
}

func (t Token) String() string {
	switch {
	case t.Typ == EOF:
		return "EOF"
	case t.Typ == Error:
		return t.Val
	case t.Typ > Keyword:
		return fmt.Sprintf("<%s>", t.Val)
	case len(t.Val) > 10:
		return fmt.Sprintf("%.10q...", t.Val)
	}

	return fmt.Sprintf("%q", t.Val)
}

const (
	None TokenType = iota
	Error
	LeftParen
	RightParen
	LeftBracket
	RightBracket
	LeftSquareBracket
	RightSquareBracket
	Semicolon
	Comma
	Bar
	Range
	Label
	Equals
	Number
	EOF
	Keyword
	Obsolete
	KwOpaque
	KwOptional
	LastUpdated
	Organization
	ContactInfo
	ModuleIdentify
	Compliance
	Definitions
	End
	Augments
	NoAccess
	WriteOnly
	Nsapaddress
	Units
	Reference
	NumEntries
	Bitstring
	Continue
	BitString
	Counter64
	Timeticks
	NotifType
	ObjGroup
	ObjIdentity
	Identifier
	Object
	Netaddr
	Gauge
	Unsigned32
	ReadWrite
	ReadCreate
	Octetstr
	Of
	Sequence
	Nul
	Ipaddr
	Binary
	Hex
	Uinteger32
	Integer
	Integer32
	Counter
	ReadOnly
	Description
	Index
	Defval
	Deprecated
	Size
	Access
	Mandatory
	Current
	Status
	Syntax
	ObjType
	TrapType
	Enterprise
	Begin
	Imports
	Exports
	Accnotify
	Convention
	Notifgroup
	DisplayHint
	From
	AgentCap
	Macro
	Implied
	Supports
	Includes
	Variation
	Revision
	NotImpl
	Objects
	Notifications
	Module
	MinAccess
	ProdRel
	WrSyntax
	CreateReq
	MandatoryGroups
	Group
	Choice
	Implicit
	Objsyntax
	Simplesyntax
	Appsyntax
	Objname
	Notifname
	Variables
	Quotestring
)

// IsSyntax returns true if the token type is one of the valid types in a
// SYNTAX declaration.
func (t TokenType) IsSyntax() bool {
	switch t {
	case Identifier, Octetstr, Integer, Netaddr,
		Ipaddr, Counter, Gauge, Timeticks, KwOpaque,
		Nul, BitString, Nsapaddress, Counter64,
		Uinteger32, Appsyntax, Objsyntax, Simplesyntax,
		Objname, Notifname, Unsigned32, Integer32:
		return true
	default:
		return false
	}
}

var lexemes = map[string]TokenType{
	"OBSOLETE":              Obsolete,
	"OPAQUE":                KwOpaque,
	"OPTIONAL":              KwOptional,
	"LAST-UPDATED":          LastUpdated,
	"ORGANIZATION":          Organization,
	"CONTACT-INFO":          ContactInfo,
	"MODULE-IDENTITY":       ModuleIdentify,
	"MODULE-COMPLIANCE":     Compliance,
	"DEFINITIONS":           Definitions,
	"END":                   End,
	"AUGMENTS":              Augments,
	"NOT-ACCESSIBLE":        NoAccess,
	"WRITE-ONLY":            WriteOnly,
	"NSAPADDRESS":           Nsapaddress,
	"UNITS":                 Units,
	"REFERENCE":             Reference,
	"NUM-ENTRIES":           NumEntries,
	"BITSTRING":             Bitstring,
	"BIT":                   Continue,
	"BITS":                  BitString,
	"COUNTER64":             Counter64,
	"TIMETICKS":             Timeticks,
	"NOTIFICATION-TYPE":     NotifType,
	"OBJECT-GROUP":          ObjGroup,
	"OBJECT-IDENTITY":       ObjIdentity,
	"IDENTIFIER":            Identifier,
	"OBJECT":                Object,
	"NETWORKADDRESS":        Netaddr,
	"GAUGE":                 Gauge,
	"GAUGE32":               Gauge,
	"UNSIGNED32":            Unsigned32,
	"READ-WRITE":            ReadWrite,
	"READ-CREATE":           ReadCreate,
	"OCTETSTRING":           Octetstr,
	"OCTET":                 Continue,
	"OF":                    Of,
	"SEQUENCE":              Sequence,
	"NULL":                  Nul,
	"IPADDRESS":             Ipaddr,
	"UINTEGER32":            Uinteger32,
	"INTEGER":               Integer,
	"INTEGER32":             Integer32,
	"COUNTER":               Counter,
	"COUNTER32":             Counter,
	"READ-ONLY":             ReadOnly,
	"DESCRIPTION":           Description,
	"INDEX":                 Index,
	"DEFVAL":                Defval,
	"DEPRECATED":            Deprecated,
	"SIZE":                  Size,
	"MAX-ACCESS":            Access,
	"ACCESS":                Access,
	"MANDATORY":             Mandatory,
	"CURRENT":               Current,
	"STATUS":                Status,
	"SYNTAX":                Syntax,
	"OBJECT-TYPE":           ObjType,
	"TRAP-TYPE":             TrapType,
	"ENTERPRISE":            Enterprise,
	"BEGIN":                 Begin,
	"IMPORTS":               Imports,
	"EXPORTS":               Exports,
	"ACCESSIBLE-FOR-NOTIFY": Accnotify,
	"TEXTUAL-CONVENTION":    Convention,
	"NOTIFICATION-GROUP":    Notifgroup,
	"DISPLAY-HINT":          DisplayHint,
	"FROM":                  From,
	"AGENT-CAPABILITIES":    AgentCap,
	"MACRO":                 Macro,
	"IMPLIED":               Implied,
	"SUPPORTS":              Supports,
	"INCLUDES":              Includes,
	"VARIATION":             Variation,
	"REVISION":              Revision,
	"NOT-IMPLEMENTED":       NotImpl,
	"OBJECTS":               Objects,
	"NOTIFICATIONS":         Notifications,
	"MODULE":                Module,
	"MIN-ACCESS":            MinAccess,
	"PRODUCT-RELEASE":       ProdRel,
	"WRITE-SYNTAX":          WrSyntax,
	"CREATION-REQUIRES":     CreateReq,
	"MANDATORY-GROUPS":      MandatoryGroups,
	"GROUP":                 Group,
	"CHOICE":                Choice,
	"IMPLICIT":              Implicit,
	"OBJECTSYNTAX":          Objsyntax,
	"SIMPLESYNTAX":          Simplesyntax,
	"APPLICATIONSYNTAX":     Appsyntax,
	"OBJECTNAME":            Objname,
	"NOTIFICATIONNAME":      Notifname,
	"VARIABLES":             Variables,
	"QUOTEDSTRING":          Quotestring,
}

const (
	eof        = byte(0xFF)
	maxASCII   = byte(0x7F)
	spaceASCII = byte(0x20)
	aASCII     = byte(0x61)
	zASCII     = byte(0x7A)
	htASCII    = byte(0x09)
	lfASCII    = byte(0x0A)
	vtASCII    = byte(0x0B)
	crASCII    = byte(0x0C)
	ffASCII    = byte(0x0D)
)

// Pos represents a byte position in the original input text from which
// this token was parsed.
type Pos int

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) (stateFn, Token)

const (
	// LongestKeyword is 21 characters long; used to preallocate
	// max keyword label buffer.
	LongestKeyword = 21
)

// Lexer holds the state of the scanner.
type Lexer struct {
	input string               // string to scan
	state stateFn              // the next lexing function to enter
	pos   Pos                  // current position in the input
	start Pos                  // start position of this item
	width Pos                  // width of last []byte read from input
	label [LongestKeyword]byte // buffer used to compare keywords
}

// next returns the next byte in the input.
func (l *Lexer) next() byte {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r := l.input[l.pos]
	l.width = 1
	l.pos += l.width
	return r
}

// peek returns but does not consume the next []byte in the input.
func (l *Lexer) peek() byte {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one []byte. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *Lexer) emit(t TokenType) Token {
	tk := Token{
		Typ: t,
		Val: l.input[l.start:l.pos],
	}
	l.start = l.pos
	return tk
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *Lexer) errorf(format string, args ...interface{}) Token {
	return Token{Error, fmt.Sprintf(format, args...)}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	var tk Token
	for l.state != nil {
		l.state, tk = l.state(l)
		if tk.Typ != None {
			return tk
		}
	}
	return Token{}
}

// PushToken puts a token back into the lexer; used to simplify
// parsing optional tokens.
func (l *Lexer) PushToken(t Token) {
	p := &push{
		prev:  l.state,
		token: t,
	}
	l.state = p.pop
}

type push struct {
	prev  stateFn
	token Token
}

func (p *push) pop(l *Lexer) (stateFn, Token) {
	return p.prev, p.token
}

// NewLexer creates a new scanner for the input string.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		state: lexSpace,
	}
}

const (
	comment   = "--"
	dashASCII = byte(0x2D)
)

func lexText(l *Lexer) (stateFn, Token) {
	if strings.HasPrefix(l.input[l.pos:], "--") {
		l.ignore()
		return lexComment, Token{}
	}

	switch r := l.next(); {
	case r == eof:
		break
	case r == '"':
		return lexQuotedString, Token{}
	case r == '\'':
		return lexNumberLiteral, Token{}
	case r == '(':
		return lexSpace, l.emit(LeftParen)
	case r == ')':
		return lexSpace, l.emit(RightParen)
	case r == '{':
		return lexSpace, l.emit(LeftBracket)
	case r == '}':
		return lexSpace, l.emit(RightBracket)
	case r == '[':
		return lexSpace, l.emit(LeftSquareBracket)
	case r == ']':
		return lexSpace, l.emit(RightSquareBracket)
	case r == ';':
		return lexSpace, l.emit(Semicolon)
	case r == ',':
		return lexSpace, l.emit(Comma)
	case r == '|':
		return lexSpace, l.emit(Bar)
	case r == '.':
		return lexRange, Token{}
	case r == ':':
		return lexEquals, Token{}
	case r <= maxASCII && r >= spaceASCII:
		return lexChars, Token{}
	default:
		return nil, l.errorf("unrecognized character: %#U", r)
	}

	return nil, l.emit(EOF)
}

func lexQuotedString(l *Lexer) (stateFn, Token) {
Loop:
	for {
		switch l.next() {
		case '\r', '\n':
			continue
		case eof:
			return nil, l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	return lexSpace, l.emit(Quotestring)
}

func lexNumberLiteral(l *Lexer) (stateFn, Token) {
	const (
		binary uint = 1 << iota
		hex
		unknown
	)
	numType := binary

Loop:
	for {
		switch r := l.next(); {
		case r == '0' || r == '1':
			numType |= binary
		case ('0' <= r && r <= '9') || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F'):
			numType |= hex
		case r == eof:
			return nil, l.errorf("unterminated literal string")
		case r == '\'':
			break Loop
		default:
			numType |= unknown
		}
	}

	r := l.next()
	if r >= aASCII && r <= zASCII { // TODO(goller): comment as toupper
		r -= 32
	}
	switch {
	// TODO: check why this is not being used
	case r == 'B' && numType&binary == binary: // TODO(goller): shoudn't this only be binary?
		return lexSpace, l.emit(Binary)
	case r == 'H' && numType&unknown == 0:
		return lexSpace, l.emit(Hex)
	case r == eof:
		return lexSpace, l.emit(EOF)
	default:
		return lexSpace, l.emit(Label)
	}
}

func lexSpace(l *Lexer) (stateFn, Token) {
LOOP:
	for {
		switch r := l.peek(); r {
		case htASCII, lfASCII, vtASCII, ffASCII, crASCII, spaceASCII, 0x85, 0xA0:
			_ = l.next()
		default:
			break LOOP

		}
	}
	l.ignore()
	return lexText, Token{}
}

// lexEquals searches for ::= otherwise it assumes a label; assumes first `:`
// already consumed.
func lexEquals(l *Lexer) (stateFn, Token) {
	if l.next() != ':' {
		l.backup()
		return lexSpace, l.emit(Label)
	}
	if l.next() != '=' {
		l.backup()
		return lexSpace, l.emit(Label)
	}
	return lexSpace, l.emit(Equals)
}

// lexEquals searches for .. otherwise it assumes a label; assumes first `.`
// already consumed.
func lexRange(l *Lexer) (stateFn, Token) {
	if l.next() == '.' {
		return lexSpace, l.emit(Range)
	}
	l.backup()
	return lexSpace, l.emit(Label)
}

// lexComment treats the rest of the line or until another '--' as a comment;
// the left comment marker is known to be present.
func lexComment(l *Lexer) (stateFn, Token) {
	l.pos += Pos(len(comment))
	var prev byte
	for {
		switch r := l.next(); {
		case r == eof:
			return nil, l.emit(EOF)
		case r == '\n':
			l.ignore()
			return lexSpace, Token{}
		case prev == dashASCII && r == dashASCII:
			l.ignore()
			return lexSpace, Token{}
		default:
			prev = r
		}
	}
}

// lexChars accumulate characters until end of token is found.
// If the token is a reserved word return the type otherwise,
// assume a label.
func lexChars(l *Lexer) (stateFn, Token) {
	// flag that tracks if all characters are 0 - 9
	isNumber := l.input[l.start] >= '0' && l.input[l.start] <= '9'

	n := 0
	l.label[n] = l.input[l.start]
	if l.label[n] >= 'a' && l.label[0] <= 'z' {
		l.label[n] -= 32
	}
	n++

LOOP:
	for {
		switch r := l.next(); {
		case r >= 'a' && r <= 'z':
			r -= 32 // uppercase
			isNumber = false
			fallthrough
		case (r >= 'A' && r <= 'Z') || r == '_' || r == '-':
			isNumber = false
			fallthrough
		case r >= '0' && r <= '9':
			if n >= LongestKeyword || n == -1 {
				n = -1 // therefore, not a keyword
			} else {
				l.label[n] = r
				n++
			}
		default:
			l.backup()
			break LOOP
		}
	}

	if n != -1 {
		if keyword, ok := lexemes[string(l.label[0:n])]; ok {
			return lexSpace, l.emit(keyword)
		}
	}

	if isNumber {
		return lexSpace, l.emit(Number)
	}

	return lexSpace, l.emit(Label)
}
