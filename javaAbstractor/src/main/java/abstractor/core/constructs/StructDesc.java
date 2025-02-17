package abstractor.core.constructs;

import java.util.*;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtClass;

public class StructDesc extends TypeDesc {
    private final CtClass<?> src;

    public final List<Field> fields;

    public StructDesc(CtClass<?> src, ArrayList<Field> fields) {
        this.src = src;
        this.fields = Collections.unmodifiableList(fields);
    }

    public Object source() { return this.src; }
    public String kind() { return "structDesc"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        // TODO: | `fields` | List of [indices](#indices) of [fields](#field) in this structure. |
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c)
            // TODO: | `fields` | List of [indices](#indices) of [fields](#field) in this structure. |
        );
    }
}
