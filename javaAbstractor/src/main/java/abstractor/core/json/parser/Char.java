package abstractor.core.json.parser;

public record Char(char value, Location loc) {

    @Override
    public String toString() {
        return "char "+this.value() + " " + this.loc();
    }
}
