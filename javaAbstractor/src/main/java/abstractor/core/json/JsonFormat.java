package abstractor.core.json;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

public class JsonFormat {
    static public JsonFormat Minimize() { return new JsonFormat( true, false, false); }
    static public JsonFormat Expand()   { return new JsonFormat(false,  true, false); }
    static public JsonFormat Relaxed()  { return new JsonFormat(false, false,  true); }
    static public JsonFormat Normal()   { return new JsonFormat(false, false, false); }

    static private boolean needsEscape(String text) {
        if (text.equalsIgnoreCase("null") ||
            text.equalsIgnoreCase("true") ||
            text.equalsIgnoreCase("false"))
            return true;

        final char[] cs = text.toCharArray();
        for (int i = 0; i < cs.length; i++) {
            final char c = cs[i];
            if (c >= 'a' && c <= 'z') continue;
            if (c >= 'A' && c <= 'Z') continue;
            if (c == '_' || c == '-' || c == '.' || c == '$') continue;
            if (c >= '0' && c <= '9' && i >= 0) continue;
            return true;
        }
        return false;
    }

    static private String escape(String text) {
        StringBuilder sb = new StringBuilder();
        sb.append('"');
        for (char c : text.toCharArray()) {
            switch (c) {
                case '\\': sb.append("\\\\"); break;
                case '"':  sb.append("\\\""); break;
                case '\b': sb.append("\\b");  break;
                case '\f': sb.append("\\f");  break;
                case '\n': sb.append("\\n");  break;
                case '\r': sb.append("\\r");  break;
                case '\t': sb.append("\\t");  break;
                default:   sb.append(c);      break;
            }
        }
        sb.append('"');
        return sb.toString();
    }

    private final boolean minimize;
    private final boolean expand;
    private final boolean relaxed;

    private JsonFormat(boolean minimize, boolean expand, boolean relaxed) {
        this.minimize = minimize;
        this.expand   = expand;
        this.relaxed  = relaxed;
    }
    
    public String format(JsonNode j) {
        return this.format(j, "");
    }
    
    public String format(JsonNode j, String indent) {
        ByteArrayOutputStream buf = new ByteArrayOutputStream();
        this.format(new PrintStream(buf), j, indent);
        return buf.toString();
    }
    
    public void format(PrintStream sb, JsonNode j) {
        this.format(sb, j, "");
    }
    
    public void format(PrintStream sb, JsonNode j, String indent) {
        this.format(sb, j, indent, 0);
    }
    
    private void format(PrintStream sb, JsonNode j, String indent, int depth) {
        if (j == null) sb.append("null");
        else if (j instanceof JsonValue  jv) this.format(sb, jv, indent, depth);
        else if (j instanceof JsonArray  ja) this.format(sb, ja, indent, depth);
        else if (j instanceof JsonObject jo) this.format(sb, jo, indent, depth);
        else sb.append("unknown_node");
    }

    private boolean simplify(int depth) {
        return !this.expand && (depth != 0 || !this.relaxed);
    }

    private void format(PrintStream sb, JsonValue j, String indent, int depth) {
        String str = j.asString();
        if (j.isString() && (!this.relaxed || needsEscape(str)))
            str = escape(str);
        sb.append(str);
    }

    private void format(PrintStream sb, JsonArray j, String indent, int depth) {
        if (j.isEmpty()) {
            sb.append(this.minimize? "[]": "[ ]");
            return;
        }

        sb.append("[");
        boolean first = true;
        final boolean simple = this.simplify(depth) && j.isSimple();
        final String indent2 = indent + "  ";
        final String separator = simple ? " " : "\n" + indent2;
        for (JsonNode elem : j) {
            if (first) first = false;
            else sb.append(",");
            if (!this.minimize) sb.append(separator);
            
            this.format(sb, elem, indent2, depth + 1);
        }
        if (!this.minimize) {
            if (simple) sb.append(" ");
            else if (this.relaxed) sb.append(",\n" + indent);
            else sb.append("\n" + indent);
        }
        sb.append("]");
    }
    
    private void format(PrintStream sb, JsonObject j, String indent, int depth) {
        if (j.isEmpty()) {
            sb.append(this.minimize? "{}": "{ }");
            return;
        }

        sb.append("{");
        boolean first = true;
        final boolean simple = this.simplify(depth) && j.isSimple();
        final String indent2 = indent + "  ";
        final String separator = simple ? " " : "\n" + indent2;
        for (Object key : j.keySet()) {
            if (first) first = false;
            else sb.append(",");
            if (!this.minimize) sb.append(separator);

            String keyStr = (String)key;
            if (!this.relaxed || needsEscape(keyStr)) {
                keyStr = escape(keyStr);
            }
            sb.append(keyStr);

            sb.append(":");
            if (!this.minimize) sb.append(" ");
            this.format(sb, j.get(key), indent2, depth + 1);         
        }
        if (!this.minimize) {
            if (simple) sb.append(" ");
            else if (this.relaxed) sb.append(",\n" + indent);
            else sb.append("\n" + indent);
        }
        sb.append("}");
    }
}
