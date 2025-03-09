package abstractor.core.json;

import java.util.ArrayList;

public class JsonArray extends ArrayList<JsonNode> implements JsonNode {
    
    public JsonArray(JsonNode ...nodes) {
        for(JsonNode node: nodes) this.add(node);
    }

    public JsonArray(Iterable<JsonNode> nodes) {
        for (JsonNode node: nodes) this.add(node);
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

    public String toString() { return JsonFormat.Normal().format(this); }
}
