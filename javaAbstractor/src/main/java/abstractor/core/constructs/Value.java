package abstractor.core.constructs;

import spoon.reflect.declaration.CtInterface;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Value extends Declaration {
    public final boolean constant;
    public final Metrics metrics;
    public final TypeDesc type;

    public Value(CtInterface<?> src, PackageCon pkg, Location loc,
        String name, boolean constant, Metrics metrics, TypeDesc type) {
        super(src, pkg, loc, name);
        this.constant = constant;
        this.metrics  = metrics;
        this.type     = type;
    }

    public String kind() { return "value"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("const",   this.constant);
        obj.putNotEmpty("metrics", index(this.metrics));
        obj.put("type",            key(this.type));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.constant, () -> ((Value)c).constant),
            Cmp.defer(this.metrics,  () -> ((Value)c).metrics),
            Cmp.defer(this.type,     () -> ((Value)c).type)
        );
    }   
}
