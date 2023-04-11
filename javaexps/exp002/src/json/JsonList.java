package json;

import java.util.ArrayList;
import java.util.Map;

import static java.lang.Math.min;

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
    public void removeOmitted() {
        for (int i = this.size()-1; i >= 0; --i) {
            final JsonObj elem = this.get(i);
            if (elem.omit()) this.remove(i);
            else elem.removeOmitted();
        }
    }

    @Override
    public boolean equals(Object o) {
        if (o instanceof JsonList other) {
            final int size = this.size();
            if (size != other.size()) return false;
            for (int i = 0; i < size; ++i) {
                final JsonObj elem1 = this.get(i);
                final JsonObj elem2 = other.get(i);
                if ((elem1 == null && elem2 != null) ||
                    (!elem1.equals(elem2)))
                    return false;
            }
            return true;
        }
        return false;
    }

    @Override
    public void assertCompare(JsonObj obj) throws Exception {
        if (!(obj instanceof JsonList other))
            throw new Exception("Expected type to be " + this.getClass().getName() + " but got " + obj.getClass().getName());

        final int size = min(this.size(), other.size());
        for (int i = 0; i < size; ++i) {
            final JsonObj elem1 = this.get(i);
            final JsonObj elem2 = other.get(i);

            if (elem1 == null && elem2 != null)
                throw new Exception("Expected value at [" + i + "] was null but got " + elem2);

            if (elem2 == null)
                throw new Exception("Expected value at [" + i + "] was " + elem1 + " but got null");

            try {
                elem1.assertCompare(elem2);
            } catch (Exception e) {
                throw new Exception("At [" + i + "]: " + e.getMessage());
            }
        }

        if (this.size() != other.size())
            throw new Exception("Expected " + this.size() + " values but got " + other.size());
    }

    @Override
    public String toString() {
        return JsonObj.toString(this);
    }
}
