package abstractor.core.constructs;

import java.util.*;

import spoon.reflect.declaration.CtElement;

import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;
import abstractor.core.log.*;
import abstractor.core.require.Require;

public class Factory<T extends Construct> implements Jsonable {
    static private final boolean logCreate = true;
    
    private final ConstructKind              conKind;

    // TODO: HashMap and HashSet still check equal, so if two CtElements have
    // the same hash, that doesn't mean the HashMap/Set will treat them as the same.

    private final HashMap<Integer, CtElement> elemHash   = new HashMap<>();
    private final HashMap<Integer, Ref<T>>    byElem     = new HashMap<>();
    private final HashMap<T,       Ref<T>>    nonElemRef = new HashMap<>();
    private final HashSet<Integer>            elemInProg = new HashSet<>();
    private final TreeSet<Ref<T>>             refSet     = new TreeSet<>();
    private final TreeSet<T>                  conSet     = new TreeSet<>();

    public Factory(ConstructKind kind) { this.conKind = kind; }

    public ConstructKind kind() { return this.conKind; }

    public String toString() { return "factory " + this.conKind; }

    public int refSize() { return this.refSet.size(); }
    public int size()    { return this.conSet.size(); }

    public TreeSet<Ref<T>> getRefSet() { return this.refSet; }
    public TreeSet<T>      getConSet() { return this.conSet; }

    public List<T> toList() {
        ArrayList<T> list = new ArrayList<>(this.conSet.size());
        for (T value : this.conSet) list.add(value);
        return Collections.unmodifiableList(list);
    }

    public T get(int index) {
        return this.conSet.stream().skip(index).findFirst().orElse(null);
    }

    public Ref<T> getRef(CtElement elem) {
        Ref<T> r = this.byElem.get(elem.hashCode());
        if (elem.hashCode() == 3002523) { // TODO: REMOVE
            System.out.println();
            System.out.println(">>(getRef) " + elem.hashCode() + " => " + r); // TODO: REMOVE
            this.printRefDebug(); // TODO: REMOVE
        }
        return r;
    }
    
    public T getExisting(T c) {
        final T other = this.conSet.floor(c);
        return c.equals(other) ? other : null;
    }

    private List<Ref<T>> findRefsForCon(T con) {
        return this.refSet.stream().filter(r -> con.equals(r.getResolved())).toList();
    }

    private void printRefDebug() {
        System.out.println("===[byElem]=========================================");
        int i = 0;
        for (int hash: this.byElem.keySet()) {
            CtElement e2 = this.elemHash.get(hash);
            System.out.println("byElem: " + i + ". " + e2.hashCode() + " => " + this.byElem.get(hash));
            i++;
        }
        System.out.println("===[refSet]=========================================");
        i = 0;
        for (Ref<?> r2: this.refSet) {
            System.out.println("refSet: " + i + ". " + r2);
            i++;
        }
        System.out.println("====================================================");
    }

    //==========================================================================

    private void addRef(Ref<T> ref) throws Exception {   
        Require.require(this.refSet.add(ref),
            "reference " + ref + " must be added at this point");

        final CtElement elem = ref.elem;
        if (elem != null) {
            Require.isNull(this.elemHash.put(elem.hashCode(), elem));
            Require.isNull(this.byElem.put(elem.hashCode(), ref));
        }

        if (ref.isResolved()) {
            final T c = ref.getResolved();
            Require.require(this.conSet.add(c),
                "resolved construct " + c + " must be added at this point");
            Require.isNull(this.nonElemRef.put(c, ref),
                "resolved construct " + c + " and reference " + ref + " must be added at this point for non-element ref");
        }
    }

    @FunctionalInterface
    public interface Creator<T extends Construct> { T create() throws Exception; }
    
    @FunctionalInterface
    public interface Finisher<T extends Construct> { void finish(Ref<T> ref, T con) throws Exception; }

    public Ref<T> create(Logger log, CtElement elem, String title, Creator<T> creator, Finisher<T> finisher) throws Exception {
        if (elem == null) return null;

        // If already "in progress" then check for if a reference already exists
        // so that we only create one and all others are references. However,
        // since references can be created other ways, we need to skip checking
        // for references if not "in progress" to start progress.
        final Ref<T> existing = this.getRef(elem);
        final boolean inProgress = this.elemInProg.contains(elem.hashCode());
        if (inProgress && existing != null) return existing;
        
        try {
            if (logCreate) {
                log.log("Adding " + title);
                log.push();
            }

            // First add a reference so that if a circular loop is hit when
            // creating the new construct, the same reference will be picked up.
            Ref<T> ref;
            if (existing != null) ref = existing;
            else {
                final Ref<T> newRef = new Ref<T>(this.conKind, elem, title);

                // TODO: REMOVE THE FOLLOWING DEBUG PRINTS
                if (elem.hashCode() == 3002523) {
                    final Ref<T> oldRef = this.refSet.floor(newRef);
                    System.out.println();
                    System.out.println("===[create]===");
                    System.out.println("newRef: " + newRef);
                    System.out.println("oldRef: " + oldRef);
                    System.out.println("existing: " + existing);
                    this.printRefDebug();
                    System.out.println(">>(create) " + elem.hashCode() + " => " + newRef);
                    System.out.println();
                } // TODO: REMOVE

                this.addRef(newRef);
                ref = newRef;
            }

            // Only set "in progress" to true here so that only we can differentiate
            // from the methods that only create a temporary reference.
            this.elemInProg.add(elem.hashCode());

            // Create a new construct for this data.
            final T newCon = creator.create();
            if (newCon == null)
                throw new Exception("Factory creator for " + this.toString() + " returned null.");
            if (!newCon.kind().equals(this.conKind))
                throw new Exception("Factory creator for " + this.toString() + " create a type with kind " + newCon.kind() + ".");

            // If an existing construct matches the new one after the new one
            // has been loaded, then there are two elements to get to the same
            // value. Resolve the reference for the existing or new construct.
            // Run finisher on both since the element is different, it may have
            // different finishing steps.
            final T other = this.getExisting(newCon);
            if (other != null) {
                ref.setResolved(other);
                if (finisher != null) finisher.finish(ref, other);
            } else {
                this.conSet.add(newCon);
                ref.setResolved(newCon);
                if (finisher != null) finisher.finish(ref, newCon);
            }

            return ref;
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
        this.byElem.remove(elem.hashCode());

        if (ref == null) return;
        this.refSet.remove(ref);

        if (!ref.isResolved()) return;
        final T con = ref.getResolved();
        this.conSet.remove(con);
        this.nonElemRef.remove(con);
    }

    public void setRefForElem(CtElement elem, Ref<T> ref) throws Exception {
        final Ref<T> existing = this.getRef(elem);
        if (existing != null) {
            if (existing.equals(ref)) return;
            throw new Exception("Ref already exists for element " + existing + " so cannot set " + ref + ".");
        }

        if (elem.hashCode() == 3002523) {
            System.out.println(">>(setRefForElem) "+elem.hashCode() + " => " + ref); // TODO: REMOVE
        }
        
        this.addRef(ref);
    }

    public Ref<T> addOfGetRefForElem(CtElement elem, String title) throws Exception {
        final Ref<T> existing = this.getRef(elem);
        if (existing != null) return existing;

        final Ref<T> ref = new Ref<T>(this.conKind, elem, title);

        if (elem.hashCode() == 3002523) {
            System.out.println(">>(addOfGetRefForElem) " + elem.hashCode() + " => " + ref); // TODO: REMOVE
        }

        this.addRef(ref);
        return ref;
    }

    static private CmpOptions resolvedCmpOptionsSingleton = null;
    private CmpOptions resolvedCmpOptions() {
        if (resolvedCmpOptionsSingleton != null) return resolvedCmpOptionsSingleton;
        CmpOptions options = new CmpOptions();
        options.useResolved = true;
        resolvedCmpOptionsSingleton = options;
        return options;
    }

    public Ref<T> addOrGetRef(T c, String context) throws Exception {
        final T other = this.getExisting(c);
        if (other != null) c = other;

        Ref<T> ref = this.nonElemRef.get(c);
        if (ref != null) return ref;

        ref = new Ref<T>(this.conKind, null, "no element ref: " + context);
        ref.setResolved(c);
        ref.setCmpOptions(resolvedCmpOptions());

        final Ref<T> otherRef = this.refSet.floor(ref);
        if (ref.equals(otherRef)) return otherRef;

        this.addRef(ref);
        return ref;
    }

    //==========================================================================

    public boolean consolidateCons(Logger log) throws Exception {
        // Copy all cons to a list and clear the set.
        final ArrayList<T> conList = new ArrayList<T>(this.conSet);
        this.conSet.clear();

        CmpOptions options = resolvedCmpOptions();
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
