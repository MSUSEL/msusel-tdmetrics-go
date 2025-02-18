package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtField;

public class Field extends ConstructImp {
    private final CtField<?> src;
    public final String name;
    public final TypeDesc type;

    public Field(CtField<?> src, String name, TypeDesc type) {
        this.src = src;
        this.name = name;
        this.type = type;
    }

    public Object source() { return this.src; }
    public String kind() { return "field"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("name", this.name);
        obj.put("type", key(this.type));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.name, () -> ((Field)c).name),
            Cmp.defer(this.type, () -> ((Field)c).type)
        );
    }   
}
