package abstractor.core.json.parser;

public class Char {
    public final char value;
    public final Location loc;

    public Char(char value, Location loc) {
        this.value = value;
        this.loc = loc;
    }

    @Override
    public String toString() {
        return "char "+this.value + " " + this.loc;
    }
}
