package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.TreeSet;

import spoon.reflect.declaration.CtElement;

import abstractor.core.json.*;

public class Factory<T extends Construct> implements Jsonable, Iterable<T> {
    private final ConstructKind conKind;
    private final TreeSet<T> set;
    private final Map<CtElement, T> byElem;
    public boolean inProgress;

    public Factory(ConstructKind kind) {
        this.conKind = kind;
        this.set = new TreeSet<T>();
        this.byElem = new HashMap<CtElement, T>();
        this.inProgress = false;
    }

    public ConstructKind kind() { return this.conKind; }

    public int size() { return this.set.size(); }

    public void clear() {
        this.set.clear();
        this.byElem.clear();
    }

    public Iterator<T> iterator() { return this.set.iterator(); }

    public List<T> toList() {
        ArrayList<T> list = new ArrayList<>(this.set.size());
        for (T value : this.set) list.add(value);
        return Collections.unmodifiableList(list);
    }

    public T get(int index) {
        int i = 0;
        for (T value : this.set) {
            if (i == index) return value;
            i++;
        }
        return null;
    }

    public T get(CtElement elem) {
        return this.byElem.get(elem);
    }
    
    public T get(T c) {
        final T other = this.set.floor(c);
        if (c.equals(other)) return other;
        return null;
    }

    public T addOrGet(CtElement elem, T c) {
        final T other = this.get(c);
        if (other != null) {
            this.addElemKey(elem, other);
            return other;
        }
        this.add(elem, c);
        return c;
    }

    public T addOrGet(T c) {
        return this.addOrGet(null, c);
    }

    public void addElemKey(CtElement elem, T c) {
        if (elem != null) this.byElem.put(elem, c);
    }
    
    public void add(CtElement elem, T c) {
        this.set.add(c);
        this.addElemKey(elem, c);
    }

    public void setIndices() {
        int index = 1;
        for (Construct o : this.set) {
            o.setIndex(index);
            index++;
        }
    }

    public JsonNode toJson(JsonHelper h) {
        JsonArray array = new JsonArray();
        for (T t : this.set) array.add(t.toJson(h));
        return array;
    }
}
