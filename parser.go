// Package mib parses mib files.
package mib

import (
	"fmt"

	"github.com/goller/mib/tokens"
)

// Parse reads
func Parse(lexer *tokens.Lexer) {
	var name, token string
	prevType, curType := tokens.Label, tokens.Label

	for curType != tokens.EOF {
		if prevType == tokens.Continue {
			prevType = curType
		} else {
			t := lexer.NextToken()
			curType = t.Typ
			token = t.Val
			prevType = curType
		}

		switch curType {
		case tokens.End:
			// TODO: Link nodes
		case tokens.Imports:
			parseImports(lexer)
		case tokens.Exports:
			curType = parseExports(lexer)
			// need token somehow
			continue
		case tokens.Label, tokens.Integer: // TODO cleanup
		case tokens.Integer32, tokens.Uinteger32: // TODO why 2 unsigned
		case tokens.Unsigned32, tokens.Counter:
		case tokens.Gauge, tokens.Counter64:
		case tokens.Ipaddr, tokens.Netaddr:
		case tokens.Nsapaddress, tokens.Objsyntax:
		case tokens.Appsyntax, tokens.Simplesyntax:
		case tokens.Objname, tokens.Notifname:
		case tokens.KwOpaque, tokens.Timeticks:
		case tokens.EOF:
			continue
		default: // for some reason we skip the token before a macro?
			name = token
			t := lexer.NextToken()
			curType = t.Typ
			token = t.Val
			if curType == tokens.Macro {
				parseMacro(lexer, name) // I think this just is to check macro syntax?
			}
		}
		name = token
		t := lexer.NextToken()
		curType = t.Typ
		token = t.Val

		/*
		 * Handle obsolete method to assign an object identifier to a
		 * module
		 */
		if prevType == tokens.Label && curType == tokens.LeftBracket {
			for curType != tokens.RightBracket && curType != tokens.EOF {
				t := lexer.NextToken()
				curType = t.Typ
				if curType == tokens.EOF {
					// TODO: module syntax error
					return
				}
			}
			t := lexer.NextToken()
			curType = t.Typ
			token = t.Val
		}

		switch curType {
		case tokens.Definitions:
			// something to do with current module using the name variable
			for {
				t := lexer.NextToken()
				curType = t.Typ
				token = t.Val
				if curType == tokens.EOF || curType == tokens.Begin {
					break
				}
			}
		case tokens.ObjType:
			parseObjectType(lexer, name)
			// TODO handle error
		case tokens.ObjGroup:
			parseObjectGroup(lexer, name)
			// TODO: somehow return objects and check for error
		case tokens.Notifgroup: // ????? NotifType ????
			parseNotifications(lexer, name)
			// TODO: somehow return objects and check for error
		case tokens.TrapType:
			parseTrapDefinition(lexer, name)
			// TODO: somehow return objects and check for error
		case tokens.NotifType:
			parseNotification(lexer, name)
			// TODO: somehow return objects and check for error
		case tokens.Compliance:
			parseCompliance(lexer, name)
			// TODO: somehow return objects and check for error
		case tokens.AgentCap:
			parseCapabilities(lexer, name)
		case tokens.Macro:
			parseMacro(lexer, name)
		case tokens.ModuleIdentify:
			parseModuleIdentity(lexer, name)
		case tokens.ObjIdentity:
			parseObjectIdentity(lexer, name)
			// TODO: somehow return objects and check for error
		case tokens.Object:
			parseObject(lexer, name)
		case tokens.Equals:
			parseASNType(lexer, name, curType, token)
		case tokens.EOF: // break
		default:
			// TODO errors of some sort
		}

	}
}

func parseExports(lexer *tokens.Lexer) tokens.TokenType {
	tk := lexer.NextToken()
	for tk.Typ != tokens.Semicolon && tk.Typ != tokens.EOF {
		tk = lexer.NextToken()
	}
	return tk.Typ
}

func parseObjectType(lexer *tokens.Lexer, name string) {
	tk := lexer.NextToken()
	if tk.Typ != tokens.Syntax {
		// TODO error
		return
	}

	tk = lexer.NextToken()
	if tk.Typ == tokens.Object {
		tk = lexer.NextToken()
		if tk.Typ != tokens.Identifier {
			// TODO error
			return
		}
		tk.Typ = tokens.ObjIdentity // perhaps compound objid type?
	}
	if tk.Typ == tokens.Label {
		// TODO write get_tc and make sure ther esult is a label
		// this label is supposed to be a TEXTUAL-CONVENTION.
		// If this TEXTUAL-CONVENTION is not known, then, this
		// would be some sort of warning.
	}
	// set np->type to type ?
	next := lexer.NextToken()
	switch tk.Typ {
	case tokens.Sequence:
		if next.Typ == tokens.Of {
			_ = lexer.NextToken() // it is unclear to me why net-snmp skips
			next = lexer.NextToken()
		}
	case tokens.Integer, tokens.Integer32,
		tokens.Unsigned32, tokens.Uinteger32, tokens.Counter,
		tokens.Gauge, tokens.BitString, tokens.Label:
		if next.Typ == tokens.LeftBracket {
			parseEnumList(lexer)
			next = lexer.NextToken()
		} else if next.Typ == tokens.LeftParen {
			parseRanges(lexer)
			next = lexer.NextToken()
		}
	case tokens.Octetstr, tokens.KwOpaque:
		// must be (SIZE (range))
		if next.Typ != tokens.LeftParen {
			break
		}
		next = lexer.NextToken()
		if next.Typ != tokens.Size {
			break // TODO: error condition
		}
		next = lexer.NextToken()
		if next.Typ != tokens.LeftParen {
			break // TODO: error condition
		}

		parseRanges(lexer)

		next = lexer.NextToken()
		if next.Typ != tokens.RightParen {
			break // TODO: error condition
		}
	case tokens.ObjIdentity, tokens.Netaddr, tokens.Ipaddr, tokens.Timeticks, tokens.Nul, tokens.Nsapaddress, tokens.Counter64:
	default:
		// TODO: error condition
	}

	if next.Typ == tokens.Units {
		tk = lexer.NextToken()
		if tk.Typ != tokens.Quotestring {
			//  TODO: error units required
		}
		// TODO: save units
		//  units := tk.Val
		next = lexer.NextToken()
	}

	if next.Typ != tokens.Access {
		// TODO error rquires access
	}

	tk = lexer.NextToken()
	switch tk.Typ {
	case tokens.ReadOnly, tokens.ReadWrite:
	case tokens.WriteOnly, tokens.NoAccess:
	case tokens.ReadCreate, tokens.Accnotify:
	default:
		// TODO: error with bad access type
	}

	// TODO: set access of node to access type
	tk = lexer.NextToken()
	if tk.Typ != tokens.Status {
		// TODO: error should be status
	}

	tk = lexer.NextToken()
	switch tk.Typ {
	case tokens.Mandatory, tokens.Current, tokens.KwOptional, tokens.Obsolete, tokens.Deprecated:
	default:
		// TODO: error with bad status
	}
	// TODO: set the status on the node

	tk = lexer.NextToken()
	for tk.Typ != tokens.Equals && tk.Typ != tokens.EOF {
		switch tk.Typ {
		case tokens.Description:
		case tokens.Reference:
		case tokens.Index:
		case tokens.Augments:
		case tokens.Defval:
		case tokens.NumEntries:
			// I can't find usage of this anywhere
			if tossObjectIdentifier(lexer) != tokens.ObjIdentity {
				// TODO: return bad object id
			}
		default:
			// TODO: errorbad optional clases
		}
		tk = lexer.NextToken()
	}

	if tk.Typ != tokens.Equals {
		// TODO: error expected equals
	}

	parseObjectID(lexer, name)
}

/*
 * Parses an enumeration list of the form:
 *        { label(value) label(value) ... }
 * The initial { has already been parsed.
 */
func parseEnumList(lexer *tokens.Lexer) {
Loop:
	for {
		switch tok := lexer.NextToken(); {
		case tok.Typ == tokens.EOF || tok.Typ == tokens.None:
			// TODO: This is an error; probably return something about expected }
			break Loop
		case tok.Typ == tokens.RightBracket:
			break Loop
		case tok.Typ == tokens.Label || tok.Typ == tokens.Deprecated:
			/*
				Some enums use "deprecated" to indicate a no longer value label.

				For example, IP-MIB's IpAddressStatusTC:
				SYNTAX     INTEGER {
					preferred(1),
					deprecated(2), -- here this is deprecated
					invalid(3),
					inaccessible(4),
					unknown(5),
					tentative(6),
					duplicate(7),
					optimistic(8)
				}
			*/
			// TODO: requires LeftParen, Number, RightParen
			// TODO: Collect up all the enumerations
			// Perhaps something like struct{Label: tok.Val, Value: ParseInt(val.Val)}
			// net-snmp does _not_ check for errors in the string to long conversion and uses 0 as value.
			next := lexer.NextToken()
			if next.Typ != tokens.LeftParen {
				// TODO: Error about expected (
			}
			next = lexer.NextToken()
			if next.Typ != tokens.Hex && next.Typ != tokens.Binary {
				// TODO: Error about expected Number
				// TODO: here we would ParseInt()
				// TODO: check on lexer for Number type
			}

			next = lexer.NextToken()
			if next.Typ != tokens.RightParen {
				// TODO: Error about expected )
			}
		}
	}
}

/*
Assumes ( has already been parsed
A. NUMBER) -- hex/binary/number should be a long
B. LOW..HIGH)
C. A or B | A or B | etc) -- the pipe can repeat
D. SIZE(A))
E. SIZE(B))
F. SIZE(C))
*/
func parseRanges(lexer *tokens.Lexer) {

}

/*
[OBJECTS {(goto) LABEL(goto) [, LABEL(goto) }(goto)] -- netsnmp has this optional
STATUS(goto) CURRENT(goto) | DEPRECATED(goto) | OBSOLETE(goto)
DESCRIPTION(goto) QUOTEDSTRING(return)
[REFERENCE QUOTEDSTRING]
EQUALS(goto)

until EQUALS|EOF
then parseObjectID
*/
func parseObjectGroup(lexer *tokens.Lexer, name string) (*OIDSegments, error) {
	err := required(lexer, tokens.Objects, tokens.LeftBracket, tokens.Label)
	if err != nil {
		return nil, err
	}

	err = untilRepeat(lexer, tokens.RightBracket, tokens.Comma, tokens.Label)
	if err != nil {
		return nil, err
	}

	if err = parseStatus(lexer); err != nil {
		return nil, err
	}

	if err = required(lexer, tokens.Description, tokens.Quotestring); err != nil {
		return nil, err
	}

	tk := lexer.NextToken()
	if tk.Typ == tokens.Reference {
		if err = required(lexer, tokens.Quotestring); err != nil {
			return nil, err
		}
		tk = lexer.NextToken()
	}

	if tk.Typ != tokens.Equals {
		return nil, NewExpectedTokenError(tokens.Token{
			Typ: tokens.Equals,
			Val: "::=",
		})
	}

	return parseObjectID(lexer, name)
}

func untilRepeat(lexer *tokens.Lexer, until tokens.TokenType, repeat ...tokens.TokenType) error {
	tk := lexer.NextToken()
	for {
		if tk.Typ == until {
			return nil
		}
		if tk.Typ == tokens.EOF || tk.Typ == tokens.None {
			return fmt.Errorf("missing token type %v", until) // TODO: stringify the contants
		}
		for i := range repeat {
			if tk.Typ != repeat[i] {
				return fmt.Errorf("missing token type %v; received %v", repeat[i], tk) // TODO: stringify the contants
			}
			tk = lexer.NextToken()
		}
	}
}

/*
 Status ::=
              "current"
            | "deprecated"
            | "obsolete"
*/
func parseStatus(lexer *tokens.Lexer) error {
	tk := lexer.NextToken()
	if tk.Typ != tokens.Status {
		return NewExpectedTokenError(tokens.Token{
			Typ: tokens.Status,
			Val: "STATUS",
		})
	}

	tk = lexer.NextToken()
	switch tk.Typ {
	case tokens.Current, tokens.Deprecated, tokens.Obsolete:
	default:
		return NewExpectedTokenError(tokens.Token{
			Typ: tokens.Current,
			Val: "current",
		})
	}
	return nil
}

/*
until EQUALS | EOF

(any order and any token is ok)
DESCRIPTION QUOTEDSTRING(return)
REFERENCE QUOTEDSTRING(return)
OBJECTS parseVarBinds(return)
then parseObjectID
*/
func parseNotification(lexer *tokens.Lexer, name string) (*Trap, error) {
	var (
		variables []string
		err       error
	)

	tk := lexer.NextToken()
	if tk.Typ == tokens.Objects {
		variables, err = parseVarBinds(lexer)
		if err != nil {
			return nil, err
		}
		tk = lexer.Next()
	}

	if tk.Typ != tokens.Status {
		return nil, NewExpectedTokenError(tokens.Token{
			Typ: tokens.Status,
			Val: "STATUS",
		})
	}

	tk = lexer.NextToken()
	switch tk.Typ {
	case tokens.Current, tokens.Deprecated, tokens.Obsolete:
	default:
		return nil, NewExpectedTokenError(tokens.Token{
			Typ: tokens.Current,
			Val: "current",
		})
	}

	if err = parseStatus(lexer); err != nil {
		return nil, err
	}
	if err = required(lexer, tokens.Description, tokens.Quotestring); err != nil {
		return nil, err
	}

	tk := lexer.NextToken()
	if tk.Typ == tokens.Reference {
		tk = lexer.NextToken()
		if tk.Typ != tokens.Quotestring {
			return nil, NewExpectedTokenError(tokens.Token{
				Typ: tokens.Quotestring,
				Val: "reference string",
			})
		}
		tk = lexer.NextToken()
	}

}

/*
until EQUALS | EOF

(any order and any token is ok)
DESCRIPTION QUOTEDSTRING(return)
REFERENCE QUOTEDSTRING(return)
OBJECTS parseVarBinds(return)
then parseObjectID
*/
// TODO: these objects need to be recorded in a new type like a NotificationGroup
func parseNotifications(lexer *tokens.Lexer, name string) (*OIDSegments, error) {
	// TODO: parseVarBinds
	err := required(lexer, tokens.Notifications, tokens.LeftBracket, tokens.Label)
	if err != nil {
		return nil, err
	}

	err = untilRepeat(lexer, tokens.RightBracket, tokens.Comma, tokens.Label)
	if err != nil {
		return nil, err
	}

	if err = parseStatus(lexer); err != nil {
		return nil, err
	}

	if err = required(lexer, tokens.Description, tokens.Quotestring); err != nil {
		return nil, err
	}

	tk := lexer.NextToken()
	if tk.Typ == tokens.Reference {
		if err = required(lexer, tokens.Quotestring); err != nil {
			return nil, err
		}
		tk = lexer.NextToken()
	}

	if tk.Typ != tokens.Equals {
		return nil, NewExpectedTokenError(tokens.Token{
			Typ: tokens.Equals,
			Val: "::=",
		})
	}

	return parseObjectID(lexer, name)
}

/*
LEFTBRACKET(return)
until RIGHTBRACKET|EOF
store if LABEL or with SYNTAX_MASK
*/
func parseVarBinds(lexer *tokens.Lexer) ([]string, error) {
	tk := lexer.NextToken()
	if tk.Typ != tokens.LeftBracket {
		return nil, NewExpectedTokenError(tokens.Token{
			Typ: tokens.LeftBracket,
			Val: "{",
		})
	}

	variables := []string{}
Loop:
	for {
		switch tk = lexer.NextToken(); {
		case tk.Typ.IsSyntax(): // I have not see why this matters in any MIB
			variables = append(variables, tk.Val)
		case tk.Typ == tokens.Label:
			variables = append(variables, tk.Val)
		case tk.Typ == tokens.RightBracket:
			break Loop
		case tk.Typ == tokens.EOF || tk.Typ == tokens.None:
			return nil, NewEOFError(tokens.Token{
				Typ: tokens.RightBracket,
				Val: "}",
			})
		}
	}
	return variables, nil
}

/*
until EQUALS | EOF
(any order and any token is ok)

DESCRIPTION QUOTEDSTRING(return)
REFERENCE QUOTEDSTRING(return)
VARIABLES parseVarBinds(return)
ENTERPRISE

then

NUMBER(return) this is the oid
*/
func parseTrapDefinition(lexer *tokens.Lexer, name string) (*Trap, error) {
	var (
		enterprise string
		variables  []string
	)
Loop:
	for {
		tk := lexer.NextToken()
		switch tk.Typ {
		case tokens.Enterprise:
			// net-snmp has {LABEL}, but, I haven't seen that anywhere
			tk := lexer.NextToken()
			if tk.Typ != tokens.Label {
				return nil, NewExpectedTokenError(tokens.Token{
					Typ: tokens.Label,
					Val: "enterprise id",
				})
			}
			enterprise = tk.Val
		case tokens.Variables:
			var err error
			variables, err = parseVarBinds(lexer)
			if err != nil {
				return nil, err
			}
		case tokens.Description:
			if err := required(lexer, tokens.Quotestring); err != nil {
				return nil, err
			}
		case tokens.Reference:
			if err := required(lexer, tokens.Quotestring); err != nil {
				return nil, err
			}
		case tokens.Equals:
			break Loop
		case tokens.EOF, tokens.None:
			return nil, NewEOFError(tokens.Token{
				Typ: tokens.Number,
				Val: "number",
			})
		}
	}
	tk := lexer.NextToken()
	if tk.Typ != tokens.Number { // TODO: hex/binary
		return nil, NewExpectedTokenError(tokens.Token{
			Typ: tokens.Number,
			Val: "number",
		})
	}

	return &Trap{
		OIDSegments: OIDSegments{
			Name: name,
			Segments: []OIDSegment{
				{
					Name: enterprise,
				},
				{
					Number: tk.Val,
				},
			},
		},
		VarBinds: variables,
	}, nil
}

/*
STATUS(goto)
CURRENT|DEPRECATED|OBSOLETE(goto)
DESCRIPTION(GOTO)
QUOTEDSTRING(goto)
[REFERENCE QUOTEDSTRING(goto)]
MODULE(goto)
	while MODULE
		[LABEL]
		[MANDATORYGROUPS LEFTBRACKET(goto) [LABEL(goto) dowhile(COMMA)] RIGHTBRACKET(goto) ]
		while GROUP || OBJECT
			if GROUP then LABEL(goto)
			if OBJECT then LABEL(goto)
				[SYNTAX(eatSyntax)] from eat or  [WRSYNTAX(eatSyntax)]
				 from eat or [MINACCESS MAXTOKEN|NOACCESS|ACCNOTIFY|READONLY|WRITEONLY|READCREATE|READWRITE otherwise goto  ]
			DESCRIPTION(goto)
			QUOTEDSTRING(goto)
goto:
	until EQUALS | EOF
    then parseObjectID
*/
func parseCompliance(lexer *tokens.Lexer, name string) {
}

/*
	INTEGER | INTEGER32 |
	UINTEGER32 | UNSIGNED32 |
	COUNTER | GAUGE |
	BITSTRING | LABEL then LEFTBRACKET(parseEnumList) | LEFTPAREN(parseRange)
	return next token
OR
	OCTETSTR | KW_OPAQUE
	[LEFTPAREN [SIZE [LEFTPAREN(parseRanges)[RIGHTPAREN]]](error)
	return next token
OR
	OBJID|NETADDR|IPADDR|TIMETICKS|NUL|NSAPADDRESS|COUNTER64
		return next token
OR error and return next token
*/
func eatSyntax(lexer *tokens.Lexer) tokens.Token {
	return tokens.Token{}
}

/*
PRODREL(goto)
QUOTEDSTRING(goto)
STATUS(goto)
CURRENT|OBSOLETE(goto)  DEPRECATED seems to not be used?
DESCRIPTION(GOTO)
QUOTEDSTRING(goto)
[REFERENCE QUOTEDSTRING(goto)]
	while SUPPORTS
		LABEL(goto)
		INCLUDES(goto)
		LEFTBRACKET(goto)  {LABEL(goto) dowhile(COMMA)]} RIGHTBRACKET(goto)

		while VARIATION
			LABEL(goto)
			[SYNTAX(eatSyntax)]
			from eat or  [WRSYNTAX(eatSyntax)]
			from eat or ACCESS ACCNOTIFY|READONLY|READWRITE|READCREATE|WRITEONLY|NOTIMPL otherwise goto  ]
			or CREATEREQ LEFTBRACKET(goto)  {LABEL(goto) dowhile(COMMA)]} RIGHTBRACKET(goto)
			of if DEFVAL while pairs of LEFT/RIGHTBRACKET and tokens between

			DESCRIPTION(goto)
			QUOTEDSTRING(goto)
    EQUALS
goto:
	until EQUALS | EOF
    then parseObjectID
*/
func parseCapabilities(lexer *tokens.Lexer, name string) {
}

/*
	LASTUPDATED(goto) QUOTEDSTRING(goto)
	ORGANIZATION(goto) QUOTEDSTRING(goto)
	CONTACTINFO(goto) QUOTEDSTRING(goto)
	DESCRIPTION(goto) QUOTEDSTRING(goto)

	while REVISION
		QUOTEDSTRING(goto)
		DESCRIPTION(goto)
		QUOTEDSTRING(goto)
	EQUALS
goto:
	until EQUALS | EOF
    then parseObjectID
*/
func parseModuleIdentity(lexer *tokens.Lexer, name string) (segments *OIDSegments, err error) {
	var tk tokens.Token
	// After attempting to find all required fields of MODULE-IDENTITY,
	// swallow all tokens until equals or EOF; afterwhich, we recover
	// the object id.
	defer func() {
		if err != nil {
			return //TODO: optimistically continue?
		}
		for {
			if tk.Typ == tokens.Equals {
				segments, err = parseObjectID(lexer, name)
				return
			}
			if tk.Typ == tokens.EOF || tk.Typ == tokens.None {
				err = NewEOFError(tokens.Token{
					Typ: tokens.Equals,
					Val: "=",
				})
				return
			}
			tk = lexer.NextToken()
		}
	}()

	notation := []tokens.TokenType{
		tokens.LastUpdated, tokens.Quotestring,
		tokens.Organization, tokens.Quotestring,
		tokens.ContactInfo, tokens.Quotestring,
		tokens.Description, tokens.Quotestring,
	}

	if err = required(lexer, notation...); err != nil {
		return
	}

	// Revision statements are optional
	for tk = lexer.NextToken(); tk.Typ == tokens.Revision; tk = lexer.NextToken() {
		err = required(lexer, tokens.Quotestring, tokens.Description, tokens.Quotestring)
		if err != nil {
			return
		}
	}
	return
}

func required(lexer *tokens.Lexer, reqs ...tokens.TokenType) error {
	for i := range reqs {
		tk := lexer.NextToken()
		if tk.Typ != reqs[i] {
			return fmt.Errorf("missing token type %v; received %v", reqs[i], tk) // TODO: stringify the contants
		}
	}
	return nil
}

/*
[OBJECTS {(goto) LABEL(goto) [, LABEL(goto) }(goto)] -- netsnmp has this optional
STATUS(goto) CURRENT(goto) | DEPRECATED(goto) | OBSOLETE(goto)
DESCRIPTION(goto) QUOTEDSTRING(return)
[REFERENCE QUOTEDSTRING]
EQUALS(goto)

until EQUALS|EOF
then parseObjectID
*/
func parseObjectIdentity(lexer *tokens.Lexer, name string) {
	parseObjectGroup(lexer, name)
}

func parseObject(lexer *tokens.Lexer, name string) (*OIDSegments, error) {
	notation := []tokens.TokenType{
		tokens.Identifier, tokens.Equals,
	}

	if err := required(lexer, notation...); err != nil {
		return nil, err
	}

	return parseObjectID(lexer, name)
}

/*
any EOF return
until EQUALS
until BEGIN
until END
*/
func parseMacro(lexer *tokens.Lexer, name string) {

}

/*
until SEMI
	LABEL (load?)
	FROM
*/
func parseImports(lexer *tokens.Lexer) ([]Import, error) {
	imports := []Import{}
	imp := Import{}
Loop:
	for {
		tk := lexer.NextToken()
		switch tk.Typ {
		case tokens.Semicolon:
			break Loop
		case tokens.Label:
			imp.Types = append(imp.Types, tk.Val)
		case tokens.From:
			tk = lexer.NextToken()
			if tk.Typ != tokens.Label {
				return nil, NewExpectedTokenError(tokens.Token{
					Typ: tokens.Label,
					Val: "FROM label",
				})
			}
			// sometimes imports are only keywords, and thus,
			// we don't need to account for them as we have hardcoded
			// keywords in the lexer.
			if len(imp.Types) == 0 {
				continue
			}
			imp.From = tk.Val
			imports = append(imports, imp)
			imp = Import{}
		case tokens.EOF:
			return nil, NewEOFError(
				tokens.Token{
					Typ: tokens.Semicolon,
					Val: ";",
				})
		}
	}
	return imports, nil
}

/*
TODO look at how to pass next token from here
	SEQUENCE|CHOICE
		until RIGHTBRACKET || depth > 0
			if LEFTBRACKET depth++
			if RIGHTBRACKET depth--
or
	LEFTBRACKET(parseObjectID) -- I don't see this syntax _anywhere_
or
	LEFTSQBRACK
		until RIGHTSQBRACK
		[IMPLICIT]
		OCTETSTR | INTEGER -- can ASN.1 types be anything else?
			until RIGHTPAREN || depth > 0
				if LEFTPAREN depth++
				if RIGHTPAREN depth--
or
	CONVENTION
		until SYNTAX
			[DISPLAYHINT QUOTEDSTRING(error)]
			[DESCRIPTION QUOTEDSTRING(error)]
		[OBJECT IDENTIFIER(error)] type = OBJID ???
		[LABEL(type = get_tc)]
		[LEFTPAREN(parseRanges)(return)]
		[LEFTBRACKET(parseEnumList)(return)
or
		[OBJECT IDENTIFIER(error)] type = OBJID ???
		[LABEL(type = get_tc)]
		[LEFTPAREN(parseRanges)(return)]
		[LEFTBRACKET(parseEnumList)(return)

*/
func parseASNType(lexer *tokens.Lexer, name string, curType tokens.TokenType, token string) {}

/*
LEFTBRACKET(error)
until RIGHTBRACKET
	NUMBER
	or LABEL [LEFTPAREN NUMBER RIGHTPAREN]
*/
func parseObjectID(lexer *tokens.Lexer, name string) (*OIDSegments, error) {
	tk := lexer.NextToken()
	if tk.Typ != tokens.LeftBracket {
		return nil, NewExpectedTokenError(tokens.Token{
			Typ: tokens.LeftBracket,
			Val: "{",
		},
		)
	}

	oid := &OIDSegments{
		Name:     name,
		Segments: make([]OIDSegment, 0, 2), // {label number} most common
	}
	depth := 0
	for {
		switch {
		case tk.Typ == tokens.LeftBracket:
			depth++
		case tk.Typ == tokens.RightBracket:
			depth--
		case tk.Typ == tokens.EOF:
			return nil, NewEOFError(tokens.Token{
				Typ: tokens.RightBracket,
				Val: "}",
			})
		case tk.Typ == tokens.Number:
			oid.Segments = append(oid.Segments, OIDSegment{
				Number: tk.Val,
			})
		case tk.Typ == tokens.Label:
			label := tk.Val
			tk = lexer.NextToken()
			if tk.Typ != tokens.LeftParen {
				oid.Segments = append(oid.Segments, OIDSegment{
					Name: label,
				})
				continue
			}
			tk = lexer.NextToken()
			if tk.Typ != tokens.Number {
				oid.Segments = append(oid.Segments, OIDSegment{
					Name: label,
				})
				continue // TODO: probably an error
			}

			num := tk.Val

			tk = lexer.NextToken()
			if tk.Typ != tokens.RightParen {
				oid.Segments = append(oid.Segments, OIDSegment{
					Name:   label,
					Number: num,
				})
				continue // TODO: probably an error
			}
			oid.Segments = append(oid.Segments, OIDSegment{
				Name:   label,
				Number: num,
			})
		}

		if depth == 0 {
			return oid, nil
		}
		tk = lexer.NextToken()
	}
}

/*

LEFTBRACKET
until RIGHTBRACKET || depth > 0
	if LEFTBRACKET depth++
	if RIGHTBRACKET depth--
*/
func tossObjectIdentifier(lexer *tokens.Lexer) tokens.TokenType {
	tk := lexer.NextToken()
	if tk.Typ != tokens.LeftBracket {
		return tokens.None
	}

	depth := 0
	for {
		switch {
		case tk.Typ == tokens.LeftBracket:
			depth++
		case tk.Typ == tokens.RightBracket:
			depth--
		case tk.Typ == tokens.EOF:
			return tokens.None
		}

		if depth == 0 {
			return tokens.ObjIdentity
		}
		tk = lexer.NextToken()
	}
}
