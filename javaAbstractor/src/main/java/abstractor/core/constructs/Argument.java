package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtField;

public class Argument extends ConstructImp {
    private final CtField<?> src;

    public Argument(CtField<?> src) {
        this.src = src;
    }

    public Object source() { return this.src; }
    public String kind() { return "argument"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        // TODO: | `name`  | ⬤ | ◯ | The optional string name for the argument. |
        // TODO: | `type`  | ◯ | ◯ | [Key](#keys) for any [type description](#type-descriptions). |
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c)
            // TODO: | `name`  | ⬤ | ◯ | The optional string name for the argument. |
            // TODO: | `type`  | ◯ | ◯ | [Key](#keys) for any [type description](#type-descriptions). |
        );
    }   
}
