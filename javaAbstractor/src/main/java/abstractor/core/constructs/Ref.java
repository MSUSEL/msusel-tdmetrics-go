package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.List;

import spoon.reflect.declaration.CtElement;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.JsonHelper;
import abstractor.core.json.JsonNode;
import abstractor.core.json.JsonObject;

public abstract class Ref<T extends Construct> extends ConstructImp {
    private final ConstructKind conKind;
    public  final CtElement     elem;
    public  final String        context;
    public  final String        name;

    // TODO: Check if the typeArgs are needed. If the element is different based
    // on the typeArgs, this can be removed. However, if the element is the same
    // inside different typeArgs context, then this will have to be kept.
    public final ArrayList<TypeDesc> typeArgs = new ArrayList<TypeDesc>();

    private T res;
    
    public Ref(ConstructKind kind, CtElement elem, String context, String name) {
        this.conKind = kind;
        this.elem    = elem;
        this.context = context;
        this.name    = name;
    }

    public Ref(ConstructKind kind, CtElement elem, String context, String name, List<TypeDesc> typeArgs) {
        this.conKind = kind;
        this.elem    = elem;
        this.context = context;
        this.name    = name;
        this.typeArgs.addAll(typeArgs);
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
        if (this.isResolved())
            throw new Exception("Attempted to overwrite the resolved construct, " + this.res + ", with " + res + " for reference " + this);
        if (res.kind() != this.conKind)
            throw new Exception("Attempted to write a resolved construct with the kind " + res.kind() + " for reference " + this + " with kind " + this.conKind);
        this.res = res;
    }

    @Override
    public JsonNode toJson(JsonHelper h) {
        if (this.isResolved()) return this.res.toJson(h);

        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("name",             this.name);
        obj.putNotEmpty("context",  this.context);
        obj.putNotEmpty("resolved", key(this.res));
        obj.putNotEmpty("typeArgs", keyList(this.typeArgs));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        if (this.isResolved()) return this.res.compareTo(c);

        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.name,         () -> ((Ref<?>)c).name),
            Cmp.defer(this.context,      () -> ((Ref<?>)c).context),
            Cmp.deferList(this.typeArgs, () -> ((Ref<?>)c).typeArgs)
        );
    }
}
