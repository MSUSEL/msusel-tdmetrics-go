package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.List;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class InterfaceInst extends ConstructImp implements TypeDesc {
    public Ref<InterfaceDecl> generic;
    public final ArrayList<Ref<? extends TypeDesc>> instanceTypes = new ArrayList<Ref<? extends TypeDesc>>();
    public Ref<InterfaceDesc> resolved;

    public InterfaceInst() {}

    public InterfaceInst(Ref<InterfaceDecl> generic, List<Ref<? extends TypeDesc>> instanceTypes, Ref<InterfaceDesc> resolved) {
        this.generic = generic;
        this.instanceTypes.addAll(instanceTypes);
        this.resolved = resolved;
    }

    public ConstructKind kind() { return ConstructKind.INTERFACE_INST; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("generic",       index(this.generic));
        obj.put("instanceTypes", keyList(this.instanceTypes));
        obj.put("resolved",      index(this.resolved));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(    this.generic,       () -> ((InterfaceInst)c).generic),
            Cmp.deferList(this.instanceTypes, () -> ((InterfaceInst)c).instanceTypes),
            Cmp.defer(    this.resolved,      () -> ((InterfaceInst)c).resolved)
        );
    }   
}
