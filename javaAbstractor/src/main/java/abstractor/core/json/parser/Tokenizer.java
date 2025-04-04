package abstractor.core.json.parser;

import java.util.Iterator;

import abstractor.core.iter.PushBackIterator;

public class Tokenizer implements Iterator<Token> {
    private final PushBackIterator<Char> src;

    public Tokenizer(Iterator<Char> src) {
        this.src = new PushBackIterator<Char>(src);
    }

    public boolean hasNext() { return this.src.hasNext(); }

    public Token next() {
        final Char start = this.src.next();
        if (isWhiteSpace(start))  return this.readWhitespace(start);
        if (start.value() == '#') return this.readComment(start);
        if (start.value() == '{') return new Token(TokenType.openCurl,  "{", start.loc());
        if (start.value() == '}') return new Token(TokenType.closeCurl, "}", start.loc());
        if (start.value() == '[') return new Token(TokenType.openSqr,   "[", start.loc());
        if (start.value() == ']') return new Token(TokenType.closeSqr,  "]", start.loc());
        if (start.value() == ':') return new Token(TokenType.colon,     ":", start.loc());
        if (start.value() == ',') return new Token(TokenType.comma,     ",", start.loc());
        if (start.value() == '-') return this.readNumber(start);
        if (isDigit(start))       return this.readNumber(start);
        if (start.value() == '"') return this.readQuote(start);
        if (start.value() == '$') return this.readIdent(start);
        if (isLetter(start))      return this.readIdent(start);
        return new Token(TokenType.error, "unexpected character (" + start.value() + ")", start.loc());
    }

    static private boolean isWhiteSpace(Char c) {
        final char v = c.value();
        return v == '\n' || v == '\r' || v == ' ' || v == '\t';
    }

    static private boolean isDigit(Char c) {
        final char v = c.value();
        return v >= '0' && v <= '9';
    }

    static private boolean isLetter(Char c) {
        final char v = c.value();
        return (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || v == '_';
    }

    static private boolean isSimpleSep(Char c) {
        final char v = c.value();
        return v == '.' || v == '-' || v == '_';
    }

    private Token readWhitespace(Char start) {
        StringBuilder buf = new StringBuilder();
        buf.append(start.value());
        while (this.src.hasNext()) {
            final Char next = this.src.next();
            if (!isWhiteSpace(next)) {
                this.src.pushBack(next);
                break;
            }
            buf.append(next.value());
        }
        return new Token(TokenType.whitespace, buf.toString(), start.loc());
    }

    private Token readComment(Char start) {
        StringBuilder buf = new StringBuilder();
        buf.append(start.value());
        while (this.src.hasNext()) {
            final Char next = this.src.next();
            buf.append(next.value());
            if (next.value() == '\n') break;
        }
        return new Token(TokenType.comment, buf.toString(), start.loc());
    }

    private Token readNumber(Char start) {
        StringBuilder buf = new StringBuilder();
        buf.append(start.value());
        boolean needsDigit = start.value() == '-';
        boolean optionalSign = false;
        boolean hasDec = false;
        boolean hasExp= false;
        while (this.src.hasNext()) {
            final Char next = this.src.next();
            if (isDigit(next)) {
                buf.append(next.value());
                needsDigit = false;
                optionalSign = false;
                continue;
            }
            if (optionalSign) {
                if (next.value() == '+' || next.value() == '-') {
                    buf.append(next.value());
                    optionalSign = false;
                    continue;
                }
            }
            if (needsDigit) break;
            if (next.value() == '.') {
                if (!hasDec && !hasExp) {
                    hasDec = true;
                    needsDigit = true;
                    buf.append(next.value());
                    continue;
                }
            }
            if (next.value() == 'e' || next.value() == 'E') {
                if (!hasExp) {
                    hasExp = true;
                    needsDigit = true;
                    optionalSign = true;
                    buf.append(next.value());
                    continue;
                }
            }
            this.src.pushBack(next);
            break;
        }
        if (needsDigit)
            return new Token(TokenType.error, "unexpected end of number", start.loc());
        TokenType type = hasDec || hasExp ? TokenType.real : TokenType.integer;
        return new Token(type, buf.toString(), start.loc());
    }

    private int readHexDigit() throws JsonException {
        if (this.src.hasNext()) {
            final Char c = this.src.next();
            if (isDigit(c)) return (int)(c.value() - '0');
            if (c.value() >= 'a' && c.value() <= 'f') return (int)(c.value() - 'a') + 10;
            if (c.value() >= 'A' && c.value() <= 'F') return (int)(c.value() - 'A') + 10;
        }
        throw new JsonException("expected a hex value.");
    }

    private Token readQuote(Char start) {
        StringBuilder buf = new StringBuilder();
        while (this.src.hasNext()) {
            final Char next = this.src.next();
            char c = next.value();
            if (next.value() == '"') {
                return new Token(TokenType.quote, buf.toString(), start.loc());
            }
            if (next.value() == '\\') {
                if (!this.src.hasNext()) break;
                final Char escape = this.src.next(); 
                switch (escape.value()) {
                    case '"':  c = '\"'; break;
                    case '\\': c = '\\'; break;
                    case '/':  c = '/';  break;
                    case 'b':  c = '\b'; break;
                    case 'f':  c = '\f'; break;
                    case 'n':  c = '\n'; break;
                    case 'r':  c = '\r'; break;
                    case 't':  c = '\t'; break;
                    case 'u':
                        try {
                            final int v1 = this.readHexDigit();
                            final int v2 = this.readHexDigit();
                            final int v3 = this.readHexDigit();
                            final int v4 = this.readHexDigit();
                            c = (char)(v1 << 12 | v2 << 8 | v3 << 4 | v4);
                            break;
                        } catch (Exception ex) {
                            return new Token(TokenType.error, "unexpected end of escaped unicode", escape.loc());
                        }
                    default:
                        return new Token(TokenType.error, "unexpected escaped value: " + escape.value(), escape.loc());
                }
            }
            buf.append(c);
        }
        return new Token(TokenType.error, "unexpected end of quoted string", start.loc());
    }

    private Token readIdent(Char start) {
        StringBuilder buf = new StringBuilder();
        buf.append(start.value());
        while (this.src.hasNext()) {
            final Char next = this.src.next();
            if (isDigit(next) || isLetter(next) || isSimpleSep(next)) { 
                buf.append(next.value());
                continue;
            }
            this.src.pushBack(next);
            break;
        }
        return new Token(TokenType.ident, buf.toString(), start.loc());
    }
}
