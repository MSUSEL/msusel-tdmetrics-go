package abstractor.core.json;

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
    
    public String toString() { return JsonFormat.Normal().format(this); }
}
