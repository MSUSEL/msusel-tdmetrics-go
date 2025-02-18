package abstractor.core.json.parser;

public enum TokenType {
    whitespace, // \s+
    comment,    // #[^\n]*\n
    openCurl,   // {
    closeCurl,  // }
    openSqr,    // [
    closeSqr,   // ]
    colon,      // :
    comma,      // ,
    quote,      // \"[^\"]*\"
    boolId,     // true|false
    nullId,     // null
    integer,    // -?[0-9]+
    real,       // -?[0-9]+(\.[0-9]+)(e/E)
    ident,      // \S+
    error       // <invalid>
}
