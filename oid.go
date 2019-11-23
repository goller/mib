package mib

import "fmt"

// OID is an Object Identifier naming any object with a globally unambiguous
// persistent name.
type OID struct {
	MIB    string // Filename of the OID
	Name   string
	Number string
	Parent *OID // I need to experiment with this... I don't like it
}

type OIDSegments struct {
	Name     string
	Segments []OIDSegment
}

type OIDSegment struct {
	Name   string
	Number string
}

func (o OIDSegment) String() string {
	return fmt.Sprintf("Name: '%s' Number '%s'", o.Name, o.Number)
}
