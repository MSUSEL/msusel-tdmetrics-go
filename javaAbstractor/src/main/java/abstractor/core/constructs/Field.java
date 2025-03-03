package abstractor.core.constructs;

import spoon.reflect.declaration.CtModifiable;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Field extends ConstructImp {
    public final String name;
    public final TypeDesc type;
    public String visibility;

    public Field(String name, TypeDesc type) {
        this.name = name;
        this.type = type;
        this.visibility = "";
    }

    public void setVisibility(CtModifiable mod) {
        this.visibility = mod.getVisibility() == null ? "" : mod.getVisibility().toString();
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
