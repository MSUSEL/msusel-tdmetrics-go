package abstractor.core.iter;

import java.util.*;

public class ExpandIterator<T> implements Iterator<T> {
    private final Iterator<? extends Iterator<? extends T>> src;
    private Iterator<? extends T> current;
    private boolean hasNextValue;
    private T nextValue;

    public ExpandIterator(Iterable<? extends Iterator<? extends T>> src) {
        this.src = src.iterator();
    }

    public ExpandIterator(Iterator<? extends Iterator<? extends T>> src) {
        this.src = src;
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
            if (this.src.hasNext()) {
                this.current = this.src.next();
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
