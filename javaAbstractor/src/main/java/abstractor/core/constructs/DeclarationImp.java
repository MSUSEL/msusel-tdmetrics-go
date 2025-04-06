package abstractor.core.constructs;

import spoon.reflect.declaration.CtModifiable;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public abstract class DeclarationImp extends ConstructImp implements Declaration {
    public final PackageCon pkg;
    public final Location loc;
    public final String name;
    
    public String visibility;

    public DeclarationImp(PackageCon pkg, Location loc, String name) {
        this.pkg = pkg;
        this.loc = loc;
        this.name = name;
        this.visibility = "";
    }

    public void setVisibility(CtModifiable mod) {
        this.visibility = mod.getVisibility() == null ? "" : mod.getVisibility().toString();
    }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        if (this.pkg != null) obj.put("package", index(pkg));
        if (this.loc != null) obj.putNotEmpty("loc", this.loc.toJson(h));
        obj.putNotEmpty("name", this.name);
        obj.putNotEmpty("vis", this.visibility);
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.name, () -> ((DeclarationImp)c).name),
            Cmp.defer(this.pkg, () -> ((DeclarationImp)c).pkg)
        );
    }
}
