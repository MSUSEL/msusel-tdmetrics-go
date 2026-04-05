package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Basic extends ConstructImp implements TypeDesc {
    public String name;

    public Basic()            { this.name = "undefined"; }
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
    public Cmp getCmp(Construct c) {
        return Cmp.or(super.getCmp(c),
            Cmp.defer(this.name, () -> ((Basic)c).name)
        );
    }   
}
