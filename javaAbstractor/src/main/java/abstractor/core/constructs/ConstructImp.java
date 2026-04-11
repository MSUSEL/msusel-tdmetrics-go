package abstractor.core.constructs;

import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;

public abstract class ConstructImp implements Construct {

    static public JsonValue key(Construct c) {
        return c == null ? JsonValue.ofNull() : JsonValue.of(c.kind().toString() + c.getIndex());
    }
    
    static public JsonValue index(Construct c) {
        return c == null ? JsonValue.ofNull() : JsonValue.of(c.getIndex());
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

    private int index = -1;
    private CmpOptions options = null;

    protected ConstructImp() { }

    public void setIndex(int index) { this.index = index; }
    public int  getIndex()          { return this.index;  }

    public abstract ConstructKind kind();

    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = new JsonObject();
        if (h.writeKinds)   obj.put("kind",  this.kind().toString());
        if (h.writeIndices) obj.put("index", this.getIndex());
        return obj;
    }

    public void setCmpOptions(CmpOptions options) { this.options = options; }
    public CmpOptions getCmpOptions()             { return this.options;    }

    public int compareTo(Construct c) {
        return Cmp.compareTo(this, c, this.options);
    }
    
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.defer(this.kind(), () -> c.kind());
    }

    @Override
    public boolean equals(Object obj) {
        return obj instanceof Construct c && this.compareTo(c) == 0;
    }

    @Override
    public String toString() {
        JsonHelper jh = new JsonHelper();
        jh.writeKinds   = true;
        jh.writeIndices = true;
        return JsonFormat.Relaxed().format(this.toJson(jh));
    }
}
