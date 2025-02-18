package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtInterface;

public class InterfaceDecl extends Declaration implements TypeDesc {
    private final CtInterface<?> src;

    public InterfaceDecl(CtInterface<?> src, PackageCon pkg, Locations.Location loc, String name) {
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
