package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.List;

import abstractor.core.cmp.Cmp;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;

public class Signature extends ConstructImp implements TypeDesc {
    public boolean variadic;
    public final ArrayList<Ref<Argument>> params  = new ArrayList<Ref<Argument>>();
    public final ArrayList<Ref<Argument>> results = new ArrayList<Ref<Argument>>();
    
    public Signature() {}

    public Signature(boolean variadic, List<Ref<Argument>> params, List<Ref<Argument>> results) {
        this.variadic = variadic;
        if (params  != null) this.params .addAll(params);
        if (results != null) this.results.addAll(results);
    }

    public ConstructKind kind() { return ConstructKind.SIGNATURE; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("variadic", this.variadic);
        obj.putNotEmpty("params",   indexList(this.params));
        obj.putNotEmpty("results",  indexList(this.results));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.defer(    this.variadic, () -> ((Signature)c).variadic),
            Cmp.deferList(this.params,   () -> ((Signature)c).params),
            Cmp.deferList(this.results,  () -> ((Signature)c).results)
        );
    }   
}
