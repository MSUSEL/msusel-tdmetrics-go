package abstractor.core.json.parser;

public class Location {
    public final int offset;
    public final int lineNo;
    public final int column;

    public Location(int offset, int lineNo, int column) {
        this.offset = offset;
        this.lineNo = lineNo;
        this.column = column;
    }

    @Override
    public String toString() {
        return "at offset " + this.offset + ", line " + this.lineNo + ", column " + this.column;
    }
}
