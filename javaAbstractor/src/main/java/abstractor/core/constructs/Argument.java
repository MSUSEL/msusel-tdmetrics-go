package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtField;

public class Argument extends Construct {
    private final CtField<?> src;

    static public Argument Create(Project proj, CtField<?> src) {
        Argument existing = proj.arguments.findWithSource(src);
        if (existing != null) return existing;

        // TODO: Get initial stuff

        Argument f = new Argument(src);
        existing = proj.arguments.tryAdd(f);
        if (existing != null) return existing;
        
        // TODO: Finish loading 

        return f;
    }

    private Argument(CtField<?> src) {
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
