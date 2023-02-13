package util;

import java.util.*;

public class NamedSet<T extends NamedObject> implements Iterable<T> {
    private final SortedMap<String, T> data;

    public NamedSet() {
        this.data = new TreeMap<>();
    }

    public boolean isEmpty() {
        return this.data.isEmpty();
    }

    public boolean contains(final T named) {
        if (named == null) return false;
        return this.data.containsKey(named.name());
    }

    public boolean contains(final String name) {
        return this.data.containsKey(name);
    }

    public Iterator<T> iterator() {
        return this.data.values().iterator();
    }

    public Iterator<String> names() {
        return this.data.keySet().iterator();
    }

    public T get(final String name) {
        return this.data.get(name);
    }

    public boolean add(final T named) {
        if (named == null) return false;
        final String name = named.name();
        if (this.data.containsKey(name)) return false;
        this.data.put(name, named);
        return true;
    }

    public boolean remove(final T named) {
        if (named == null) return false;
        final String name = named.name();
        T other = this.data.get(name);
        if (named != other) return false;
        return this.data.containsKey(name);
    }

    public boolean remove(final String name) {
        return this.data.containsKey(name);
    }
}
