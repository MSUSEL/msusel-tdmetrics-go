package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.List;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class MethodInst extends ConstructImp implements Method {
    public Ref<MethodDecl> generic;
    public Ref<ObjectInst> receiver;
    public final ArrayList<Ref<? extends TypeDesc>> instanceTypes = new ArrayList<Ref<? extends TypeDesc>>();
    public Ref<Signature> resolved;

    public MethodInst() {}

    public MethodInst(Ref<MethodDecl> generic, Ref<ObjectInst> receiver,
        List<Ref<TypeDesc>> instanceTypes, Ref<Signature> resolved) {
        this.generic  = generic;
        this.receiver = receiver;
        this.instanceTypes.addAll(instanceTypes);
        this.resolved = resolved;
    }

    public ConstructKind kind() { return ConstructKind.METHOD_INST; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("generic",       index(this.generic));
        obj.put("instanceTypes", keyList(this.instanceTypes));
        obj.put("receiver",      index(this.receiver));
        obj.put("resolved",      index(this.resolved));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(    this.generic,       () -> ((MethodInst)c).generic),
            Cmp.deferList(this.instanceTypes, () -> ((MethodInst)c).instanceTypes),
            Cmp.defer(    this.resolved,      () -> ((MethodInst)c).resolved)
        );
    }   
}
