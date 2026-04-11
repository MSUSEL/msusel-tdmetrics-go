package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.List;

import abstractor.core.cmp.Cmp;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;

public class StructDesc extends ConstructImp implements TypeDesc {
    public final ArrayList<Ref<Field>> fields = new ArrayList<Ref<Field>>();

    public StructDesc() {} 

    public StructDesc(List<Ref<Field>> fields) {
        if (fields != null) this.fields.addAll(fields);
    }

    public ConstructKind kind() { return ConstructKind.STRUCT_DESC; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("fields", indexList(this.fields));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.deferList(this.fields, () -> ((StructDesc)c).fields)
        );
    }
}
