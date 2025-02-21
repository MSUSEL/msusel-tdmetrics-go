package abstractor.core.json.parser;

import java.util.Iterator;
import java.util.Map;

public class TokenCleaner implements Iterator<Token> {
    private final Iterator<Token> src;
    private Token next;
    private Map<String, TokenType> remap;

    public TokenCleaner(Iterator<Token> src) {
        this.src = src;
        this.next = null;
        this.remap = Map.of(
            "null",  TokenType.nullId,
            "true",  TokenType.boolId,
            "false", TokenType.boolId
        );
    }

    private Token readNext() {
        final Token next = this.src.next();
        if (next.token() == TokenType.ident) {
            TokenType override = this.remap.get(next.value());
            if (override != null) return new Token(override, next.value(), next.loc());
            return next;
        }
        if (next.token() == TokenType.whitespace) return null;
        if (next.token() == TokenType.comment) return null;
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
