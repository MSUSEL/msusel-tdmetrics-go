package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.TreeSet;

import spoon.reflect.declaration.CtElement;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;
import abstractor.core.log.*;

public class Factory<T extends Construct> implements Jsonable {
    static private final boolean logCreate = true;
    
    private final ConstructKind              conKind;
    public  final HashMap<CtElement, Ref<T>> byElem     = new HashMap<CtElement, Ref<T>>();
    public  final HashMap<T,         Ref<T>> nonElemRef = new HashMap<T, Ref<T>>();
    public  final TreeSet<Ref<T>>            refSet     = new TreeSet<Ref<T>>();
    public  final TreeSet<T>                 conSet     = new TreeSet<T>();

    public Factory(ConstructKind kind) { this.conKind = kind; }

    public ConstructKind kind() { return this.conKind; }

    public String toString() { return "factory " + this.conKind; }

    public int refSize() { return this.refSet.size(); }
    public int size()    { return this.conSet.size(); }

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
            if (!newCon.kind().equals(this.conKind))
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

    public void removeElem(Logger log, CtElement elem, String title) {
        if (logCreate) log.log("Removing " + title);

        final Ref<T> ref = this.getRef(elem);
        this.byElem.remove(elem);

        if (ref == null) return;
        this.refSet.remove(ref);

        if (!ref.isResolved()) return;
        final T con = ref.getResolved();
        this.conSet.remove(con);
        this.nonElemRef.remove(con);
    }

    private List<Ref<T>> findRefsForCon(T con) {
        return this.refSet.stream().filter(r -> con.equals(r.getResolved())).toList();
    }

    public void setRefForElem(CtElement elem, Ref<T> ref) throws Exception {
        final Ref<T> existing = this.getRef(elem);
        if (existing != null) {
            if (existing.equals(ref)) return;
            throw new Exception("Ref already exists for element " + existing + " so cannot set " + ref + ".");
        }
        this.refSet.add(ref);
        this.byElem.put(elem, ref);
    }

    public Ref<T> addOfGetRefForElem(CtElement elem, String title) throws Exception {
        final Ref<T> existing = this.getRef(elem);
        if (existing != null) return existing;

        final Ref<T> newRef = new Ref<T>(this.conKind, elem, title);
        this.refSet.add(newRef);
        this.byElem.put(elem, newRef);
        return newRef;
    }

    public Ref<T> addOrGetRef(T c) throws Exception {
        final T other = this.getExisting(c);
        if (other != null) c = other;

        Ref<T> ref = this.nonElemRef.get(c);
        if (ref != null) return ref;

        ref = new Ref<T>(this.conKind, null, "no element ref");
        ref.setResolved(c);

        this.refSet.add(ref);
        this.conSet.add(c);
        this.nonElemRef.put(c, ref);
        return ref;
    }

    public boolean consolidateCons(Logger log) throws Exception {
        // Copy all cons to a list and clear the set.
        final ArrayList<T> conList = new ArrayList<T>(this.conSet);
        this.conSet.clear();

        CmpOptions options = new CmpOptions();
        options.useResolved = true;
        for (T con : conList) con.setCmpOptions(options);
        for (Ref<T> ref : this.refSet) ref.setCmpOptions(options);

        boolean collision = false;
        for (T con : conList) {
            T existing = this.conSet.floor(con);
            if (existing == null || !existing.equals(con)) {
                this.conSet.add(con);
                continue;
            }

            // Found another construct that is equal so move all references over
            // to the existing construct since the duplicate is about to be removed.
            collision = true;
            List<Ref<T>> refs = this.findRefsForCon(con);
            for (Ref<T> ref : refs) ref.setResolved(existing);
            con.setIndex(-100);

            // TODO: Need to handle any non-references that need
            // to be moved over. There shouldn't be any, but double check.

        }
        return collision;
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
        for (T t : this.conSet) {
            JsonNode node = t.toJson(h);
            if (h.writeRefs) {
                JsonObject obj;
                if (node instanceof JsonObject o) obj = o;
                else {
                    obj = new JsonObject();
                    obj.put("resolved", node);
                }

                List<Ref<T>> refs = this.findRefsForCon(t);
                JsonArray refList = new JsonArray();
                for (Ref<T> r : refs) refList.add(r.refToJson(h));
                obj.put("refs", refList);
                node = obj;
            }
            array.add(node);
        }
        return array;
    }
}
