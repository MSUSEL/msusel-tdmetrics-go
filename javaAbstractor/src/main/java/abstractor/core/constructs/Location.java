package abstractor.core.constructs;

import abstractor.core.json.*;
import spoon.reflect.cu.SourcePosition;

public class Location implements Jsonable {
    public final SourcePosition pos;
    public int offset;

    public Location(SourcePosition pos) {
        this.pos = pos;
    }

    public JsonNode toJson(JsonHelper h) {
        return JsonValue.of(this.offset);
    }
}
