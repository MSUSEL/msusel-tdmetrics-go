package abstractor.core.constructs;

import java.util.*;

import spoon.reflect.declaration.CtElement;

import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;
import abstractor.core.log.*;
import abstractor.core.require.Require;

public class Factory<T extends Construct> implements Jsonable {
    static private final boolean logCreate = true;

    private final ConstructKind conKind;

    private final HashMap<CtElement, Ref<T>> byElem     = new HashMap<>();
    private final HashMap<T,         Ref<T>> nonElemRef = new HashMap<>();
    private final HashSet<CtElement>         elemInProg = new HashSet<>();
    private final TreeSet<Ref<T>>            refSet     = new TreeSet<>();
    private final TreeSet<T>                 conSet     = new TreeSet<>();

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

    public Ref<T> getRefByElem(CtElement elem) {
        return this.byElem.get(elem);
    }
    
    public T getExisting(T c) {
        final T other = this.conSet.floor(c);
        return c.equals(other) ? other : null;
    }

    private List<Ref<T>> findRefsForCon(T con) {
        return this.refSet.stream().filter(r -> con.equals(r.getResolved())).toList();
    }

    //==========================================================================

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
        final Ref<T> existing = this.getRefByElem(elem);
        final boolean inProgress = this.elemInProg.contains(elem);
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
                this.addRefWithElem(newRef);
                ref = newRef;
            }

            // Only set "in progress" to true here so that only we can differentiate
            // from the methods that only create a temporary reference.
            this.elemInProg.add(elem);

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
                Require.require(this.conSet.add(newCon));
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

        final Ref<T> ref = this.getRefByElem(elem);
        this.byElem.remove(elem);

        if (ref == null) return;
        this.refSet.remove(ref);

        if (!ref.isResolved()) return;
        final T con = ref.getResolved();
        this.conSet.remove(con);
        this.nonElemRef.remove(con);
    }

    /**
     * Adds a new reference that has an element in it.
     * 
     * This should only be used by the factory when
     * adding newly created references with elements.
     */
    private void addRefWithElem(Ref<T> ref) throws Exception {
        final CtElement elem = ref.elem;
        Require.notNull(elem, "element may not be null when adding the reference " + ref);
        Require.require(this.refSet.add(ref), "reference " + ref + " must be added at this point");
        Require.isNull(this.byElem.put(elem, ref));
    }

    /**
     * Sets an existing reference to an element this it doesn't have in it.
     * If an element already exists as a reference, this it will be checked
     * that the reference isn't changing instead.
     * 
     * This returns an existing equivalent reference set of that element or
     * the given reference if it was added.
     *
     * For example, when an array is instantiated for a specific element type,
     * the instantiated array then has the element for the array set for it.
     */
    public Ref<T> setRefForElem(CtElement elem, Ref<T> ref) throws Exception {
        final Ref<T> existing = this.getRefByElem(elem);
        if (existing != null) {
            Require.equal(existing, ref,
                "reference already exists for element " + existing + " so cannot set " + ref);
            return existing;
        }

        final Ref<T> otherRef = this.refSet.floor(ref);
        if (ref.equals(otherRef)) {
            Require.isNull(this.byElem.put(elem, otherRef));
            return otherRef;
        }

        Require.require(this.refSet.add(ref));
        Require.isNull(this.byElem.put(elem, ref));
        return ref;
    }

    /**
     * Gets existing reference for the given element.
     * If no reference for that element exists, then one will be created, added, and returned.
     *
     * This is used to create a reference before the actual creation of the construct is called.
     * For example when creating a reference for something pending to be created later, like a package.
     */
    public Ref<T> addOrGetRefForElem(CtElement elem, String title) throws Exception {
        final Ref<T> existing = this.getRefByElem(elem);
        if (existing != null) return existing;

        final Ref<T> ref = new Ref<T>(this.conKind, elem, title);
        this.addRefWithElem(ref);
        return ref;
    }

    /**
     * Gets the reference for the given construct.
     * If no reference for the given construct exists, then a new reference with
     * no element is created for this construct and set as resolved with the construct.
     *
     * This is used when a construct is generated or baked such that there is
     * no element, or at least no element yet, for the construct.
     */
    public Ref<T> addOrGetRef(T c, String context) throws Exception {
        final T other = this.getExisting(c);
        if (other != null) c = other;

        final Ref<T> ref = this.nonElemRef.get(c);
        if (ref != null) return ref;

        final Ref<T> newRef = new Ref<T>(this.conKind, null, "no element ref: " + context);
        newRef.setResolved(c);
        newRef.setCmpOptions(resolvedCmpOptions());

        final Ref<T> otherRef = this.refSet.floor(newRef);
        if (newRef.equals(otherRef)) return otherRef;

        // construct may already exist because it was added with an element,
        // like what happens with adding an `int` from the AST and one from the baker.
        this.conSet.add(c);

        Require.require(this.refSet.add(newRef),
            "reference " + newRef + " must be added at this point for non-element ref, otherwise it should have returned before now");       
        Require.isNull(this.nonElemRef.put(c, newRef),
            "resolved construct " + c + " and reference " + newRef + " must be added at this point for non-element ref");
        return newRef;
    }

    //==========================================================================

    static private CmpOptions resolvedCmpOptionsSingleton = null;
    static private CmpOptions resolvedCmpOptions() {
        if (resolvedCmpOptionsSingleton != null) return resolvedCmpOptionsSingleton;
        CmpOptions options = new CmpOptions();
        options.useResolved = true;
        resolvedCmpOptionsSingleton = options;
        return options;
    }

    /**
     * Change all the comparison options to use the resolved.
     */
    public void setToCompareResolved() {
        final CmpOptions options = resolvedCmpOptions();
        for (T con : this.conSet) con.setCmpOptions(options);
        for (Ref<T> ref : this.refSet) ref.setCmpOptions(options);
        
        // The non-element references is no longer useful but the changed comparisons
        // could cause issues if someone tried to use it so just clear it out.
        this.nonElemRef.clear();
    }

    public boolean consolidateCons(Logger log) throws Exception {
        // Copy all cons to a list and clear the set so that only
        // the unique cons can be re-added in the new sort order.
        final ArrayList<T> conList = new ArrayList<T>(this.conSet);
        this.conSet.clear();

        boolean collision = false;
        for (T con : conList) {
            T existing = this.conSet.floor(con);
            if (existing == null || !existing.equals(con)) {
                // No conflict found, so add the construct into set.
                Require.require(this.conSet.add(con));
                continue;
            }

            // Found another construct that is equal so move all references over
            // to the existing construct since the duplicate is about to be removed.
            collision = true;
            List<Ref<T>> refs = this.findRefsForCon(con);
            for (Ref<T> ref : refs) ref.setResolved(existing);
            con.setIndex(-100);
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
