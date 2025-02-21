package abstractor.core.iter;

import java.util.Iterator;

public class Iter {
    private Iter() { }

    public static <T> Iterable<T> Single(T value) {
        return new IterableWrapper<T>(() -> SingleIterator(value) );
    }

    public static <T> Iterator<T> SingleIterator(T value) {
        class SingleIter implements Iterator<T> {
            private final T value;
            private boolean read;
    
            public SingleIter(T value) { this.value = value; }
    
            @Override
            public boolean hasNext() { return !this.read; }
        
            @Override
            public T next() {
                this.read = true;
                return this.value;
            }
        }
        return new SingleIter(value);
    }
    
    public static <T> Iterable<T> Empty() {
        return new IterableWrapper<T>(() -> EmptyIterator() );
    }

    public static <T> Iterator<T> EmptyIterator() {
        class EmptyIter implements Iterator<T> {
            @Override
            public boolean hasNext() { return false; }
        
            @Override
            public T next() { return null; }
        }
        return new EmptyIter();
    }
    
    public static <T> Iterable<T> Array(T[] values) {
        return new IterableWrapper<T>(() -> ArrayIterator(values) );
    }

    public static <T> Iterator<T> ArrayIterator(T[] values) {
        class ArrayIter implements Iterator<T> {
            private final T[] values;
            private int index;
    
            public ArrayIter(T[] values) { this.values = values; }
    
            @Override
            public boolean hasNext() { return this.index < this.values.length; }
    
            @Override
            public T next() {
                final T result = this.values[this.index];
                this.index++;
                return result;
            }
        }
        return new ArrayIter(values);
    }
}
