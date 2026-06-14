package abstractor.core.constructs;

import abstractor.core.AbstractorException;
import abstractor.core.ElementKey;
import abstractor.core.cmp.*;
import abstractor.core.json.*;
import abstractor.core.require.Require;

public class Ref<T extends Construct> extends ConstructImp {
    private final ConstructKind conKind;
    public  final ElementKey    elemKey;
    public  final String        context;

    private T res;

    public Ref(ConstructKind kind, ElementKey elemKey, String context) throws Exception {
        Require.notBlank(context, "may not have a blank reference context.");
        this.conKind = kind;
        this.elemKey = elemKey;
        this.context = context;
    }

    public ConstructKind kind() { return this.conKind; }

    @Override
    public int getIndex() {
        return this.isResolved() ? this.res.getIndex() : super.getIndex();
    }

    public T getResolved() { return this.res; }
    public boolean isResolved() { return this.res != null; }

    public T mustGetResolved() throws Exception {
        if (this.res == null)
            throw new AbstractorException("Expected the " + this.kind() + " reference is resolved for " + this);
        return this.res;
    }

    public void setResolved(T res) throws Exception {
        if (res == null)
            throw new Exception("Attempted to write null as the resolved construct to the reference " + this);
        if (!res.kind().equals(this.conKind))
            throw new Exception("Attempted to write a resolved construct with the kind " + res.kind() + " for reference " + this + " with kind " + this.conKind);
        if (this.isResolved() && !this.res.equals(res)) {
            throw new Exception("Attempted to overwrite the resolved construct, " + this.res + ", with " + res + " for reference " + this);
        }
        this.res = res;
    }

    public JsonNode refToJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("ref",     true);
        obj.put("context", this.context);

        final boolean showExtras = false;
        if (showExtras) {
            obj.put("refHash", this.hashCode());
            obj.put("cmpOptions", String.valueOf(this.getCmpOptions()));
            if (this.isResolved())
                obj.put("resHash", this.res.hashCode());
            if (this.elemKey != null)
                obj.put("elemKey", this.elemKey.toString());
        }
        return obj;
    }

    @Override
    public JsonNode toJson(JsonHelper h) {
        return (!h.refSkipResolve && this.isResolved()) ? this.res.toJson(h) : this.refToJson(h);
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        // Check if both are using the same resolved comparison option since
        // otherwise A.compareTo(B) will not be the negation of B.compareTo(A)
        // because of which CmpOptions are being used.
        Cmp opCmp = Cmp.defer(CmpOptions.shouldUseResolved(this.getCmpOptions()),
            () -> CmpOptions.shouldUseResolved(c.getCmpOptions()), "useResolved");

        if (options.useResolved) {
            return Cmp.or("Ref", super.getCmp(c, options), opCmp,
                Cmp.defer(this.res, () -> ((Ref<?>)c).res, "resolved")
            );
        }

        return Cmp.or("Ref", super.getCmp(c, options), opCmp,
            Cmp.defer(this.elemKey, () -> ((Ref<?>)c).elemKey, "elemKey"),
            Cmp.defer(this.context, () -> ((Ref<?>)c).context, "context")
        );
    }
}
