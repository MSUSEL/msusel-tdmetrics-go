package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Value extends DeclarationImp {
    public boolean  constant;
    public Metrics  metrics;
    public TypeDesc type;

    public Value() {}

    public Value(PackageCon pkg, Location loc,
        String name, boolean constant, Metrics metrics, TypeDesc type) {
        super(pkg, loc, name);
        this.constant = constant;
        this.metrics  = metrics;
        this.type     = type;
    }

    public ConstructKind kind() { return ConstructKind.VALUE; }

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
