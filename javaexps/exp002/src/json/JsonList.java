package json;

import java.util.ArrayList;

public final class JsonList extends ArrayList<JsonObj> implements JsonObj {

    public JsonList() {}

    public JsonList(Iterable elements) {
        elements.forEach((e) -> this.add(JsonObj.convert(e)));
    }

    public boolean addValue(Object value) {
        return super.add(JsonObj.convert(value));
    }

    public JsonObj setValue(int index, Object value) {
        return super.set(index, JsonObj.convert(value));
    }

    @Override
    public void write(StringBuilder buf) {
        buf.append("[");
        boolean first = true;
        for (JsonObj elem : this) {
            if (first) first = false;
            else buf.append(",");
            JsonObj.write(buf, elem);
        }
        buf.append("]");
    }
}
