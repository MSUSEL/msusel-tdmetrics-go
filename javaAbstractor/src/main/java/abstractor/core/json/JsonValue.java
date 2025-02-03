package abstractor.core.json;

import java.io.PrintStream;

public class JsonValue implements JsonNode {

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

    public String asString() {
        return this.isNull()? "null": this.value.toString();
    }
    
    public int asInt() {
        if (this.isInt()) return (int)this.value;
        if (this.isDouble()) return (int)((double)this.value);
        return 0;
    }
    
    public double asDouble() {
        if (this.isDouble()) return (double)this.value;
        if (this.isInt()) return (int)this.value;
        return 0.0;
    }
    
    public boolean asBool() {
        if (this.isBool()) return (boolean)this.value;
        if (this.isInt()) return (int)this.value > 0;
        if (this.isDouble()) return (double)this.value > 0.0;
        return this.value == "true";
    }

    public String toString() {
        if (this.isNull()) return "null";
        if (this.isString())
            return "\""+this.value.toString()
                .replaceAll("\\", "\\\\").replaceAll("\"", "\\\"")+"\"";
        return this.value.toString();
    }
    
    public void toString(PrintStream sb, boolean minimize, String indent) {
        sb.append(this.toString());
    }
}
