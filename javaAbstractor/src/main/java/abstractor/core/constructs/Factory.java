package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.TreeSet;

import spoon.reflect.declaration.CtElement;

import abstractor.core.json.*;
import abstractor.core.log.Logger;

public class Factory<T extends Construct> implements Jsonable {
    private final ConstructKind conKind;
    private final TreeSet<T> set;
    private final HashMap<CtElement, T> byElem;

    public Factory(ConstructKind kind) {
        this.conKind = kind;
        this.set = new TreeSet<T>();
        this.byElem = new HashMap<CtElement, T>();
    }

    public ConstructKind kind() { return this.conKind; }

    public int size() { return this.set.size(); }

    public T get(int index) {
        int i = 0;
        for (T value : this.set) {
            if (i == index) return value;
            i++;
        }
        return null;
    }

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

    public interface ConstructCreator<T extends Construct> { T create(); }

    public interface FinishConstruct<T extends Construct> { void finish(T con); }

    public T create(Logger log, CtElement elem, String title,
        ConstructCreator<T> c, FinishConstruct<T> f) {
        final T existing = this.byElem.get(elem);
        if (existing != null) return existing;
        
        if (log != null) {
            log.log("Adding " + title);
            log.push();
        }
        try {
            final T newCon = c.create();

            T other = this.set.floor(newCon);
            if (newCon.equals(other)) return other;
            this.set.add(newCon);
            this.byElem.put(elem, newCon);

            if (f != null) f.finish(newCon);
            return newCon;
        } finally {
            if (log != null) log.pop();
        }
    }
    
    public T create(Logger log, CtElement elem, String title,
        ConstructCreator<T> c) {
        return this.create(log, elem, title, c, null);
    }

    public T addOrGet(T c) {
        T other = this.set.floor(c);
        if (c.equals(other)) return other;
        this.set.add(c);
        return c;
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
