package abstractor.core.json.parser;

public class Token {
    public final TokenType token;
    public final String value;
    public final Location loc;

    public Token(TokenType token, String value, Location loc) {
        this.token = token;
        this.value = value;
        this.loc = loc;
    }

    @Override
    public String toString() {
        return this.token + " " + this.value + " " + this.loc;
    }
}
