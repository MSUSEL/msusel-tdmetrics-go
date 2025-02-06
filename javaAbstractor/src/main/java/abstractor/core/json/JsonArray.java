package abstractor.core.json;

import java.io.*;
import java.util.ArrayList;

public class JsonArray extends ArrayList<JsonNode> implements JsonNode {
    
    public JsonArray(JsonNode ...nodes) {
        for(JsonNode node: nodes) this.add(node);
    }
    
    public void add(boolean value) { this.add(JsonValue.of(value)); }

    public void add(int value) { this.add(JsonValue.of(value)); }
    
    public void add(double value) { this.add(JsonValue.of(value)); }
    
    public void add(String value) { this.add(JsonValue.of(value)); }

    public void addNotEmpty(JsonNode value) {
        if (!value.isEmpty()) this.add(value);
    }

    public void addNotEmpty(boolean value) { this.addNotEmpty(JsonValue.of(value)); }

    public void addNotEmpty(int value) { this.addNotEmpty(JsonValue.of(value)); }
    
    public void addNotEmpty(double value) { this.addNotEmpty(JsonValue.of(value)); }
    
    public void addNotEmpty(String value) { this.addNotEmpty(JsonValue.of(value)); }

    public boolean isSimple() {
        for (JsonNode elem : this) {
            if (!(elem instanceof JsonValue))
                return false;
        }
        return true;
    }

    public void toString(PrintStream sb, boolean minimize, String indent) {
        if (isEmpty()) {
            sb.append(minimize? "[]": "[ ]");
            return;
        }

        sb.append("[");
        boolean first = true;
        final boolean simple = this.isSimple();
        final String indent2 = indent + "  ";
        final String separator = simple ? " " : "\n"+indent2;
        for (JsonNode elem : this) {
            if (first) first = false;
            else sb.append(",");
            if (!minimize) sb.append(separator);
            
            elem.toString(sb, minimize, indent2);            
        }
        if (!minimize) sb.append(simple ? " " : "\n"+indent);
        sb.append("]");
    }

    public String toString(boolean minimize) {
        ByteArrayOutputStream buf = new ByteArrayOutputStream();
        this.toString(new PrintStream(buf), minimize, "");
        return buf.toString();
    }

    public String toString() {
        return this.toString(false);
    }
}
