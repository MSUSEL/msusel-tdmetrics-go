package abstractor.core.constructs;

import java.util.List;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class MethodInst extends ConstructImp implements Method {
    public final MethodDecl generic;
    public final ObjectInst receiver;
    public final List<TypeDesc> instanceTypes;
    public final Signature resolved;

    public MethodInst(MethodDecl generic, ObjectInst receiver,
        List<TypeDesc> instanceTypes, Signature resolved) {
        this.generic = generic;
        this.receiver = receiver;
        this.instanceTypes = unmodifiableList(instanceTypes);
        this.resolved = resolved;
    }

    public String kind() { return "methodInst"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("generic",       index(this.generic));
        obj.put("instanceTypes", indexList(this.instanceTypes));
        obj.put("receiver",      index(this.receiver));
        obj.put("resolved",      index(this.resolved));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.generic,           () -> ((MethodInst)c).generic),
            Cmp.deferList(this.instanceTypes, () -> ((MethodInst)c).instanceTypes),
            Cmp.defer(this.resolved,          () -> ((MethodInst)c).resolved)
        );
    }   
}
