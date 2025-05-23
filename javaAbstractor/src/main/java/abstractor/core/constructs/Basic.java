package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Basic extends ConstructImp implements TypeDesc {
    public final String name;

    public Basic(String name) { this.name = name; }

    public ConstructKind kind() { return ConstructKind.BASIC; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        if (obj.size() <= 0) return JsonValue.of(this.name);
        obj.put("name", this.name);
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.name, () -> ((Basic)c).name)
        );
    }   
}
