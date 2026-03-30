package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.TreeSet;
import java.util.function.Supplier;
import java.util.function.Consumer;

import spoon.reflect.declaration.CtElement;

import abstractor.core.json.*;
import abstractor.core.log.*;

public class Factory<T extends Construct> implements Jsonable, Iterable<T> {
    static private final boolean logCreate = false;

    private final ConstructKind conKind;
    private final TreeSet<T>    set = new TreeSet<T>();
    private final Supplier<T>   creator;
    private final HashMap<CtElement, T> byElem = new HashMap<CtElement, T>();

    public Factory(ConstructKind kind, Supplier<T> creator) {
        this.conKind = kind;
        this.creator = creator;
    }

    public ConstructKind kind() { return this.conKind; }

    public String toString() { return this.conKind.toString(); }

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
        return this.set.stream().skip(index).findFirst().orElse(null);
    }

    public T get(CtElement elem) { return this.byElem.get(elem); }
    
    public T get(T c) {
        final T other = this.set.floor(c);
        return c.equals(other) ? other : null;
    }

    public T create(Logger log, CtElement elem, String title, Consumer<T> loader) throws Exception {
        if (elem == null) return null;
        final T existing = this.get(elem);
        if (existing != null) return existing;

        try {
            if (logCreate) {
                log.log("Adding " + title);
                log.push();
            }

            final T newCon = this.creator.get();
            if (newCon == null)
                throw new Exception("Factory creator for " + this.toString() + " returned null.");
            this.add(elem, newCon);

            loader.accept(newCon);

            return newCon;
        } finally {
            if (logCreate) log.pop();
        }
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
