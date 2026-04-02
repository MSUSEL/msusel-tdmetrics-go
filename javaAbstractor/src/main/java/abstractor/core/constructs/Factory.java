package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.TreeSet;

import spoon.reflect.declaration.CtElement;

import abstractor.core.json.*;
import abstractor.core.log.*;

public class Factory<T extends Construct> implements Jsonable, Iterable<T> {
    static private final boolean logCreate = false;
    
    private final ConstructKind              conKind;
    private final TreeSet<Ref<T>>            refSet     = new TreeSet<Ref<T>>();
    private final TreeSet<T>                 conSet     = new TreeSet<T>();
    private final HashMap<CtElement, Ref<T>> byElem     = new HashMap<CtElement, Ref<T>>();
    private final HashMap<T,         Ref<T>> nonElemRef = new HashMap<T, Ref<T>>();

    public Factory(ConstructKind kind) { this.conKind = kind; }

    public ConstructKind kind() { return this.conKind; }

    public String toString() { return "factory " + this.conKind; }

    public int refSize() { return this.refSet.size(); }
    public int size()    { return this.conSet.size(); }

    public void clear() {
        this.refSet.clear();
        this.conSet.clear();
        this.byElem.clear();
        this.nonElemRef.clear();
    }

    public Iterator<Ref<T>> refIterator() { return this.refSet.iterator(); }
    public Iterator<T>      iterator()    { return this.conSet.iterator(); }

    public List<T> toList() {
        ArrayList<T> list = new ArrayList<>(this.conSet.size());
        for (T value : this.conSet) list.add(value);
        return Collections.unmodifiableList(list);
    }

    public T get(int index) {
        return this.conSet.stream().skip(index).findFirst().orElse(null);
    }

    public Ref<T> getRef(CtElement elem) { return this.byElem.get(elem); }
    
    public T getExisting(T c) {
        final T other = this.conSet.floor(c);
        return c.equals(other) ? other : null;
    }

    @FunctionalInterface
    public interface Creator<T extends Construct> { T create() throws Exception; }
    
    @FunctionalInterface
    public interface Finisher<T extends Construct> { void finish(Ref<T> ref, T con) throws Exception; }

    public Ref<T> create(Logger log, CtElement elem, String title, Creator<T> creator, Finisher<T> finisher) throws Exception {
        if (elem == null) return null;

        // Check if a resistance already exists.
        final Ref<T> existing = this.getRef(elem);
        if (existing != null) return existing;
        
        try {
            if (logCreate) {
                log.log("Adding " + title);
                log.push();
            }

            // First add a reference so that if a circular loop is hit when
            // creating the new construct, the same reference will be picked up.
            final Ref<T> newRef = new Ref<T>(this.conKind, elem, title);
            this.refSet.add(newRef);
            this.byElem.put(elem, newRef);

            // Create a new construct for this data.
            final T newCon = creator.create();
            if (newCon == null)
                throw new Exception("Factory creator for " + this.toString() + " returned null.");
            if (newCon.kind() != this.conKind)
                throw new Exception("Factory creator for " + this.toString() + " create a type with kind " + newCon.kind() + ".");

            // If an existing construct matches the new one after the new one
            // has been loaded, then there are two elements to get to the same
            // value. Resolve the reference for the existing or new construct.
            final T other = this.getExisting(newCon);
            if (other != null) {
                newRef.setResolved(other);
            } else {
                this.conSet.add(newCon);
                newRef.setResolved(newCon);
                if (finisher != null) finisher.finish(newRef, newCon);
            }

            return newRef;
        } finally {
            if (logCreate) log.pop();
        }
    }

    public Ref<T> create(Logger log, CtElement elem, String title, Creator<T> creator) throws Exception {
        return this.create(log, elem, title, creator, null);
    }

    public void setRefForElem(CtElement elem, Ref<T> ref) throws Exception {
        final Ref<T> existing = this.getRef(elem);
        if (existing != null) {
            if (existing == ref) return;
            throw new Exception("Ref already exists for element " + existing + " so cannot set " + ref);
        }
        this.refSet.add(ref);
        this.byElem.put(elem, ref);
    }

    public Ref<T> addOfGetRefForElem(CtElement elem, String title) {
        final Ref<T> existing = this.getRef(elem);
        if (existing != null) return existing;

        final Ref<T> newRef = new Ref<T>(this.conKind, elem, title);
        this.refSet.add(newRef);
        this.byElem.put(elem, newRef);
        return newRef;
    }

    public Ref<T> addOrGetRef(T c) {
        final T other = this.getExisting(c);
        if (other != null) c = other;

        Ref<T> ref = this.nonElemRef.get(c);
        if (ref != null) return ref;

        ref = new Ref<T>(this.conKind, null, "no element ref");
        this.conSet.add(c);
        this.nonElemRef.put(c, ref);
        return ref;
    }

    public void setIndices() {
        int index = 1;
        for (Construct o : this.conSet) {
            o.setIndex(index);
            index++;
        }
    }

    public JsonNode toJson(JsonHelper h) {
        JsonArray array = new JsonArray();
        for (T t : this.conSet) array.add(t.toJson(h));
        return array;
    }
}
