package abstractor.core.iter;

import java.util.Iterator;
import java.util.Stack;

public class PushBackIterator<T> implements Iterator<T> {
    private final Iterator<T> src;
    private final Stack<T> pending = new Stack<T>();

    public PushBackIterator(Iterator<T> src) { this.src = src; }

    public boolean hasNext() {
        return !this.pending.empty() || this.src.hasNext();
    }

    public void pushBack(T c) { this.pending.push(c); }

    public T next() {
        return this.pending.empty() ? this.src.next() : this.pending.pop();
    }
}
