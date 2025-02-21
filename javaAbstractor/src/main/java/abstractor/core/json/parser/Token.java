package abstractor.core.json.parser;

public record Token(TokenType token, String value, Location loc) {

    @Override
    public String toString() {
        return this.token() + " " + this.value() + " " + this.loc();
    }
}
