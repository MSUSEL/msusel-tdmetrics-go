package abstractor.core.constructs;

import java.util.*;

import abstractor.core.cmp.*;
import abstractor.core.iter.*;
import abstractor.core.json.*;

public class StructDesc extends ConstructImp implements TypeDesc {
    public final ArrayList<Ref<Field>> fields = new ArrayList<>();

    // pin is mostly used for shadows to keep them different
    // while shadow objects have fields added on demand.
    public Ref<? extends Construct> pin;

    public StructDesc(List<Ref<Field>> fields) {
        if (fields != null) this.fields.addAll(fields);
    }

    public StructDesc(List<Ref<Field>> fields, Ref<? extends Construct> pin) {
        if (fields != null) this.fields.addAll(fields);
        this.pin = pin;
    }

    public ConstructKind kind() { return ConstructKind.STRUCT_DESC; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("fields", indexList(this.fields));
        obj.putNotEmpty("pin",    key(this.pin));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.deferList(this.fields, () -> ((StructDesc)c).fields),
            Cmp.defer(    this.pin,    () -> ((StructDesc)c).pin)
        );
    }

    @Override
    public Iterator<Ref<? extends Construct>> subConstructs() {
        return new Bundle<Ref<? extends Construct>>(this.fields);
    }
}
