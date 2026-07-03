package abstractor.core.constructs;

import java.util.*;

import abstractor.core.cmp.*;
import abstractor.core.iter.*;
import abstractor.core.json.*;

public class Selection extends ConstructImp {
    public String name;
    public Ref<? extends Construct> origin;

    public Selection() {}

    public Selection(String name, Ref<? extends Construct> origin) {
        this.name   = name;
        this.origin = origin;
    }

    public ConstructKind kind() { return ConstructKind.SELECTION; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("name",   this.name);
        obj.put("origin", key(this.origin));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.defer(this.name,   () -> ((Selection)c).name),
            Cmp.defer(this.origin, () -> ((Selection)c).origin)
        );
    }

    @Override
    public Iterator<Construct> subConstructs() {
        return Iter.SingleIterator(this.origin);
    } 
}
