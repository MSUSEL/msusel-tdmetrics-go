package abstractor.core.constructs;

import spoon.reflect.declaration.CtField;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Field extends ConstructImp {
    public final String name;
    public final TypeDesc type;
    public final String visibility;

    public Field(CtField<?> src, String name, TypeDesc type) {
        super(src);
        this.name = name;
        this.type = type;
        this.visibility = src.getVisibility() == null ? "" : src.getVisibility().toString();
    }

    public String kind() { return "field"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("name", this.name);
        obj.put("type", key(this.type));
        obj.putNotEmpty("visibility", this.visibility);
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
