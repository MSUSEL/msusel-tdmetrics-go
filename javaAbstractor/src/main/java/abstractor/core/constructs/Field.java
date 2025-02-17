package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtField;

public class Field extends Construct {
    private final CtField<?> src;

    public Field(CtField<?> src) {
        this.src = src;
    }

    public Object source() { return this.src; }
    public String kind() { return "field"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);

        // TODO: | `name`     | ◯ | ◯ | The string name for the field. |
        // TODO: | `type`     | ◯ | ◯ | [Key](#keys) for any [type description](#type-descriptions). |        
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c)
            
            // TODO: | `name`     | ◯ | ◯ | The string name for the field. |
            // TODO: | `type`     | ◯ | ◯ | [Key](#keys) for any [type description](#type-descriptions). |
        );
    }   
}
