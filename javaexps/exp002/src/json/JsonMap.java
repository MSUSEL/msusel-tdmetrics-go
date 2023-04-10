package json;

import java.util.Map;
import java.util.TreeMap;

public class JsonMap extends TreeMap<String, JsonObj> implements JsonObj {
    private boolean omitOnEmpty;

    public JsonMap() {}

    public JsonMap(Map<String, ?> data) {
        for (String key : data.keySet())
            this.put(key, data.get(key));
    }

    public void put(String key, Object value) {
        super.put(key, JsonObj.convert(value));
    }

    public JsonMap addMap(String key) {
        JsonMap m = new JsonMap();
        super.put(key, m);
        return m;
    }

    public JsonList addList(String key) {
        JsonList m = new JsonList();
        super.put(key, m);
        return m;
    }

    public JsonMap with(String key, Object value) {
        this.put(key, value);
        return this;
    }

    public JsonMap withOmitOnDefault(String key, Object value) {
        JsonObj jObj = JsonObj.convert(value);
        jObj.setOmitOnDefault(true);
        this.put(key, jObj);
        return this;
    }

    @Override
    public void setOmitOnDefault(boolean omit) {
        this.omitOnEmpty = omit;
    }

    @Override
    public boolean omit() {
        if (!this.omitOnEmpty) return false;
        for (JsonObj val : this.values()) {
            if (!val.omit()) return false;
        }
        return true;
    }

    @Override
    public void write(StringBuilder buf) {
        buf.append("{");
        boolean first = true;
        for (Map.Entry<String, JsonObj> pair : this.entrySet()) {
            if (pair.getValue().omit()) continue;
            if (first) first = false;
            else buf.append(",");
            buf.append(JsonObj.escape(pair.getKey()));
            buf.append(":");
            JsonObj.write(buf, pair.getValue());
        }
        buf.append("}");
    }

    @Override
    public boolean equals(Object o) {
        if (o instanceof JsonMap other) {
            if (this.size() != other.size()) return false;
            for (Map.Entry<String, JsonObj> pair : this.entrySet()) {
                JsonObj elem1 = pair.getValue();
                JsonObj elem2 = other.get(pair.getKey());
                if (elem1 == null) return elem2 == null;
                if (elem2 == null) return false;
                if (!elem1.equals(elem2)) return false;
            }
            return true;
        }
        return false;
    }

    @Override
    public String toString() {
        return JsonObj.toString(this);
    }
}
