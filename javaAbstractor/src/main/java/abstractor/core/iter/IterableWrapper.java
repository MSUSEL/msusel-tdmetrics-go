package abstractor.core.iter;

import java.util.Iterator;

public class IterableWrapper<T> implements Iterable<T> {
    public interface IterableFn<T> { Iterator<T> iterator(); }
    
    private final IterableFn<T> iter;
    public IterableWrapper(IterableFn<T> iter) { this.iter = iter; }

    @Override
    public Iterator<T> iterator() { return this.iter.iterator(); }
}
