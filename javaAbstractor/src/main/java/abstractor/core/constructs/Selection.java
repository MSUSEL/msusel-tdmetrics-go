package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtField;

public class Selection extends ConstructImp implements TypeDesc {
    private final CtField<?> src;

    public Selection(CtField<?> src) {
        this.src = src;
    }

    public Object source() { return this.src; }
    public String kind() { return "selection"; }

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
