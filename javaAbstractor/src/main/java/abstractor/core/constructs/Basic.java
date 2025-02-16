package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtField;

public class Basic extends TypeDesc {
    private final CtField<?> src;
    private final String name;

    static public Basic Create(Project proj, CtField<?> src) {
        Basic existing = proj.basics.findWithSource(src);
        if (existing != null) return existing;

        String name = "TODO"; // TODO: Get name

        Basic f = new Basic(src, name);
        existing = proj.basics.tryAdd(f);
        if (existing != null) return existing;
        return f;
    }

    private Basic(CtField<?> src, String name) {
        this.src = src;
        this.name = name;
    }

    public Object source() { return this.src; }
    public String kind() { return "basic"; }

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
