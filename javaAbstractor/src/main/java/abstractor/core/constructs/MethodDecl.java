package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtMethod;

public class MethodDecl extends Declaration {
    private final CtMethod<?> src;
    public final ObjectDecl receiver;

    public MethodDecl(CtMethod<?> src, PackageCon pkg, ObjectDecl receiver, Locations.Location loc, String name) {
        super(pkg, loc, name);
        this.receiver = receiver;
        this.src = src;
    }

    public Object source() { return this.src; }
    public String kind() { return "method"; }
    
    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        // TODO: | `instances`  | ⬤ | List of [indices](#indices) to [method instances](#method-instance). |
        // TODO: | `metrics`    | ⬤ | The [index](#indices) of the [metrics](#metrics) for this method. |
        obj.putNotEmpty("receiver", index(this.receiver));
        // TODO: | `signature`  | ◯ | The [index](#indices) of the [signature](#signature) for this method. |
        // TODO: | `typeParams` | ⬤ | List of [indices](#indices) to [type parameters](#type-parameter) if this method is generic. | 
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.receiver, () -> ((MethodDecl)c).receiver)
            // TODO: | `signature`  | ◯ | The [index](#indices) of the [signature](#signature) for this method. |
            // TODO: | `typeParams` | ⬤ | List of [indices](#indices) to [type parameters](#type-parameter) if this object is generic. |
        );
    }
}
