package abstractor.core.json;

import java.io.PrintStream;

public class JsonValue implements JsonNode {

    public static String escape(String text) {
        StringBuilder sb = new StringBuilder();
        sb.append('"');
        for (char c : text.toCharArray()) {
            switch (c) {
                case '\\': sb.append("\\\\"); break;
                case '"':  sb.append("\\\""); break;
                default:   sb.append(c);      break;
            }
        }
        sb.append('"');
        return sb.toString();
    }

    public static JsonValue of(String value) { return new JsonValue(value); }

    public static JsonValue of(int value) { return new JsonValue(value); }

    public static JsonValue of(double value) { return new JsonValue(value); }

    public static JsonValue of(boolean value) { return new JsonValue(value); }
    
    private static final JsonValue nullSingleton = new JsonValue(null);

    public static JsonValue ofNull() { return nullSingleton; }

    private final Object value;

    private JsonValue(Object value) { this.value = value; }

    public boolean isString() { return this.value instanceof String; }
    
    public boolean isInt() { return this.value instanceof Integer; }
    
    public boolean isDouble() { return this.value instanceof Double; }
    
    public boolean isBool() { return this.value instanceof Boolean; }
    
    public boolean isNull() { return this.value == null; }

    public boolean isEmpty() {
        if (this.isBool()) return !(boolean)this.value;
        if (this.isInt()) return (int)this.value == 0;
        if (this.isDouble()) return (double)this.value == 0.0;
        if (this.isString()) return ((String)this.value).isEmpty();
        return true;
    }

    public String asString() {
        return this.isNull()? "null": this.value.toString();
    }
    
    public int asInt() {
        if (this.isInt()) return (int)this.value;
        if (this.isDouble()) return (int)((double)this.value);
        if (this.isBool()) return (boolean)this.value ? 1 : 0;
        if (this.isString()) {
            try {
                return Integer.parseInt((String)this.value);
            } catch(Exception ex) {}
            try {
                return (int)Double.parseDouble((String)this.value);
            } catch(Exception ex) {}
            if (Boolean.parseBoolean((String)this.value)) return 1;
        }
        return 0;
    }
    
    public double asDouble() {
        if (this.isDouble()) return (double)this.value;
        if (this.isInt()) return (double)((int)this.value);
        if (this.isBool()) return (boolean)this.value ? 1.0 : 0.0;
        if (this.isString()) {
            try {
                return Double.parseDouble((String)this.value);
            } catch(Exception ex) {}
            if (Boolean.parseBoolean((String)this.value)) return 1.0;
        }
        return 0.0;
    }
    
    public boolean asBool() {
        if (this.isBool()) return (boolean)this.value;
        if (this.isInt()) return (int)this.value != 0;
        if (this.isDouble()) return (double)this.value != 0.0;
        if (this.isString()) {
            if (Boolean.parseBoolean((String)this.value)) return true;
            try {
                return Double.parseDouble((String)this.value) != 0.0;
            } catch(Exception ex) {}
        }
        return false;
    }

    public String toString() {
        if (this.isNull()) return "null";
        if (this.isString()) return escape(this.value.toString());
        return this.value.toString();
    }
    
    public String toString(boolean minimize) {
        return this.toString();
    }
    
    public void toString(PrintStream sb, boolean minimize, String indent) {
        sb.append(this.toString());
    }
}
