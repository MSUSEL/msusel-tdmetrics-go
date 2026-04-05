package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class TypeParam extends ConstructImp implements TypeDesc {
    public String        name;
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
    public Cmp getCmp(Construct c) {
        return Cmp.or(super.getCmp(c),
            Cmp.defer(this.name, () -> ((TypeParam)c).name),
            Cmp.defer(this.type, () -> ((TypeParam)c).type)
        );
    }
}
