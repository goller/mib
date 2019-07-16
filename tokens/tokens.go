package tokens

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// tokenType is one of the specific MIB token types.
type tokenType int

// token is a string in a MIB file with an identified meaning.
type token struct {
	typ  tokenType
	pos  Pos
	val  string
	line int
}

func (t token) String() string {
	switch {
	case t.typ == tokenEOF:
		return "EOF"
	case t.typ == tokenError:
		return t.val
	case t.typ > tokenKeyword:
		return fmt.Sprintf("<%s>", t.val)
	case len(t.val) > 10:
		return fmt.Sprintf("%.10q...", t.val)
	}

	return fmt.Sprintf("%q", t.val)
}

const (
	tokenError tokenType = iota
	tokenLeftParen
	tokenRightParen
	tokenLeftBracket
	tokenRightBracket
	tokenLeftSquareBracket
	tokenRightSquareBracket
	tokenSemicolon
	tokenComma
	tokenBar
	tokenRange
	tokenLabel
	tokenEquals
	tokenEOF
	tokenKeyword
	tokenObsolete
	tokenKwOpaque
	tokenKwOptional
	tokenLastUpdated
	tokenOrganization
	tokenContactInfo
	tokenModuleIdentify
	tokenCompliance
	tokenDefinitions
	tokenEnd
	tokenAugments
	tokenNoAccess
	tokenWriteOnly
	tokenNsapaddress
	tokenUnits
	tokenReference
	tokenNumEntries
	tokenBitstring
	tokenContinue
	tokenBitString
	tokenCounter64
	tokenTimeticks
	tokenNotifType
	tokenObjGroup
	tokenObjIdentity
	tokenIdentifier
	tokenObject
	tokenNetaddr
	tokenGauge
	tokenUnsigned32
	tokenReadWrite
	tokenReadCreate
	tokenOctetstr
	tokenOf
	tokenSequence
	tokenNul
	tokenIpaddr
	tokenBinary
	tokenHex
	tokenUinteger32
	tokenInteger
	tokenInteger32
	tokenCounter
	tokenReadOnly
	tokenDescription
	tokenIndex
	tokenDefval
	tokenDeprecated
	tokenSize
	tokenAccess
	tokenMandatory
	tokenCurrent
	tokenStatus
	tokenSyntax
	tokenObjType
	tokenTrapType
	tokenEnterprise
	tokenBegin
	tokenImports
	tokenExports
	tokenAccnotify
	tokenConvention
	tokenNotifgroup
	tokenDisplayHint
	tokenFrom
	tokenAgentCap
	tokenMacro
	tokenImplied
	tokenSupports
	tokenIncludes
	tokenVariation
	tokenRevision
	tokenNotImpl
	tokenObjects
	tokenNotifications
	tokenModule
	tokenMinAccess
	tokenProdRel
	tokenWrSyntax
	tokenCreateReq
	tokenMandatoryGroups
	tokenGroup
	tokenChoice
	tokenImplicit
	tokenObjsyntax
	tokenSimplesyntax
	tokenAppsyntax
	tokenObjname
	tokenNotifname
	tokenVariables
	tokenQuotestring
)

var lexemes = map[string]tokenType{
	"OBSOLETE":              tokenObsolete,
	"OPAQUE":                tokenKwOpaque,
	"OPTIONAL":              tokenKwOptional,
	"LAST-UPDATED":          tokenLastUpdated,
	"ORGANIZATION":          tokenOrganization,
	"CONTACT-INFO":          tokenContactInfo,
	"MODULE-IDENTITY":       tokenModuleIdentify,
	"MODULE-COMPLIANCE":     tokenCompliance,
	"DEFINITIONS":           tokenDefinitions,
	"END":                   tokenEnd,
	"AUGMENTS":              tokenAugments,
	"NOT-ACCESSIBLE":        tokenNoAccess,
	"WRITE-ONLY":            tokenWriteOnly,
	"NSAPADDRESS":           tokenNsapaddress,
	"UNITS":                 tokenUnits,
	"REFERENCE":             tokenReference,
	"NUM-ENTRIES":           tokenNumEntries,
	"BITSTRING":             tokenBitstring,
	"BIT":                   tokenContinue,
	"BITS":                  tokenBitString,
	"COUNTER64":             tokenCounter64,
	"TIMETICKS":             tokenTimeticks,
	"NOTIFICATION-TYPE":     tokenNotifType,
	"OBJECT-GROUP":          tokenObjGroup,
	"OBJECT-IDENTITY":       tokenObjIdentity,
	"IDENTIFIER":            tokenIdentifier,
	"OBJECT":                tokenObject,
	"NETWORKADDRESS":        tokenNetaddr,
	"GAUGE":                 tokenGauge,
	"GAUGE32":               tokenGauge,
	"UNSIGNED32":            tokenUnsigned32,
	"READ-WRITE":            tokenReadWrite,
	"READ-CREATE":           tokenReadCreate,
	"OCTETSTRING":           tokenOctetstr,
	"OCTET":                 tokenContinue,
	"OF":                    tokenOf,
	"SEQUENCE":              tokenSequence,
	"NULL":                  tokenNul,
	"IPADDRESS":             tokenIpaddr,
	"UINTEGER32":            tokenUinteger32,
	"INTEGER":               tokenInteger,
	"INTEGER32":             tokenInteger32,
	"COUNTER":               tokenCounter,
	"COUNTER32":             tokenCounter,
	"READ-ONLY":             tokenReadOnly,
	"DESCRIPTION":           tokenDescription,
	"INDEX":                 tokenIndex,
	"DEFVAL":                tokenDefval,
	"DEPRECATED":            tokenDeprecated,
	"SIZE":                  tokenSize,
	"MAX-ACCESS":            tokenAccess,
	"ACCESS":                tokenAccess,
	"MANDATORY":             tokenMandatory,
	"CURRENT":               tokenCurrent,
	"STATUS":                tokenStatus,
	"SYNTAX":                tokenSyntax,
	"OBJECT-TYPE":           tokenObjType,
	"TRAP-TYPE":             tokenTrapType,
	"ENTERPRISE":            tokenEnterprise,
	"BEGIN":                 tokenBegin,
	"IMPORTS":               tokenImports,
	"EXPORTS":               tokenExports,
	"ACCESSIBLE-FOR-NOTIFY": tokenAccnotify,
	"TEXTUAL-CONVENTION":    tokenConvention,
	"NOTIFICATION-GROUP":    tokenNotifgroup,
	"DISPLAY-HINT":          tokenDisplayHint,
	"FROM":                  tokenFrom,
	"AGENT-CAPABILITIES":    tokenAgentCap,
	"MACRO":                 tokenMacro,
	"IMPLIED":               tokenImplied,
	"SUPPORTS":              tokenSupports,
	"INCLUDES":              tokenIncludes,
	"VARIATION":             tokenVariation,
	"REVISION":              tokenRevision,
	"NOT-IMPLEMENTED":       tokenNotImpl,
	"OBJECTS":               tokenObjects,
	"NOTIFICATIONS":         tokenNotifications,
	"MODULE":                tokenModule,
	"MIN-ACCESS":            tokenMinAccess,
	"PRODUCT-RELEASE":       tokenProdRel,
	"WRITE-SYNTAX":          tokenWrSyntax,
	"CREATION-REQUIRES":     tokenCreateReq,
	"MANDATORY-GROUPS":      tokenMandatoryGroups,
	"GROUP":                 tokenGroup,
	"CHOICE":                tokenChoice,
	"IMPLICIT":              tokenImplicit,
	"OBJECTSYNTAX":          tokenObjsyntax,
	"SIMPLESYNTAX":          tokenSimplesyntax,
	"APPLICATIONSYNTAX":     tokenAppsyntax,
	"OBJECTNAME":            tokenObjname,
	"NOTIFICATIONNAME":      tokenNotifname,
	"VARIABLES":             tokenVariables,
	"QUOTEDSTRING":          tokenQuotestring,
}

const eof = -1

// Trimming spaces.
const (
	spaceChars = " \t\r\n" // These are the space characters defined by Go itself.
)

// Pos represents a byte position in the original input text from which
// this token was parsed.
type Pos int

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name    string     // the name of the input; used only for error reports
	input   string     // string to scan
	state   stateFn    // the next lexing function to enter
	pos     Pos        // current position in the input
	start   Pos        // start position of this item
	width   Pos        // width of last rune read from input
	lastPos Pos        // position of most recent item returned by nextItem
	tokens  chan token // channel of scanned tokens
	line    int        // 1+number of newlines seen
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// emit passes an item back to the client.
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.start, l.input[l.start:l.pos], l.line}
	// Some items contain text internally. If so, count their newlines.
	/*switch t {
	  case itemText:
	      l.line += strings.Count(l.input[l.start:l.pos], "\n")
	  }*/
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{tokenError, l.start, fmt.Sprintf(format, args...), l.line}
	return nil
}

// nextToken returns the next token from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextToken() token {
	item := <-l.tokens
	l.lastPos = item.pos
	return item
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make(chan token),
		line:   1,
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexSpace; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.tokens)
}

const (
	comment = "--"
)

// lexText scans until an opening action delimiter, "{{".
func lexText(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], "--") {
		l.ignore()
		return lexComment
	}

	switch r := l.next(); {
	case r == eof:
		break
	case r == '"':
		return lexQuotedString
	case r == '\'':
		return lexNumberLiteral
	case r == '(':
		l.emit(tokenLeftParen)
		return lexSpace
		// TODO: count depth
	case r == ')':
		l.emit(tokenRightParen)
		// TODO: count depth
		return lexSpace
	case r == '{':
		l.emit(tokenLeftBracket)
		// TODO: count depth
		return lexSpace
	case r == '}':
		l.emit(tokenRightBracket)
		// TODO: count depth
		return lexSpace
	case r == '[':
		l.emit(tokenLeftSquareBracket)
		// TODO: count depth
		return lexSpace
	case r == ']':
		l.emit(tokenRightSquareBracket)
		// TODO: count depth
		return lexSpace
	case r == ';':
		l.emit(tokenSemicolon)
		return lexSpace
	case r == ',':
		l.emit(tokenComma)
		return lexSpace
	case r == '|':
		l.emit(tokenBar)
		return lexSpace
	case r == '.':
		return lexRange
	case r == ':':
		return lexEquals
	case r <= unicode.MaxASCII && unicode.IsPrint(r):
		return lexChars
	default:
		return l.errorf("unrecognized character: %#U", r)
	}

	l.emit(tokenEOF)
	return nil
}

func lexQuotedString(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\r', '\n':
			continue
		case eof:
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	l.emit(tokenQuotestring)
	return lexSpace
}

func lexNumberLiteral(l *lexer) stateFn {
	const (
		binary uint = 1 << iota
		hex
		unknown
	)
	var numType uint = binary

Loop:
	for {
		switch r := l.next(); {
		case r == '0' || r == '1':
			numType |= binary
		case ('0' <= r && r <= '9') || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F'):
			numType |= hex
		case r == eof:
			return l.errorf("unterminated literal string")
		case r == '\'':
			break Loop
		default:
			numType |= unknown
		}
	}

	switch r := unicode.ToUpper(l.next()); {
	case r == 'B' && numType&binary == binary:
		l.emit(tokenBinary)
	case r == 'H' && numType&unknown == 0:
		l.emit(tokenHex)
	case r == eof:
		l.emit(tokenLabel)
		l.emit(tokenEOF)
	default:
		l.emit(tokenLabel)
	}

	return lexSpace
}

func lexSpace(l *lexer) stateFn {
	for unicode.IsSpace(l.peek()) {
		_ = l.next()
	}
	l.ignore()
	return lexText
}

// lexEquals searches for ::= otherwise it assumes a label; assumes first `:`
// already consumed.
func lexEquals(l *lexer) stateFn {
	if l.next() != ':' {
		l.backup()
		l.emit(tokenLabel)
		return lexSpace
	}
	if l.next() != '=' {
		l.backup()
		l.emit(tokenLabel)
		return lexSpace
	}
	l.emit(tokenEquals)
	return lexSpace
}

// lexEquals searches for .. otherwise it assumes a label; assumes first `.`
// already consumed.
func lexRange(l *lexer) stateFn {
	if l.next() == '.' {
		l.emit(tokenRange)
		return lexSpace
	}
	l.backup()
	l.emit(tokenLabel)
	return lexSpace
}

// lexComment treats the rest of the line or until another '--' as a comment;
// the left comment marker is known to be present.
func lexComment(l *lexer) stateFn {
	l.pos += Pos(len(comment))
	var prev rune
	for {
		switch r := l.next(); {
		case r == eof:
			l.emit(tokenEOF)
			return nil
		case r == '\n':
			l.ignore()
			return lexSpace
		case prev == '-' && r == '-':
			l.ignore()
			return lexSpace
		default:
			prev = r
		}
	}
}

// lexChars accumulate characters until end of token is found.
// If the token is a reserved word return the type otherwise,
// assume a label.
func lexChars(l *lexer) stateFn {
	switch r := l.next(); {
	case (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= '0' && r <= '9') ||
		r == '_' ||
		r == '-':
	default:
		l.emit(tokenLabel)
		return lexSpace
	}

LOOP:
	for {
		switch r := l.next(); {
		case (r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') ||
			r == '_' ||
			r == '-':
		default:
			l.backup()
			break LOOP
		}
	}
	if keyword, ok := lexemes[strings.ToUpper(l.input[l.start:l.pos])]; ok {
		l.emit(keyword)
		return lexSpace
	}

	l.emit(tokenLabel)
	return lexSpace
}