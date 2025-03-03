package abstractor.core.constructs;

import java.util.Collections;
import java.util.List;
import java.util.SortedSet;
import java.util.TreeSet;

import abstractor.core.json.*;

public abstract class ConstructImp implements Construct {

    static public JsonValue key(Construct c) {
        return c == null ? JsonValue.ofNull() : JsonValue.of(c.kind() + c.getIndex());
    }
    
    static public JsonValue index(Construct c) {
        return c == null ? JsonValue.ofNull() : JsonValue.of(c.getIndex());
    }

    static public JsonNode keyList(Iterable<? extends Construct> os) {
        if (os == null) return JsonValue.ofNull();
        
        JsonArray indices = new JsonArray();
        for (Construct o: os) indices.add(key(o));
        return indices;
    }

    static public JsonNode indexList(Iterable<? extends Construct> os) {
        if (os == null) return JsonValue.ofNull();

        JsonArray indices = new JsonArray();
        for (Construct o: os) indices.add(index(o));
        return indices;
    }

    static public JsonNode keySet(Iterable<? extends Construct> os) {
        if (os == null) return JsonValue.ofNull();

        TreeSet<JsonValue> set = new TreeSet<JsonValue>();
        for (Construct o: os) set.add(key(o));
        
        JsonArray indices = new JsonArray();
        for (JsonValue s: set) indices.add(s);
        return indices;
    }

    static public JsonNode indexSet(Iterable<? extends Construct> os) {
        if (os == null) return JsonValue.ofNull();

        TreeSet<JsonValue> set = new TreeSet<JsonValue>();
        for (Construct o: os) set.add(index(o));

        JsonArray indices = new JsonArray();
        for (JsonValue s: set) indices.add(s);
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
        if (h.writeKinds) obj.put("kind", this.kind());
        if (h.writeIndices) obj.put("index", this.getIndex());
        return obj;
    }

    public abstract String kind();

    public int compareTo(Construct c) {
        return this.kind().compareTo(c.kind());
    }

    @Override
    public boolean equals(Object obj) {
        return obj instanceof Construct c && this.compareTo(c) == 0;
    }
}
