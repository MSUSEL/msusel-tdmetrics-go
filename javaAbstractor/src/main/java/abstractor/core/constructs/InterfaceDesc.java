package abstractor.core.constructs;

import spoon.reflect.declaration.CtField;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class InterfaceDesc extends ConstructImp implements TypeDesc {

    public InterfaceDesc(CtField<?> src) {
        super(src);
    }

    public String kind() { return "interfaceDesc"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);

        // TODO: | `abstracts` | ⬤ | ◯ | List of [indices](#indices) to [abstracts](#abstract). |
        // TODO: | `approx`    | ⬤ | ◯ | List of [keys](#keys) to any [type description](#type-descriptions) for approximate constraints. |
        // TODO: | `exact`     | ⬤ | ◯ | List of [keys](#keys) to any [type description](#type-descriptions) for exact constraints. |
        // TODO: | `hint`      | ◯ | ⬤ | A string indicating if the interface is a stand-in for a type, e.g. `pointer`, `chan`, `list` |
        // TODO: | `index`     | ◯ | ⬤ | The [index](#indices) of this interface in the projects' `interfaceDescs` list. |
        // TODO: | `inherits`  | ⬤ | ◯ | List of [indices](#indices) to inherited [interfaces](#interface-description). |
        // TODO: | `kind`      | ◯ | ⬤ | `interfaceDesc` |
        // TODO: | `package`   | ⬤ | ◯ | The [index](#indices) to the [package](#package) this interface is pinned to. |

        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c)
            // TODO: Fill out
        );
    }   
}
