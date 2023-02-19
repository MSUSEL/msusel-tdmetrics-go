package json;

import java.util.Map;
import java.util.TreeMap;

public class JsonMap extends TreeMap<String, JsonObj> implements JsonObj {

    public JsonMap() {}

    public JsonMap(Map<String, ?> data) {
        for (String key : data.keySet())
            this.put(key, data.get(key));
    }

    public void put(String key, Object value) {
        super.put(key, JsonObj.convert(value));
    }

    public void write(StringBuilder buf) {
        buf.append("{");
        boolean first = true;
        for (Map.Entry<String, JsonObj> pair : this.entrySet()) {
            if (first) first = false;
            else buf.append(",");
            buf.append(JsonObj.escape(pair.getKey()));
            buf.append(":");
            JsonObj.write(buf, pair.getValue());
        }
        buf.append("}");
    }
}
