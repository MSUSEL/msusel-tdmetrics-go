package abstractor.core.constructs;

import java.util.*;

import abstractor.core.cmp.*;
import abstractor.core.iter.*;
import abstractor.core.json.*;
import abstractor.core.require.Require;

public class InterfaceInst extends ConstructImp implements TypeDesc {
    public Ref<InterfaceDecl> generic;
    public final ArrayList<Ref<? extends TypeDesc>> instanceTypes = new ArrayList<>();
    public Ref<InterfaceDesc> resolved;

    public InterfaceInst() {}

    public InterfaceInst(Ref<InterfaceDecl> generic, List<Ref<? extends TypeDesc>> instanceTypes, Ref<InterfaceDesc> resolved) throws Exception {
        this.generic = generic;
        if (instanceTypes != null) this.instanceTypes.addAll(instanceTypes);
        this.resolved = resolved;
        
        if (generic.isResolved()) {
            final ArrayList<Ref<TypeParam>> tp = generic.getResolved().typeParams;
            final int tpSize = tp.size();
            final int taSize = instanceTypes != null ? instanceTypes.size() : 0;
            Require.equal(tpSize, taSize, () -> "The interface's type params count (" + tpSize + ") " +
                "must match the type arguments (" + taSize + "):\n"+
                "  type params: " + tp + "\n"+
                "  instance types: " + instanceTypes + "\n"+
                "  instance: " + this + "\n"+
                "  generic: " + this.generic.getResolved());
        }
    }

    public ConstructKind kind() { return ConstructKind.INTERFACE_INST; }

    public boolean matchesGeneric() {
        if (!generic.isResolved()) return false;
        final ArrayList<Ref<TypeParam>> tp = generic.getResolved().typeParams;
        for (int i = 0; i < tp.size(); i++) {
            if (!tp.get(i).equals(instanceTypes.get(i))) return false;
        }
        return true;
    }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("generic",       index(this.generic));
        obj.put("instanceTypes", keyList(this.instanceTypes));
        obj.put("resolved",      index(this.resolved));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.defer(    this.generic,       () -> ((InterfaceInst)c).generic),
            Cmp.deferList(this.instanceTypes, () -> ((InterfaceInst)c).instanceTypes),
            Cmp.defer(    this.resolved,      () -> ((InterfaceInst)c).resolved)
        );
    }

    @Override
    public Iterator<Ref<? extends Construct>> subConstructs() {
        return new Bundle<Ref<? extends Construct>>().
            add(this.generic).
            add(this.instanceTypes).
            add(this.resolved);
    }
}
