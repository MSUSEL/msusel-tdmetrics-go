package abstractor.core.constructs;

import java.util.List;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class Signature extends ConstructImp implements TypeDesc {
    public final boolean variadic;
    public final List<Argument> params;
    public final List<Argument> results;
    
    public Signature(boolean variadic, List<Argument> params, List<Argument> results) {
        this.variadic = variadic;
        this.params   = unmodifiableList(params);
        this.results  = unmodifiableList(results);
    }

    public String kind() { return "signature"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("variadic", this.variadic);
        obj.putNotEmpty("params",   indexList(this.params));
        obj.putNotEmpty("results",  indexList(this.results));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.variadic,    () -> ((Signature)c).variadic),
            Cmp.deferList(this.params,  () -> ((Signature)c).params),
            Cmp.deferList(this.results, () -> ((Signature)c).results)
        );
    }   
}
