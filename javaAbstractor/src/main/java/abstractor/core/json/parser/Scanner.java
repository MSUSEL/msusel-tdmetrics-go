package abstractor.core.json.parser;

import java.util.Iterator;

public class Scanner implements Iterator<Char> {
    private final String src;
    private int offset;
    private int lineNo;
    private int column;

    public Scanner(String src) {
        this.src = src;
        this.offset = 0;
        this.lineNo = 1;
        this.column = 1;
    }

    public boolean hasNext() {
        return this.offset < this.src.length();
    }

    public Char next() {
        char cur = this.src.charAt(this.offset);
        this.offset++;
        this.column++;
        if (cur == '\n') {
            this.lineNo++;
            this.column = 1;
        }
        return new Char(cur, new Location(this.offset, this.lineNo, this.column));
    }
}
