package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtClass;

public class ObjectDecl extends Declaration {
    private final CtClass<?> src;


    //| `data`       | ◯ | ◯ | The [index](#indices) of the [structure description](#structure-description). |
    //| `exported`   | ◯ | ⬤ | True if the scope is "exported". |
    //| `instances`  | ⬤ | ◯ | List of [indices](#indices) to [object instances](#object-instance). |
    //| `methods`    | ⬤ | ◯ | List of [indices](#indices) to [methods](#method) that have this object as a receiver. |
    //| `typeParams` | ⬤ | ◯ | List of [indices](#indices) to [type parameters](#type-parameter) if this object is generic. |
    //| `interface`  | ◯ | ◯ | The [index](#indices) to the [interface description](#interface-description) that this object matches with. |


    public ObjectDecl(Project proj, CtClass<?> src, Package pkg) {
        super(pkg, new Location(src.getPosition()), src.getSimpleName());
        this.src = src;
    }

    public Object source() { return this.src; }
    public String kind() { return "object"; }
    
    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);

        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c)
        );
    }
}
