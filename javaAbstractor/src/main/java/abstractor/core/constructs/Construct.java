package abstractor.core.constructs;

import java.util.TreeSet;

import abstractor.core.json.*;

public abstract class Construct implements Comparable<Construct>, Jsonable {

    static public JsonArray keyList(Iterable<? extends Construct> os) {
        JsonArray indices = new JsonArray();
        for (Construct o: os) indices.add(o.getKey());
        return indices;
    }

    static public JsonArray indexList(Iterable<? extends Construct> os) {
        JsonArray indices = new JsonArray();
        for (Construct o: os) indices.add(o.getIndex());
        return indices;
    }

    static public JsonArray keySet(Iterable<? extends Construct> os) {
        TreeSet<String> set = new TreeSet<String>();
        for (Construct o: os) set.add(o.getKey());
        
        JsonArray indices = new JsonArray();
        for (String s: set) indices.add((String)s);
        return indices;
    }

    static public JsonArray indexSet(Iterable<? extends Construct> os) {
        TreeSet<Integer> set = new TreeSet<Integer>();
        for (Construct o: os) set.add(o.getIndex());

        JsonArray indices = new JsonArray();
        for (int s: set) indices.add(s);
        return indices;
    }

    private int index;

    public void setIndex(int index) { this.index = index; }
    public int getIndex() { return this.index; }
    public String getKey() { return this.kind() + this.index; }

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
