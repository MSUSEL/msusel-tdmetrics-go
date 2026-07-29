package abstractor.core.constructs;

import java.util.Iterator;

import abstractor.core.cmp.Cmp;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.iter.Iter;
import abstractor.core.json.*;

public class Argument extends ConstructImp {
    public String                  name;
    public Ref<? extends TypeDesc> type;

    public Argument() {}

    public Argument(String name, Ref<? extends TypeDesc> type) throws Exception {
        this.name = name;
        this.type = type;
        if(this.type == null)
            throw new Exception("Argument type may not be null (name: " + this.name + ")");
    }
    
    public Argument(Ref<? extends TypeDesc> type) throws Exception {
        this("", type);
    }

    public ConstructKind kind() { return ConstructKind.ARGUMENT; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("name", this.name);
        obj.put(        "type", key(this.type));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or("Argument", super.getCmp(c, options),
            Cmp.defer(this.name, () -> ((Argument)c).name, "name"),
            Cmp.defer(this.type, () -> ((Argument)c).type, "type")
        );
    }

    @Override
    public Iterator<Ref<? extends Construct>> subConstructs() {
        return Iter.SingleIterator(this.type);
    }
}
