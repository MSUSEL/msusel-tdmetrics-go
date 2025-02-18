package abstractor.core.constructs;

import java.util.TreeSet;

import abstractor.core.json.*;

public abstract class ConstructImp implements Construct {

    static public JsonValue key(Construct c) {
        return c == null ? JsonValue.ofNull() : JsonValue.of(c.kind() + c.getIndex());
    }
    
    static public JsonValue index(Construct c) {
        return c == null ? JsonValue.ofNull() : JsonValue.of(c.getIndex());
    }

    static public JsonArray keyList(Iterable<? extends Construct> os) {
        JsonArray indices = new JsonArray();
        for (Construct o: os) indices.add(key(o));
        return indices;
    }

    static public JsonArray indexList(Iterable<? extends Construct> os) {
        JsonArray indices = new JsonArray();
        for (Construct o: os) indices.add(index(o));
        return indices;
    }

    static public JsonArray keySet(Iterable<? extends Construct> os) {
        TreeSet<JsonValue> set = new TreeSet<JsonValue>();
        for (Construct o: os) set.add(key(o));
        
        JsonArray indices = new JsonArray();
        for (JsonValue s: set) indices.add(s);
        return indices;
    }

    static public JsonArray indexSet(Iterable<? extends Construct> os) {
        TreeSet<JsonValue> set = new TreeSet<JsonValue>();
        for (Construct o: os) set.add(index(o));

        JsonArray indices = new JsonArray();
        for (JsonValue s: set) indices.add(s);
        return indices;
    }

    private int index;

    public void setIndex(int index) { this.index = index; }
    public int getIndex() { return this.index; }

    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = new JsonObject();
        if (h.writeKinds) obj.put("kind", this.kind());
        if (h.writeIndices) obj.put("index", this.getIndex());
        return obj;
    }

    public abstract Object source();
    public abstract String kind();

    public int compareTo(Construct c) {
        return this.kind().compareTo(c.kind());
    }

    @Override
    public boolean equals(Object obj) {
        return obj instanceof Construct c && this.compareTo(c) == 0;
    }
}
