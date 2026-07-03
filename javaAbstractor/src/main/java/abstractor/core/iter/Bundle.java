package abstractor.core.iter;

import java.util.*;

public class Bundle<T> implements Iterator<T> {
    final private LinkedList<Iterator<? extends T>> src = new LinkedList<>();
    private Iterator<? extends T> current;
    private boolean hasNextValue;
    private T nextValue;

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
    
    private void seekNext() {
        if (this.hasNextValue) return;
        while (true) {
            if (this.current != null) {
                if (this.current.hasNext()) {
                    this.nextValue = this.current.next();
                    this.hasNextValue = true;
                    return;
                }
                this.current = null;
            }
            if (!this.src.isEmpty()) {
                this.current = this.src.pollFirst();
                continue;
            }
            return;
        }
    }

    @Override
    public boolean hasNext() {
        this.seekNext();
        return this.hasNextValue;
    }

    @Override
    public T next() {
        this.seekNext();
        T result = this.nextValue;
        this.nextValue = null;
        this.hasNextValue = false;
        return result;
    }
}
