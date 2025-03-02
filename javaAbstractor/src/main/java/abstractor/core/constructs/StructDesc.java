package abstractor.core.constructs;

import java.util.List;

import spoon.reflect.declaration.CtClass;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class StructDesc extends ConstructImp implements TypeDesc {
    public final List<Field> fields;

    public StructDesc(CtClass<?> src, List<Field> fields) {
        super(src);
        this.fields = unmodifiableList(fields);
    }

    public String kind() { return "structDesc"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("fields", indexList(this.fields));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.deferList(this.fields, () -> ((StructDesc)c).fields)
        );
    }
}
