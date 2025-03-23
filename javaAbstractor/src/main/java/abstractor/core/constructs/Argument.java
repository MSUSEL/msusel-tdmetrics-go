package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Argument extends ConstructImp {
    public final String name;
    public final TypeDesc type;

    public Argument(String name, TypeDesc type) {
        this.name = name;
        this.type = type;
    }
    
    public Argument(TypeDesc type) {
        this("", type);
    }

    public ConstructKind kind() { return ConstructKind.ARGUMENT; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("name", this.name);
        obj.put("type", key(this.type));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.name, () -> ((Argument)c).name),
            Cmp.defer(this.type, () -> ((Argument)c).type)
        );
    }   
}
