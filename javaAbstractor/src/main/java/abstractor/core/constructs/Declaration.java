package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public abstract class Declaration extends ConstructImp {
    public final PackageCon pkg;
    public final Locations.Location loc;
    public final String name;

    public Declaration(PackageCon pkg, Locations.Location loc, String name) {
        this.pkg = pkg;
        this.loc = loc;
        this.name = name;
    }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        if (this.pkg != null) obj.put("package", index(pkg));
        obj.putNotEmpty("loc", this.loc.toJson(h));
        obj.putNotEmpty("name", this.name);
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.name, () -> ((Declaration)c).name),
            Cmp.defer(this.pkg, () -> ((Declaration)c).pkg)
        );
    }
}
