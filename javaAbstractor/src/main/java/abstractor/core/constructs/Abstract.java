package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Abstract extends ConstructImp {
    public final String name;
    public final Signature signature;

    public Abstract(String name, Signature signature) {
        this.name = name;
        this.signature = signature;
    }

    public ConstructKind kind() { return ConstructKind.ABSTRACT; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("name", this.name);
        obj.put("signature", index(this.signature));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.name,      () -> ((Abstract)c).name),
            Cmp.defer(this.signature, () -> ((Abstract)c).signature)
        );
    }   
}
