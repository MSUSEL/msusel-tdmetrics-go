package abstractor.core.constructs;

import java.util.*;

import abstractor.core.cmp.*;
import abstractor.core.json.*;

/**
 * Instantiation key is a placeholder that takes the place of a CtElement
 * in the factory for a construct that is being instantiated, the target.
 */
public class InstKey implements Comparable<InstKey>, CmpGetter<InstKey>, Jsonable {
    public final Construct target;
    public final ArrayList<Construct> typeArgs = new ArrayList<>();

    public InstKey(Construct target, Collection<? extends Construct> rest) {
        this.target = target;
        this.typeArgs.addAll(rest);
    }

    public ConstructKind kind() {
        return this.target.kind();
    }

    public int compareTo(InstKey c) {
        return Cmp.compareTo(this, c, this.target.getCmpOptions());
    }
    
    @Override
    public Cmp getCmp(InstKey g, CmpOptions options) {
        return Cmp.or(
            Cmp.defer(this.target, () -> g.target),
            Cmp.deferList(this.typeArgs, () -> g.typeArgs)
        );
    }

    @Override
    public boolean equals(Object obj) {
        return obj instanceof InstKey c && this.compareTo(c) == 0;
    }

    public JsonNode toJson(JsonHelper h) {
        JsonArray indices = new JsonArray();
        for (Construct o: this.typeArgs) indices.add(o.toJson(h));

        JsonObject obj = new JsonObject();
        obj.put("target",  this.target.toJson(h));
        obj.put("typeArgs", indices);
        return obj;
    }

    @Override
    public String toString() {
        JsonHelper jh = new JsonHelper();
        jh.writeKinds     = true;
        jh.writeIndices   = true;
        jh.writeRefs      = true;
        jh.refSkipResolve = true;
        return JsonFormat.Relaxed().format(this.toJson(jh));
    }
}
