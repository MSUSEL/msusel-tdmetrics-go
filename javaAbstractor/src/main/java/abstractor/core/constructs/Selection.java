package abstractor.core.constructs;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Selection extends ConstructImp {
    public final String name;
    public final Construct origin;

    public Selection(String name, Construct origin) {
        this.name = name;
        this.origin = origin;
    }

    public String kind() { return "selection"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("name", this.name);
        obj.put("origin", key(this.origin));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.name,   () -> ((Selection)c).name),
            Cmp.defer(this.origin, () -> ((Selection)c).origin)
        );
    }   
}
