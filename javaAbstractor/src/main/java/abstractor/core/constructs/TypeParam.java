package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtField;

public class TypeParam extends ConstructImp implements TypeDesc {
    private final CtField<?> src;

    public TypeParam(CtField<?> src) {
        this.src = src;
    }

    public Object source() { return this.src; }
    public String kind() { return "typeParam"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        // TODO: Fill out
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c)
            // TODO: Fill out
        );
    }   
}
