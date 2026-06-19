package abstractor.core.json;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

public class JsonFormat {
    static public JsonFormat Minimize() { return new JsonFormat(Style.minimize); }
    static public JsonFormat Expand()   { return new JsonFormat(Style.expand); }
    static public JsonFormat Relaxed()  { return new JsonFormat(Style.relaxed); }
    static public JsonFormat Inline()   { return new JsonFormat(Style.inline); }
    static public JsonFormat Normal()   { return new JsonFormat(Style.normal); }

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

    static final String constEmpty   = "";
    static final String constSingle  = " ";
    static final String constIndent  = "  "; 
    static final String constNewLine = "\n";

    private enum Style {
        minimize, expand, relaxed, inline, normal;

        public String emptySep() {
            return this == minimize ? constEmpty : constSingle;
        }

        public boolean canInline(boolean isSimple, int depth) {
            if (this == minimize) return true;
            if (this == inline) return true;
            if (this == expand) return false;
            if (this == relaxed && depth == 0) return false;
            return isSimple;
        }

        public String firstSep(boolean canInline, String indent2) {
            if (this == minimize) return constEmpty;
            return canInline ? constSingle : constNewLine + indent2;
        }

        public String restSep(boolean canInline, String indent2) {
            return "," + (this == minimize || canInline ? constEmpty : constNewLine + indent2);
        }

        public String tailSep(boolean inline, String indent) {
            if (this == minimize) return constEmpty;
            if (inline) return constSingle;
            if (this == relaxed) return "," + constNewLine + indent;
            return constNewLine + indent;
        }

        public String keyValSep() {
            return ":" + (this == minimize ? constEmpty : constSingle);
        }
        
        public String getString(String key) {
            if (this != relaxed && this != inline) return escape(key);
            if (needsEscape(key)) return escape(key);
            return key;
        }
    }

    private final Style style;

    private JsonFormat(Style style) { this.style = style; }
    
    public String format(JsonNode j) {
        return this.format(j, constEmpty);
    }
    
    public String format(JsonNode j, String indent) {
        ByteArrayOutputStream buf = new ByteArrayOutputStream();
        this.format(new PrintStream(buf), j, indent);
        return buf.toString();
    }
    
    public void format(PrintStream sb, JsonNode j) {
        this.format(sb, j, constEmpty);
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

    private void format(PrintStream sb, JsonValue j, String indent, int depth) {
        sb.append(this.style.getString(j.asString()));
    }

    private void format(PrintStream sb, JsonArray j, String indent, int depth) {
        if (j.isEmpty()) {
            sb.append("[" + this.style.emptySep() + "]");
            return;
        }

        final boolean inline   = this.style.canInline(j.isSimple(), depth);
        final String  indent2  = indent + constIndent;
        final String  firstSep = this.style.firstSep(inline, indent2);
        final String  restSep  = this.style.restSep( inline, indent2);
        final String  tailSep  = this.style.tailSep( inline, indent);
        
        sb.append("[");
        boolean first = true;
        for (JsonNode elem : j) {
            sb.append(first ? firstSep : restSep);
            first = false;
            this.format(sb, elem, indent2, depth + 1);
        }
        sb.append(tailSep);
        sb.append("]");
    }

    private void format(PrintStream sb, JsonObject j, String indent, int depth) {
        if (j.isEmpty()) {
            sb.append("{" + this.style.emptySep() + "}");
            return;
        }

        final boolean inline    = this.style.canInline(j.isSimple(), depth);
        final String  indent2   = indent + constIndent;
        final String  firstSep  = this.style.firstSep(inline, indent2);
        final String  restSep   = this.style.restSep( inline, indent2);
        final String  tailSep   = this.style.tailSep( inline, indent);
        final String  keyValSep = this.style.keyValSep();
        
        sb.append("{");
        boolean first = true;
        for (String key : j.keySet()) {
            sb.append(first ? firstSep : restSep);
            first = false;
            sb.append(this.style.getString(key));
            sb.append(keyValSep);
            this.format(sb, j.get(key), indent2, depth + 1);
        }
        sb.append(tailSep);
        sb.append("}");
    }
}
