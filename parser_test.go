// Package mib parses mib files.
package mib

import (
	"reflect"
	"testing"

	"github.com/goller/mib/tokens"
)

func Test_tossObjectIdentifier(t *testing.T) {
	tests := []struct {
		name string
		mib  string
		want tokens.TokenType
	}{
		{
			name: "empty obj identifier is not an error",
			mib:  "{}",
			want: tokens.ObjIdentity,
		},
		{
			name: "no left bracket returns none",
			mib:  "OBJECT",
			want: tokens.None,
		},
		{
			name: "left with no right bracket returns none",
			mib:  "{",
			want: tokens.None,
		},
		{
			name: "deep brackets return object identity",
			mib:  "{{{}}}",
			want: tokens.ObjIdentity,
		},
		{
			name: "unbalanced deep brackets return none",
			mib:  "{{iso 3}",
			want: tokens.None,
		},
		{
			name: "well-formed object id returns object identity",
			mib:  "{iso 2}",
			want: tokens.ObjIdentity,
		},
		{
			name: "extra tokens are not processed",
			mib:  "{iso 2} ::",
			want: tokens.ObjIdentity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := tokens.NewLexer(tt.mib)
			if got := tossObjectIdentifier(lexer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tossObjectIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseObjectID(t *testing.T) {
	tests := []struct {
		name    string
		label   string
		mib     string
		want    *OIDSegments
		wantErr bool
	}{
		{
			name:  "label, number",
			label: "org",
			mib:   "{iso 3}",
			want: &OIDSegments{
				Name: "org",
				Segments: []OIDSegment{
					{
						Name: "iso",
					},
					{
						Number: "3",
					},
				},
			},
		},
		{
			name:  "label, label with number, label with number, number",
			label: "internet",
			mib:   "{iso org(3) dod(6) 1}",
			want: &OIDSegments{
				Name: "internet",
				Segments: []OIDSegment{
					{
						Name: "iso",
					},
					{
						Name:   "org",
						Number: "3",
					},
					{
						Name:   "dod",
						Number: "6",
					},
					{
						Number: "1",
					},
				},
			},
		},
		{
			name:  "label, number, number",
			label: "rptrInfoHealth",
			mib:   " { snmpDot3RptrMgt 0 4 }",
			want: &OIDSegments{
				Name: "rptrInfoHealth",
				Segments: []OIDSegment{
					{
						Name: "snmpDot3RptrMgt",
					},
					{
						Number: "0",
					},
					{
						Number: "4",
					},
				},
			},
		},
		{
			name:  "number",
			label: "iso",
			mib:   "{ 1 }",
			want: &OIDSegments{
				Name: "iso",
				Segments: []OIDSegment{
					{
						Number: "1",
					},
				},
			},
		},
		{
			name:  "label alias",
			label: "howdy",
			mib:   "{ doody }",
			want: &OIDSegments{
				Name: "howdy",
				Segments: []OIDSegment{
					{
						Name: "doody",
					},
				},
			},
		},
		{
			name:    "no left bracket",
			label:   "error",
			mib:     "error",
			wantErr: true,
		},
		{
			name:    "no terminating right bracket",
			label:   "error",
			mib:     "{error",
			wantErr: true,
		},
		{
			name:  "no number in parens still parses",
			label: "internet",
			mib:   "{iso org() 1}",
			want: &OIDSegments{
				Name: "internet",
				Segments: []OIDSegment{
					{
						Name: "iso",
					},
					{
						Name: "org",
					},
					{
						Number: "1",
					},
				},
			},
		},
		{
			name:  "no right parens still parses",
			label: "internet",
			mib:   "{iso org(3 1}",
			want: &OIDSegments{
				Name: "internet",
				Segments: []OIDSegment{
					{
						Name: "iso",
					},
					{
						Name:   "org",
						Number: "3",
					},
					{
						Number: "1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := tokens.NewLexer(tt.mib)
			got, err := parseObjectID(lexer, tt.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseObjectID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseObjectID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseImports(t *testing.T) {
	tests := []struct {
		name    string
		mib     string
		want    []Import
		wantErr bool
	}{
		{
			name: "three modules with few non-keywords",
			mib: `	MODULE-IDENTITY, OBJECT-TYPE, Integer32, Unsigned32,
					Gauge32, Counter32, Counter64, IpAddress, mib-2
													FROM SNMPv2-SMI
					MODULE-COMPLIANCE, OBJECT-GROUP    FROM SNMPv2-CONF
					InetAddress, InetAddressType,
					InetPortNumber                     FROM INET-ADDRESS-MIB;
		`,
			want: []Import{
				{
					From:  "SNMPv2-SMI",
					Types: []string{"mib-2"},
				},
				{
					From:  "INET-ADDRESS-MIB",
					Types: []string{"InetAddress", "InetAddressType", "InetPortNumber"},
				},
			},
		},
		{
			name:    "no from label is an error",
			mib:     `mib-2 FROM`,
			wantErr: true,
		},
		{
			name:    "no terminating semicolon should return error",
			mib:     `no semicolon`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := tokens.NewLexer(tt.mib)
			got, err := parseImports(lexer)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseImports() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseImports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseModuleIdentity(t *testing.T) {
	tests := []struct {
		name         string
		label        string
		mib          string
		wantSegments *OIDSegments
		wantErr      bool
	}{
		{
			name:    "last updated is required",
			label:   "error",
			mib:     "error",
			wantErr: true,
		},
		{
			name:  "revision must have quoted string",
			label: "error",
			mib: `
			LAST-UPDATED "199505241811Z"
			ORGANIZATION "IETF SNMPv2 Working Group"
			CONTACT-INFO "Chris Goller"
			DESCRIPTION "Error"
			REVISION`,
			wantErr: true,
		},
		{
			name:  "module identity must have equals",
			label: "error",
			mib: `
			LAST-UPDATED "199505241811Z"
			ORGANIZATION "IETF SNMPv2 Working Group"
			CONTACT-INFO "Chris Goller"
			DESCRIPTION "Error"`,
			wantErr: true,
		},
		{
			name:  "module identity ignores extra labels",
			label: "test",
			mib: `
			LAST-UPDATED "199505241811Z"
			ORGANIZATION "IETF SNMPv2 Working Group"
			CONTACT-INFO "Chris Goller"
			DESCRIPTION "Error"
			HOWDY "DOODY"
			::= { experimental xx }
			`,
			wantSegments: &OIDSegments{
				Name: "test",
				Segments: []OIDSegment{
					{
						Name: "experimental",
					},
					{
						Name: "xx",
					},
				},
			},
		},
		{
			name:  "rfc2578",
			label: "fizbin",
			mib: `
			LAST-UPDATED "199505241811Z"
			ORGANIZATION "IETF SNMPv2 Working Group"
			CONTACT-INFO
					"        Marshall T. Rose
	 
					 Postal: Dover Beach Consulting, Inc.
							 420 Whisman Court
							 Mountain View, CA  94043-2186
							 US
	 
						Tel: +1 415 968 1052
						Fax: +1 415 968 2510
	 
					 E-mail: mrose@dbc.mtview.ca.us"
	 
			DESCRIPTION
					"The MIB module for entities implementing the xxxx
					protocol."
			REVISION      "9505241811Z"
			DESCRIPTION
					"The latest version of this MIB module."
			REVISION      "9210070433Z"
			DESCRIPTION
					"The initial version of this MIB module, published in
					RFC yyyy."
		-- contact IANA for actual number
			::= { experimental xx }
`,
			wantSegments: &OIDSegments{
				Name: "fizbin",
				Segments: []OIDSegment{
					{
						Name: "experimental",
					},
					{
						Name: "xx",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := tokens.NewLexer(tt.mib)
			gotSegments, err := parseModuleIdentity(lexer, tt.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseModuleIdentity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSegments, tt.wantSegments) {
				t.Errorf("parseModuleIdentity() = %v, want %v", gotSegments, tt.wantSegments)
			}
		})
	}
}

func Test_parseObject(t *testing.T) {
	tests := []struct {
		name    string
		label   string
		mib     string
		want    *OIDSegments
		wantErr bool
	}{
		{
			name:  "mib definition from rfc1156",
			label: "mib",
			mib:   "IDENTIFIER ::= { mgmt 1 }",
			want: &OIDSegments{
				Name: "mib",
				Segments: []OIDSegment{
					{
						Name: "mgmt",
					},
					{
						Number: "1",
					},
				},
			},
		},
		{
			name:    "identifier is required",
			mib:     "HOWDY",
			wantErr: true,
		},
		{
			name:    "equals is required",
			mib:     "IDENTIFIER",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := tokens.NewLexer(tt.mib)
			got, err := parseObject(lexer, tt.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseObjectGroup(t *testing.T) {
	tests := []struct {
		name    string
		label   string
		mib     string
		want    *OIDSegments
		wantErr bool
	}{
		{
			name:  "snmpGroup definition from rfc2580",
			label: "snmpGroup",
			mib: `
       OBJECTS { snmpInPkts,
                 snmpInBadVersions,
                 snmpInASNParseErrs,
                 snmpBadOperations,
                 snmpSilentDrops,
                 snmpProxyDrops,
                 snmpEnableAuthenTraps }
       STATUS  current
       DESCRIPTION
               "A collection of objects providing basic instrumentation
               and control of an agent."
       REFERENCE "reference"
      ::= { snmpMIBGroups 8 }
`,
			want: &OIDSegments{
				Name: "snmpGroup",
				Segments: []OIDSegment{
					{
						Name: "snmpMIBGroups",
					},
					{
						Number: "8",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := tokens.NewLexer(tt.mib)
			got, err := parseObjectGroup(lexer, tt.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseObjectGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseObjectGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseNotifications(t *testing.T) {
	tests := []struct {
		name    string
		label   string
		mib     string
		want    *OIDSegments
		wantErr bool
	}{
		{
			name:  "snmpBasicNotificationsGroup definition from rfc2580",
			label: "snmpBasicNotificationsGroup",
			mib: `
			NOTIFICATIONS { coldStart, authenticationFailure }
			STATUS        current
			DESCRIPTION
					"The two notifications which an agent is required to
					implement."
		   ::= { snmpMIBGroups 7 }
`,
			want: &OIDSegments{
				Name: "snmpBasicNotificationsGroup",
				Segments: []OIDSegment{
					{
						Name: "snmpMIBGroups",
					},
					{
						Number: "7",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := tokens.NewLexer(tt.mib)
			got, err := parseNotifications(lexer, tt.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseObjectGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseObjectGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseVarBinds(t *testing.T) {
	tests := []struct {
		name    string
		mib     string
		want    []string
		wantErr bool
	}{
		{
			name: "3 variables",
			mib:  "{ l1,  l2, l3}",
			want: []string{"l1", "l2", "l3"},
		},
		{
			name: "variable with syntax; I have never seen this in MIBs",
			mib:  "{GAUGE32}",
			want: []string{"GAUGE32"},
		},
		{
			name:    "no left bracket is an error",
			mib:     "label",
			wantErr: true,
		},
		{
			name:    "no terminating right bracket is an error",
			mib:     "{label",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := tokens.NewLexer(tt.mib)
			got, err := parseVarBinds(lexer)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseVarBinds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseVarBinds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTrapDefinition(t *testing.T) {
	tests := []struct {
		name    string
		label   string
		mib     string
		want    *Trap
		wantErr bool
	}{
		{
			name:  "rfc1215 usage example",
			label: "myLinkDown",
			mib: `
			ENTERPRISE  myEnterprise
			VARIABLES   { ifIndex }
			DESCRIPTION
						"A myLinkDown trap signifies that the sending
						SNMP application entity recognizes a failure
						in one of the communications links represented
						in the agent's configuration."
			::= 2`,
			want: &Trap{
				OIDSegments: OIDSegments{
					Name: "myLinkDown",
					Segments: []OIDSegment{
						{
							Name: "myEnterprise",
						},
						{
							Number: "2",
						},
					},
				},
				VarBinds: []string{"ifIndex"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := tokens.NewLexer(tt.mib)
			got, err := parseTrapDefinition(lexer, tt.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTrapDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTrapDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}
