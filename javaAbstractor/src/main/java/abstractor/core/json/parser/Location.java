package abstractor.core.json.parser;

public record Location(int offset, int lineNo, int column) {
    
    @Override
    public String toString() {
        return "at offset " + this.offset() + ", line " + this.lineNo() + ", column " + this.column();
    }
}
