package abstractor.core.json;

import java.io.*;
import java.util.ArrayList;

public class JsonArray extends ArrayList<JsonNode> implements JsonNode {
    
    public JsonArray(JsonNode ...nodes) {
        for(JsonNode node: nodes) this.add(node);
    }

    public boolean isSimple() {
        for (JsonNode elem : this) {
            if (!(elem instanceof JsonValue))
                return false;
        }
        return true;
    }

    public void toString(PrintStream sb, boolean minimize, String indent) {
        if (isEmpty()) {
            sb.append("[]");
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

    public String toString() {
        ByteArrayOutputStream buf = new ByteArrayOutputStream();
        this.toString(new PrintStream(buf), false, "");
        return buf.toString();
    }
}
