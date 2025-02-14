package abstractor.core.constructs;

import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtClass;

public class ObjectDecl extends Declaration {
    private final CtClass<?> src;

    public final TreeSet<MethodDecl> methodDecls = new TreeSet<MethodDecl>();

    static public ObjectDecl Create(Project proj, CtClass<?> src) {
        ObjectDecl existing = proj.objectDecls.findWithSource(src);
        if (existing != null) return existing;

        final Locations.Location loc = proj.locations.create(src.getPosition());
        final Package pkg = Package.Create(proj, src.getPackage());
        final String name = src.getSimpleName();

        // TODO: Handle enum?
        //if (c instanceof CtEnum<?> e) {}

        ObjectDecl od = new ObjectDecl(src, pkg, loc, name);
        existing = proj.objectDecls.tryAdd(od);
        if (existing != null) return existing;
        
        pkg.objectDecls.add(od);
        return od;
    }

    private ObjectDecl(CtClass<?> src, Package pkg, Locations.Location loc, String name) {
        super(pkg, loc, name);
        this.src = src;
    }

    public Object source() { return this.src; }
    public String kind() { return "object"; }
    
    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        // TODO: | `data`       | ◯ | The [index](#indices) of the [structure description](#structure-description). |
        // TODO: | `instances`  | ⬤ | List of [indices](#indices) to [object instances](#object-instance). |
        obj.putNotEmpty("methods", Construct.indexList(this.methodDecls));
        // TODO: | `typeParams` | ⬤ | List of [indices](#indices) to [type parameters](#type-parameter) if this object is generic. |
        // TODO: | `interface`  | ◯ | The [index](#indices) to the [interface description](#interface-description) that this object matches with. 
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c)
            // TODO: | `data`       | ◯ | The [index](#indices) of the [structure description](#structure-description). |
            // TODO: | `typeParams` | ⬤ | List of [indices](#indices) to [type parameters](#type-parameter) if this object is generic. |
        );
    }
}
