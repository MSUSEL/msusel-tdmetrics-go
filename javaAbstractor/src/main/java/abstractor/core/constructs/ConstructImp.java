package abstractor.core.constructs;

import java.util.Collections;
import java.util.List;
import java.util.SortedSet;
import java.util.TreeSet;

import abstractor.core.json.*;

public abstract class ConstructImp implements Construct {

    static public JsonValue key(Construct c) {
        if (c == null) return JsonValue.ofNull();
        if (c instanceof Reference ref && ref.isResolved())
            return key(ref.getResolved());
        return JsonValue.of(c.kind().toString() + c.getIndex());
    }
    
    static public JsonValue index(Construct c) {
        if (c == null) return JsonValue.ofNull();
        if (c instanceof Reference ref && ref.isResolved())
            return index(ref.getResolved());
        return JsonValue.of(c.getIndex());
    }

    static public <T extends Construct> JsonNode keyList(Iterable<T> os) {
        if (os == null) return JsonValue.ofNull();

        JsonArray indices = new JsonArray();
        for (Construct o: os) indices.add(key(o));
        return indices;
    }

    static public <T extends Construct> JsonNode indexList(Iterable<T> os) {
        if (os == null) return JsonValue.ofNull();

        JsonArray indices = new JsonArray();
        for (Construct o: os) indices.add(index(o));
        return indices;
    }

    static public <T extends Construct> JsonNode keySet(Iterable<T> os) {
        if (os == null) return JsonValue.ofNull();

        TreeSet<T> set = new TreeSet<T>();
        for (T o: os) set.add(o);

        JsonArray indices = new JsonArray();
        for (T s: set) indices.add(key(s));
        return indices;
    }

    static public <T extends Construct> JsonNode indexSet(Iterable<T> os) {
        if (os == null) return JsonValue.ofNull();

        TreeSet<T> set = new TreeSet<T>();
        for (T o: os) set.add(o);

        JsonArray indices = new JsonArray();
        for (T s: set) indices.add(index(s));
        return indices;
    }

    static public <T> List<T> unmodifiableList(List<T> list) {
        return list == null ? null : Collections.unmodifiableList(list);
    }

    static public <T> SortedSet<T> unmodifiableSortedSet(SortedSet<T> set) {
        return set == null ? null : Collections.unmodifiableSortedSet(set);
    }

    private int index;

    protected ConstructImp() { }

    public void setIndex(int index) { this.index = index; }
    public int getIndex() { return this.index; }

    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = new JsonObject();
        if (h.writeKinds) obj.put("kind", this.kind().toString());
        if (h.writeIndices) obj.put("index", this.getIndex());
        return obj;
    }

    public abstract ConstructKind kind();

    public int compareTo(Construct c) {
        return this.kind().compareTo(c.kind());
    }

    @Override
    public boolean equals(Object obj) {
        return obj instanceof Construct c && this.compareTo(c) == 0;
    }

    @Override
    public String toString() {
        JsonHelper jh = new JsonHelper();
        jh.writeKinds = true;
        jh.writeIndices = true;
        return JsonFormat.Relaxed().format(this.toJson(jh));
    }
}
