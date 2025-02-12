package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtInterface;

public class InterfaceDecl extends Declaration {
    private final CtInterface<?> src;

    static public InterfaceDecl Create(Project proj, CtInterface<?> src) {
        InterfaceDecl existing = proj.interfaceDecls.findWithSource(src);
        if (existing != null) return existing;

        final Locations.Location loc = proj.locations.create(src.getPosition());
        final Package pkg = Package.Create(proj, src.getPackage());
        final String name = src.getSimpleName();

        InterfaceDecl id = proj.interfaceDecls.tryAdd(new InterfaceDecl(src, pkg, loc, name));
        pkg.interfaceDecls.add(id);
        return id;
    }

    private InterfaceDecl(CtInterface<?> src, Package pkg, Locations.Location loc, String name) {
        super(pkg, loc, name);
        this.src = src;
    }

    public Object source() { return this.src; }
    public String kind() { return "interfaceDecl"; }
    
    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        // TODO: | `instances`  | ⬤ | List of [indices](#indices) to [interface instances](#interface-instance). |
        // TODO: | `interface`  | ◯ | The [index](#indices) to the declared [interface](#interface-description) type. |
        // TODO: | `typeParams` | ⬤ | List of [indices](#indices) to [type parameters](#type-parameter) if this interface is generic. |
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c)
            // TODO: | `interface`  | ◯ | The [index](#indices) to the declared [interface](#interface-description) type. |
            // TODO: | `typeParams` | ⬤ | List of [indices](#indices) to [type parameters](#type-parameter) if this object is generic. |
        );
    }
}
