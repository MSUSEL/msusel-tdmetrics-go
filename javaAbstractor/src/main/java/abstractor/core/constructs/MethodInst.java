package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.List;

import abstractor.core.cmp.Cmp;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;

public class MethodInst extends ConstructImp implements Method {
    public Ref<MethodDecl> generic;
    // receiver is either the generic's ObjectDecl (when the class isn't
    // instantiated at this call site) or an ObjectInst (when it is), or
    // null only when the declaring class produced neither decl nor inst.
    // Serialized as a key so consumers can tell ObjectDecl from ObjectInst.
    public Ref<? extends TypeDesc> receiver;

    public final ArrayList<Ref<? extends TypeDesc>> instanceTypes = new ArrayList<>();
    public Ref<Signature> resolved;

    public MethodInst() {}

    public MethodInst(Ref<MethodDecl> generic, Ref<? extends TypeDesc> receiver,
        List<Ref<? extends TypeDesc>> instanceTypes, Ref<Signature> resolved) {
        this.generic  = generic;
        this.receiver = receiver;
        if (instanceTypes != null) this.instanceTypes.addAll(instanceTypes);
        this.resolved = resolved;
    }

    public ConstructKind kind() { return ConstructKind.METHOD_INST; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put(        "generic",       index(this.generic));
        obj.put(        "instanceTypes", keyList(this.instanceTypes));
        obj.putNotEmpty("receiver",      key(this.receiver));
        obj.put(        "resolved",      index(this.resolved));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.defer(    this.generic,       () -> ((MethodInst)c).generic),
            Cmp.deferList(this.instanceTypes, () -> ((MethodInst)c).instanceTypes),
            Cmp.defer(    this.resolved,      () -> ((MethodInst)c).resolved)
        );
    }   
}
