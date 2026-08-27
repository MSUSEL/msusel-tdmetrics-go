package abstractor.core.constructs;

import java.util.*;

import abstractor.core.cmp.*;
import abstractor.core.iter.*;
import abstractor.core.json.*;

public class TypeParam extends ConstructImp implements TypeDesc {
    public String                  name;
    public Ref<? extends TypeDesc> type;
    
    public TypeParam() {}

    public TypeParam(String name, Ref<? extends TypeDesc> type) {
        this.name = name;
        this.type = type;
    }

    public ConstructKind kind() { return ConstructKind.TYPE_PARAM; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("name", this.name);
        obj.put("type", key(this.type));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.defer(this.name, () -> ((TypeParam)c).name),
            Cmp.defer(this.type, () -> ((TypeParam)c).type)
        );
    }

    @Override
    public Iterator<Ref<? extends Construct>> subConstructs() {
        return Iter.SingleIterator(this.type);
    }
}
