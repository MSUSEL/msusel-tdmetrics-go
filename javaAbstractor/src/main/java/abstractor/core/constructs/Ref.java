package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.List;

import spoon.reflect.declaration.CtElement;

import abstractor.core.cmp.Cmp;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.JsonHelper;
import abstractor.core.json.JsonNode;
import abstractor.core.json.JsonObject;

public class Ref<T extends Construct> extends ConstructImp {
    private final ConstructKind conKind;
    public  final CtElement     elem;
    public  final String        context;

    // TODO: Check if the typeArgs are needed. If the element is different based
    // on the typeArgs, this can be removed. However, if the element is the same
    // inside different typeArgs context, then this will have to be kept.
    public final ArrayList<TypeDesc> typeArgs = new ArrayList<TypeDesc>();

    private T res;
    
    public Ref(ConstructKind kind, CtElement elem, String context) throws Exception {
        this(kind, elem, context, null);
    }

    public Ref(ConstructKind kind, CtElement elem, String context, List<TypeDesc> typeArgs) throws Exception {
        if (context.isBlank())
            throw new Exception("May not have a blank reference context.");
        this.conKind = kind;
        this.elem    = elem;
        this.context = context;
        if (typeArgs != null) this.typeArgs.addAll(typeArgs);
    }

    public ConstructKind kind() { return this.conKind; }

    @Override
    public int getIndex() {
        return this.isResolved() ? this.res.getIndex() : super.getIndex();
    }

    public T getResolved() { return this.res; }
    public boolean isResolved() { return this.res != null; }

    public void setResolved(T res) throws Exception {
        if (res == null)
            throw new Exception("Attempted to write null as the resolved construct to the reference " + this);
        if (!res.kind().equals(this.conKind))
            throw new Exception("Attempted to write a resolved construct with the kind " + res.kind() + " for reference " + this + " with kind " + this.conKind);
        if (this.isResolved() && !this.res.equals(res))
            throw new Exception("Attempted to overwrite the resolved construct, " + this.res + ", with " + res + " for reference " + this);
        this.res = res;
    }

    public JsonNode refToJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("ref",     true);
        obj.put("context", this.context);
        obj.putNotEmpty("typeArgs", keyList(this.typeArgs));

        //obj.put("refHash", this.hashCode());
        //if (this.isResolved())
        //    obj.put("resHash", this.res.hashCode());
        //if (this.elem != null) {
        //    obj.put("elemHash", this.elem.hashCode());
        //    obj.put("elemType", this.elem.getClass().getSimpleName());
        //}
        return obj;
    }

    @Override
    public JsonNode toJson(JsonHelper h) {
        return (!h.refSkipResolve && this.isResolved()) ? this.res.toJson(h) : this.refToJson(h);
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        if (options.useResolved) {
            return Cmp.or(super.getCmp(c, options),
                Cmp.defer(this.res, () -> ((Ref<?>)c).res)
            );
        }

        return Cmp.or(super.getCmp(c, options),
            Cmp.deferHash(this.elem,     () -> ((Ref<?>)c).elem),
            Cmp.defer(    this.context,  () -> ((Ref<?>)c).context),
            Cmp.deferList(this.typeArgs, () -> ((Ref<?>)c).typeArgs)
        );
    }
}
