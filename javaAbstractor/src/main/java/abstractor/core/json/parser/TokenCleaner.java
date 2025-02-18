package abstractor.core.json.parser;

import java.util.Iterator;

public class TokenCleaner implements Iterator<Token> {
    private final Iterator<Token> src;
    private Token next;

    public TokenCleaner(Iterator<Token> src) {
        this.src = src;
        this.next = null;
    }

    private Token readNext() {
        final Token next = this.src.next();
        if (next.token == TokenType.ident) {
            if (next.value.equals("null"))  return new Token(TokenType.nullId, next.value, next.loc);
            if (next.value.equals("true"))  return new Token(TokenType.boolId, next.value, next.loc);
            if (next.value.equals("false")) return new Token(TokenType.boolId, next.value, next.loc);
            return next;
        }
        if (next.token == TokenType.whitespace) return null;
        if (next.token == TokenType.comment) return null;
        return next;
    }

    private void seekNext() {
        if (this.next != null) return;
        while (this.src.hasNext()) {
            this.next = this.readNext();
            if (this.next != null) return;
        }
    }

    public boolean hasNext() {
        this.seekNext();
        return this.next != null;
    }

    public Token next() {
        this.seekNext();
        final Token cur = this.next;
        this.next = null;
        return cur;
    }
}
