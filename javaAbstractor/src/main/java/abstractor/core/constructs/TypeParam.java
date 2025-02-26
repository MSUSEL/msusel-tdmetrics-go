package abstractor.core.constructs;

import spoon.reflect.declaration.CtTypeParameter;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class TypeParam extends ConstructImp implements TypeDesc {
    public final String name;
    public final TypeDesc type;
    
    public TypeParam(CtTypeParameter src, String name, TypeDesc type) {
        super(src);
        this.name = name;
        this.type = type;
    }

    public String kind() { return "typeParam"; }

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
            Cmp.defer(this.name, () -> ((TypeParam)c).name),
            Cmp.defer(this.type, () -> ((TypeParam)c).type)
        );
    }   
}
