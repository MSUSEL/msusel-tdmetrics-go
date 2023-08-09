package json;

import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.util.TreeMap;
import java.util.Set;
import java.util.HashSet;
import java.util.stream.Collectors;

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
    public void removeOmitted() {
        List<String> keys = new ArrayList();
        for (Map.Entry<String, JsonObj> pair : this.entrySet()) {
            final JsonObj elem = pair.getValue();
            if (elem.omit()) keys.add(pair.getKey());
            else elem.removeOmitted();
        }
        for (String key: keys) this.remove(key);
    }

    @Override
    public boolean equals(Object o) {
        if (o instanceof JsonMap other) {
            if (this.size() != other.size()) return false;
            for (Map.Entry<String, JsonObj> pair : this.entrySet()) {
                final JsonObj elem1 = pair.getValue();
                final JsonObj elem2 = other.get(pair.getKey());
                if (elem1 != null) {
                   if (!elem1.equals(elem2)) return false;
                } else if (elem2 != null) return false;
            }
            return true;
        }
        return false;
    }

    @Override
    public void assertCompare(JsonObj obj) throws Exception {
        if (!(obj instanceof JsonMap other))
            throw new Exception("Expected type to be " + this.getClass().getName() + " but got " + obj.getClass().getName());

        Set<String> missing = new HashSet<>();
        Set<String> extra   = new HashSet<>(other.keySet());
        for (Map.Entry<String, JsonObj> pair : this.entrySet()) {
            final String key = pair.getKey();
            if (!other.containsKey(key)) {
                missing.add(key);
                continue;
            }
            extra.remove(key);

            final JsonObj elem1 = pair.getValue();
            final JsonObj elem2 = other.get(key);
            if (elem1 == null && elem2 != null)
                throw new Exception("Expected value at " + JsonObj.escape(key) + " was null but got " + elem2);

            if (elem2 == null)
                throw new Exception("Expected value at " + JsonObj.escape(key) + " was " + elem1 + " but got null");

            try {
                elem1.assertCompare(elem2);
            } catch (Exception e) {
                throw new Exception("At " + JsonObj.escape(key) + ": " + e.getMessage());
            }
        }

        if (missing.size() > 0 || extra.size() > 0) {
            String missingStr = missing.stream().map(JsonObj::escape).collect(Collectors.joining(", ", "[", "]"));
            String extraStr = extra.stream().map(JsonObj::escape).collect(Collectors.joining(", ", "[", "]"));
            throw new Exception("Expected missing keys " + missingStr + " but got extra keys " + extraStr);
        }
    }

    @Override
    public String toString() {
        return JsonObj.toString(this);
    }
}
