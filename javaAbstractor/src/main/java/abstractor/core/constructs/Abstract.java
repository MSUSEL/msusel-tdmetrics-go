package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtField;

public class Abstract extends ConstructImp {
    private final CtField<?> src;

    public Abstract(CtField<?> src) {
        this.src = src;
    }

    public Object source() { return this.src; }
    public String kind() { return "abstract"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        // TODO: | `name`      | ◯ | ◯ | The string name for the abstract. |
        // TODO: | `signature` | ◯ | ◯ | [Index](#indices) for the [signature](#signature). |
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c)
            
            // TODO: | `name`      | ◯ | ◯ | The string name for the abstract. |
            // TODO: | `signature` | ◯ | ◯ | [Index](#indices) for the [signature](#signature). |
        );
    }   
}
