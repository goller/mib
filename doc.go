package mib

/*
type ColumnInfo struct {
        OIDInfo
        // If true this column is an index.  This is used to determine if it should
        // be a tag in the generated metric.
        Index bool
}

type TableInfo struct {
        // Numeric OID to ColumnInfo
        Columns map[string]ColumnInfo
}

type OIDInfo struct {
        // Name of the MIB; ex: "IF-MIB"
        // I don't see this being used... here for parity
        MIBName string
        // Textual representation of the final path of the OID; ex: "ifDescr"
        // Used for naming in the generated metrics.
        OIDText string
        // Numerical representation of the full OID; ex: ".1.3.6.1.2.1.2.2.1.2"
        // This is used for identification.
        OID string
        // Textual Convention; ex: PhysAddress
        // Maybe an enum int of known types.  We use this for some mac and ip
        // address conversion, though I think it might not work right.
        TextualConvention string
}

type MIB interface {
        // LookupTable looks up the information for a conceptual table.
        // oid="IF-MIB::ifTable"
        // oid=".1.3.6.1.2.1.2.2"
        LookupTableInfo(oid string) (TableInfo, error)

        // LookupOID looks up an OID in the MIB.  The oid can be either the textual
        // or the numeric representation of an oid.
        //
        // oid="IF-MIB::ifTable"
        // oid=".1.3.6.1.2.1.2.2"
        LookupOIDInfo(oid string) (OIDInfo, error)
}


*/
/*
Tables matters because we walk it because we want all the data in the "bundle"
We get the names for each column so we can name them correctly
If a column is indexed then it is converted to a tag

Need to know the OIDs of the values because
*/
