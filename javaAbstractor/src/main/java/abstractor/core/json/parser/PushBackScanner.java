package abstractor.core.json.parser;

import java.util.Iterator;
import java.util.Stack;

public class PushBackScanner implements Iterator<Char> {
    private final Iterator<Char> src;
    private final Stack<Char> pending = new Stack<Char>();

    public PushBackScanner(Iterator<Char> src) { this.src = src; }

    public boolean hasNext() {
        return !this.pending.empty() || this.src.hasNext();
    }

    public void pushBack(Char c) { this.pending.push(c); }

    public Char next() {
        return this.pending.empty() ? this.src.next() : this.pending.pop();
    }
}
