package abstractor.core.constructs;

import spoon.reflect.declaration.CtModifiable;

import abstractor.core.cmp.Cmp;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;
import abstractor.core.require.Require;

public abstract class DeclarationImp extends ConstructImp implements Declaration {
    public Ref<PackageCon> pkg;
    public Location        loc;
    public String          name;
    public String          visibility;
    public boolean         isStatic;

    public DeclarationImp(Ref<PackageCon> pkg, Location loc, String name) throws Exception {
        Require.notNull(pkg, "a declaration may not have a null package");
        Require.notBlank(name, "a declaration may not have a blank name");
        this.pkg        = pkg;
        this.loc        = loc;
        this.name       = name;
        this.visibility = "";
        this.isStatic   = false;
    }

    public void setVisibility(CtModifiable mod) {
        this.visibility = mod.getVisibility() == null ? "" : mod.getVisibility().toString();
    }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        if (this.pkg != null) obj.put("package", index(pkg));
        if (this.loc != null) obj.putNotEmpty("loc", this.loc.toJson(h));
        obj.putNotEmpty("name",   this.name);
        obj.putNotEmpty("vis",    this.visibility);
        obj.putNotEmpty("static", this.isStatic);
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.defer(this.name, () -> ((DeclarationImp)c).name),
            Cmp.defer(this.pkg,  () -> ((DeclarationImp)c).pkg)
        );
    }
}
