package abstractor.core.constructs;

import spoon.reflect.declaration.CtModifiable;

import java.util.TreeSet;

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

    public Ref<? extends Construct> nest;
    public final TreeSet<Ref<? extends TypeDesc>> nestedTypes = new TreeSet<>();

    public DeclarationImp(Ref<PackageCon> pkg, Location loc, String name) throws Exception {
        Require.notNull(pkg, "a declaration may not have a null package");
        Require.notBlank(name, "a declaration may not have a blank name");
        this.pkg        = pkg;
        this.loc        = loc;
        this.name       = name;
        this.visibility = "";
        this.isStatic   = false;
    }

    public void setNest(Ref<? extends Construct> nest) throws Exception {
        if (nest == null) return;
        Require.isNull(this.nest);
        this.nest = nest;
    }

    public void setVisibility(CtModifiable mod) {
        this.visibility = mod.getVisibility() == null ? "" : mod.getVisibility().toString();
    }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        if (this.loc != null) obj.putNotEmpty("loc", this.loc.toJson(h));
        obj.putNotEmpty("package", index(pkg));
        obj.putNotEmpty("name",    this.name);
        obj.putNotEmpty("vis",     this.visibility);
        obj.putNotEmpty("static",  this.isStatic);
        obj.putNotEmpty("nest",    key(this.nest));
        obj.putNotEmpty("nested",  keySet(this.nestedTypes));
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
