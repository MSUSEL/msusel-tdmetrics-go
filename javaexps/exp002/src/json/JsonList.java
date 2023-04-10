package json;

import java.util.ArrayList;

public final class JsonList extends ArrayList<JsonObj> implements JsonObj {
    private boolean omitOnEmpty;

    public JsonList() { }

    public JsonList(Iterable<?> elements) {
        elements.forEach((e) -> this.add(JsonObj.convert(e)));
    }

    public JsonMap addMap() {
        JsonMap m = new JsonMap();
        this.add(m);
        return m;
    }

    public JsonList with(Object value) {
        super.add(JsonObj.convert(value));
        return this;
    }

    public JsonList withOmitOnDefault(Object value) {
        JsonObj jObj = JsonObj.convert(value);
        jObj.setOmitOnDefault(true);
        this.add(jObj);
        return this;
    }

    @Override
    public void setOmitOnDefault(boolean omit) {
        this.omitOnEmpty = omit;
    }

    @Override
    public boolean omit() {
        if (!this.omitOnEmpty) return false;
        for (JsonObj elem : this) {
            if (!elem.omit()) return false;
        }
        return true;
    }

    @Override
    public void write(StringBuilder buf) {
        buf.append("[");
        boolean first = true;
        for (JsonObj elem : this) {
            if (elem.omit()) continue;
            if (first) first = false;
            else buf.append(",");
            JsonObj.write(buf, elem);
        }
        buf.append("]");
    }

    @Override
    public boolean equals(Object o) {
        if (o instanceof JsonList other) {
            if (this.size() != other.size()) return false;
            for (int i = 0; i < this.size(); ++i) {
                JsonObj elem1 = this.get(i);
                JsonObj elem2 = other.get(i);
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
