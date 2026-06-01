package testData.java.test1005;

import java.util.*;

public class Foo<T> {
    static Foo<Integer> FooInt() {
        return new Foo<Integer>("int");
    }

    final public String name;
    final public Map<String, T> mapping = new TreeMap<>();

    public Foo(String name) { this.name = name; }
    
    public void add(String key, T value) {
        this.mapping.put(key, value);
    }

    public <S> S lookup(String key) {
        final T value = this.mapping.get(key);
        if (value instanceof S) {
            @SuppressWarnings("unchecked")
            return (S)value;
        }
        return null;
    }

    public int lookupInt(String key) {
        return this.<Integer>lookup(key);
    }
    
    public String lookupString(String key) {
        return this.<String>lookup(key);
    }
    
    public <S extends T> S lookupT(String key) {
        return this.<S>lookup(key);
    }
    
    public int count() {
        return this.mapping.size();
    }
}
