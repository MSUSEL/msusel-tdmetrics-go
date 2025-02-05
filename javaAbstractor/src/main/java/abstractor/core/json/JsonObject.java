package abstractor.core.json;

import java.io.*;
import java.util.HashMap;

public class JsonObject extends HashMap<String, JsonNode> implements JsonNode {
    
    public boolean isSimple() {
        for (JsonNode elem : this.values()) {
            if (!(elem instanceof JsonValue))
                return false;
        }
        return true;
    }

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
        for (String key : this.keySet()) {
            if (first) first = false;
            else sb.append(",");
            if (!minimize) sb.append(separator);

            sb.append(JsonValue.escape(key));
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

    public String toString() {
        return this.toString(false);
    }
}
