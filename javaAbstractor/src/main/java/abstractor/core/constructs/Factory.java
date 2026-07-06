package abstractor.core.constructs;

import java.util.*;
import java.util.function.Predicate;

import abstractor.core.ElementKey;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;
import abstractor.core.log.*;
import abstractor.core.require.Require;

public class Factory<T extends Construct> implements Jsonable {
    static private final boolean logCreate = true;

    private final ConstructKind conKind;

    private final TreeMap<ElementKey, Ref<T>> byElem     = new TreeMap<>();
    private final TreeMap<T,          Ref<T>> nonElemRef = new TreeMap<>();
    private final TreeSet<ElementKey>         elemInProg = new TreeSet<>();
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

    public Ref<T> getRefByElem(ElementKey elem) {
        return this.byElem.get(elem);
    }
    
    public T getExisting(T c) {
        final T other = this.conSet.floor(c);
        return c.equals(other) ? other : null;
    }

    //==========================================================================

    @FunctionalInterface
    public interface Creator<T extends Construct> { T create() throws Exception; }
    
    @FunctionalInterface
    public interface Finisher<T extends Construct> { void finish(Ref<T> ref, T con) throws Exception; }

    public Ref<T> create(Logger log, ElementKey elemKey, String title, Creator<T> creator, Finisher<T> finisher) throws Exception {
        Require.notNull(elemKey, "factories' create methods require non-null element keys");

        // If already "in progress" then check for if a reference already exists
        // so that we only create one and all others are references. However,
        // since references can be created other ways, we need to skip checking
        // for references if not "in progress" to start progress.
        //
        // `elemInProg` is set when "create" is called to differentiate between
        // just creating references and creating the actual construct. It is set
        // then never cleared so that this method will only let one in then continue
        // to use the reference correctly. It is "in progress" since if the construct is
        // recursive it will be in progress when the same construct calls "create" from
        // inside the creator function.
        final Ref<T> existing = this.getRefByElem(elemKey);
        final boolean inProgress = this.elemInProg.contains(elemKey);
        if (inProgress && existing != null) {
            //log.log("Got reference for " + title);
            return existing;
        }
        
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
                final Ref<T> newRef = new Ref<T>(this.conKind, elemKey, title);
                this.addRefWithElem(newRef);
                ref = newRef;
            }

            // Only set "in progress" to true here so that only we can differentiate
            // from the methods that only create a temporary reference.
            this.elemInProg.add(elemKey);

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

    public Ref<T> create(Logger log, ElementKey elemKey, String title, Creator<T> creator) throws Exception {
        return this.create(log, elemKey, title, creator, null);
    }

    public void removeIf(Logger log, Predicate<T> predicate) {
        if (predicate == null) return;
        List<T> toRemove = new LinkedList<>();
        for (T con : this.conSet) {
            if (predicate.test(con)) toRemove.add(con);
        }
        for (T con : toRemove) this.remove(log, con);
    }

    public void remove(Logger log, T con) {
        if (con == null) return;
        if (logCreate) log.log("Removing " + con);

        con.setIndex(-100);
        this.conSet.remove(con);
        this.nonElemRef.remove(con);

        Predicate<Ref<T>> refRemover = (Ref<T> ref) -> {
            return ref != null && ref.isResolved() && ref.getResolved().equals(con);
        };

        this.refSet.removeIf(refRemover);

        final Iterator<Map.Entry<ElementKey, Ref<T>>> it = this.byElem.entrySet().iterator();
        while (it.hasNext()) {
            if (refRemover.test(it.next().getValue())) it.remove();
        }
    }

    /**
     * Adds a new reference that has an element in it.
     * 
     * This should only be used by the factory when
     * adding newly created references with elements.
     */
    private void addRefWithElem(Ref<T> ref) throws Exception {
        final ElementKey elemKey = ref.elemKey;
        Require.notNull(elemKey, "element key may not be null when adding the reference " + ref);
        Require.require(this.refSet.add(ref), "reference " + ref + " must be added at this point");
        Require.isNull(this.byElem.put(elemKey, ref));
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
    public Ref<T> setRefForElem(ElementKey elemKey, Ref<T> ref) throws Exception {
        final Ref<T> existing = this.getRefByElem(elemKey);
        if (existing != null) {
            Require.equal(existing, ref,
                "reference already exists for element " + existing + " so cannot set " + ref);
            return existing;
        }

        final Ref<T> otherRef = this.refSet.floor(ref);
        if (ref.equals(otherRef)) {
            Require.isNull(this.byElem.put(elemKey, otherRef));
            return otherRef;
        }

        Require.require(this.refSet.add(ref));
        Require.isNull(this.byElem.put(elemKey, ref));
        return ref;
    }

    /**
     * Gets existing reference for the given element.
     * If no reference for that element exists, then one will be created, added, and returned.
     *
     * This is used to create a reference before the actual creation of the construct is called.
     * For example when creating a reference for something pending to be created later, like a package.
     */
    public Ref<T> addOrGetRefForElem(ElementKey elemKey, String title) throws Exception {
        final Ref<T> existing = this.getRefByElem(elemKey);
        if (existing != null) return existing;

        final Ref<T> ref = new Ref<T>(this.conKind, elemKey, title);
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
    public Ref<T> addOrGetRef(T c, List<Ref<? extends TypeDesc>> typeArgs, String context) throws Exception {
        final T other = this.getExisting(c);
        if (other != null) c = other;

        final Ref<T> ref = this.nonElemRef.get(c);
        if (ref != null) return ref;

        final ElementKey elemKey = new ElementKey(null, typeArgs);
        final Ref<T> newRef = new Ref<T>(this.conKind, elemKey, "no element ref: " + context);
        newRef.setResolved(c);
        newRef.setCmpOptions(CmpOptions.resolvedCmp());

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

                JsonArray refList = new JsonArray();
                for (Ref<T> ref : this.refSet) {
                    if (t.equals(ref.getResolved()))
                        refList.add(ref.refToJson(h));
                }
                obj.put("refs", refList);
                node = obj;
            }
            array.add(node);
        }
        return array;
    }
}
