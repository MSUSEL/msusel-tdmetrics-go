package abstractor.core.iter;

import java.util.*;

public class Bundle<T> implements Iterator<T> {
    final private List<Iterator<? extends T>> src = new LinkedList<>();
    final private ExpandIterator<T> expander = new ExpandIterator<T>(src.iterator());
    
    public Bundle() { }
    public Bundle(Iterator<? extends T> values) { this.add(values); }
    public Bundle(Iterable<? extends T> values) { this.add(values); }
    public Bundle(T[] values)                   { this.add(values); }
    public Bundle(T value)                      { this.add(value); }
    
    public Bundle<T> add(Iterator<? extends T> values) {
        if (values != null) this.src.add(values);
        return this;
    }

    public Bundle<T> add(Iterable<? extends T> values) {
        if (values == null) return this;
        return this.add(values.iterator());
    }

    public Bundle<T> add(T[] values) {
        if (values == null) return this;
        return this.add(Iter.ArrayIterator(values));
    }

    public Bundle<T> add(T value) {
        if (value == null) return this;
        return this.add(Iter.SingleIterator(value));
    }

    @Override public boolean hasNext() { return this.expander.hasNext(); }
    @Override public T next() { return this.expander.next(); }
}
