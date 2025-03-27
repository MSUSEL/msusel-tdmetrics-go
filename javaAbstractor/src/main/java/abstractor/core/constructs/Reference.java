package abstractor.core.constructs;

import java.util.Collections;
import java.util.List;

import spoon.reflect.declaration.CtElement;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.JsonHelper;
import abstractor.core.json.JsonNode;
import abstractor.core.json.JsonObject;

public abstract class Reference<T extends Construct> extends ConstructImp {
    public final CtElement elem;
    public final String context;
    public final String name;
    public final List<TypeDesc> typeArguments;

    private T res;
    
    public Reference(CtElement elem, String context, String name, List<TypeDesc> typeArguments) {
        this.elem = elem;
        this.context = context;
        this.name = name;
        this.typeArguments = Collections.unmodifiableList(typeArguments);
    }

    public T getResolved() { return this.res; }
    public boolean isResolved() { return this.res != null; }

    public boolean setResolved(T res) {
        if (this.res == null && res != null) {
            this.res = res;
            return true;
        }
        return false;
    }

    public abstract ConstructKind unresolvedKind();
    
    public ConstructKind kind() {
        if (this.isResolved()) return this.res.kind();
        return this.unresolvedKind();
    }

    @Override
    public int getIndex() {
        if (this.isResolved()) return this.res.getIndex();
        return super.getIndex();
    }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("name", this.name);
        obj.putNotEmpty("context", this.context);
        obj.putNotEmpty("typeArgs", keyList(this.typeArguments));
        obj.putNotEmpty("resolved", key(this.res));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            // Skip `() -> super.compareTo(c),` so `kind` is not used in comparison to ensure a stable ordering.
            Cmp.defer(this.unresolvedKind(), () -> ((Reference<?>)c).unresolvedKind()),
            Cmp.defer(this.name, () -> ((Reference<?>)c).name),
            Cmp.defer(this.context, () -> ((Reference<?>)c).context),
            Cmp.deferList(this.typeArguments, () -> ((Reference<?>)c).typeArguments)
        );
    }
}
