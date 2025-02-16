package abstractor.core.constructs;

import java.util.*;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtClass;

public class StructDesc extends TypeDesc {
    private final CtClass<?> src;

    public final List<Field> fields;

    static public StructDesc Create(Project proj, CtClass<?> src) {
        StructDesc existing = proj.structDescs.findWithSource(src);
        if (existing != null) return existing;

        // TODO: Handle enum?
        //if (c instanceof CtEnum<?> e) {}
        
        ArrayList<Field> fields = new ArrayList<Field>();

        StructDesc sd = new StructDesc(src, fields);
        existing = proj.structDescs.tryAdd(sd);
        if (existing != null) return existing;

        return sd;
    }

    private StructDesc(CtClass<?> src, ArrayList<Field> fields) {
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
