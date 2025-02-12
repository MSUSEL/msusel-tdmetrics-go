package abstractor.core.json;

import java.io.PrintStream;
import java.io.ByteArrayOutputStream;
import java.util.TreeMap;

public class JsonObject extends TreeMap<String, JsonNode> implements JsonNode {
    
    public boolean isSimple() {
        if (this.size() > 20) return false;
        for (JsonNode elem : this.values()) {
            if (!(elem instanceof JsonValue))
                return false;
        }
        return true;
    }

    public void put(String key, boolean value) { this.put(key, JsonValue.of(value)); }

    public void put(String key, int value) { this.put(key, JsonValue.of(value)); }
    
    public void put(String key, double value) { this.put(key, JsonValue.of(value)); }
    
    public void put(String key, String value) { this.put(key, JsonValue.of(value)); }

    public void putNotEmpty(String key, JsonNode value) {
        if (!value.isEmpty()) this.put(key, value);
    }

    public void putNotEmpty(String key, boolean value) { this.putNotEmpty(key, JsonValue.of(value)); }

    public void putNotEmpty(String key, int value) { this.putNotEmpty(key, JsonValue.of(value)); }
    
    public void putNotEmpty(String key, double value) { this.putNotEmpty(key, JsonValue.of(value)); }
    
    public void putNotEmpty(String key, String value) { this.putNotEmpty(key, JsonValue.of(value)); }

    public void toString(PrintStream sb, boolean minimize, String indent) {
        if (isEmpty()) {
            sb.append(minimize? "{}": "{ }");
            return;
        }

        sb.append("{");
        boolean first = true;
        final boolean simple = this.isSimple();
        final String indent2 = indent + "  ";
        final String separator = simple ? " " : "\n"+indent2;
        for (Object key : this.keySet()) {
            if (first) first = false;
            else sb.append(",");
            if (!minimize) sb.append(separator);

            sb.append(JsonValue.escape((String)key));
            sb.append(":");
            if (!minimize) sb.append(" ");
            this.get(key).toString(sb, minimize, indent2);            
        }
        if (!minimize) sb.append(simple ? " " : "\n"+indent);
        sb.append("}");
    }
    
    public String toString(boolean minimize) {
        ByteArrayOutputStream buf = new ByteArrayOutputStream();
        this.toString(new PrintStream(buf), minimize, "");
        return buf.toString();
    }

    public String toString() { return this.toString(false); }
}
