package abstractor.core.constructs;

import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtClass;

public class ObjectDecl extends Declaration implements TypeDesc {
    private final CtClass<?> src;
    public final StructDesc struct;
    public final TreeSet<MethodDecl> methodDecls = new TreeSet<MethodDecl>();

    public ObjectDecl(CtClass<?> src, StructDesc struct, PackageCon pkg, Locations.Location loc, String name) {
        super(pkg, loc, name);
        this.src = src;
        this.struct = struct;
    }

    public Object source() { return this.src; }
    public String kind() { return "object"; }
    
    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("data", this.struct.getIndex());
        // TODO: | `instances`  | ⬤ | List of [indices](#indices) to [object instances](#object-instance). |
        obj.putNotEmpty("methods", indexList(this.methodDecls));
        // TODO: | `typeParams` | ⬤ | List of [indices](#indices) to [type parameters](#type-parameter) if this object is generic. |
        // TODO: | `interface`  | ◯ | The [index](#indices) to the [interface description](#interface-description) that this object matches with. 
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.struct, () -> ((ObjectDecl)c).struct)
            // TODO: | `typeParams` | ⬤ | List of [indices](#indices) to [type parameters](#type-parameter) if this object is generic. |
        );
    }
}
