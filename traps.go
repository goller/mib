package mib

// Trap defines a data schema that are sent from a client to a collector.
// TRAP-TYPE is an SNMPv1 macro defined in RFC1215.
// NOTIFICATION-TYPE is an SNMPv2 macro defined in RFC2578.
type Trap struct {
	OIDSegments
	VarBinds []string
}
