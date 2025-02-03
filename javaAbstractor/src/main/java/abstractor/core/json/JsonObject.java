package abstractor.core.json;

import java.io.*;
import java.util.HashMap;

public class JsonObject extends HashMap<String, JsonNode> implements JsonNode {
    
    public void toString(PrintStream sb, boolean minimize, String indent) {
        if (isEmpty()) {
            sb.append("{}");
            return;
        }

        sb.append("{");
        boolean first = true;
        final String indent2 = indent + "  ";
        for (String key : this.keySet()) {
            if (first) first = false;
            else sb.append(",");
            if (!minimize) sb.append("\n"+indent2);
            sb.append(key);
            sb.append(":");
            if (!minimize) sb.append(" ");
            this.get(key).toString(sb, minimize, indent2);            
        }
        if (!minimize) sb.append("\n"+indent);
        sb.append("}");
    }

    public String toString() {
        ByteArrayOutputStream buf = new ByteArrayOutputStream();
        this.toString(new PrintStream(buf), false, "");
        return buf.toString();
    }
}
