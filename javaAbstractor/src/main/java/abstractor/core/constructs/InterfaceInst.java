package abstractor.core.constructs;

import java.util.List;

import spoon.reflect.declaration.CtField;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class InterfaceInst extends ConstructImp implements TypeDesc {
    public final InterfaceDecl generic;
    public final List<TypeDesc> instanceTypes;
    public final InterfaceDesc resolved;
    
    public InterfaceInst(CtField<?> src, InterfaceDecl generic, List<TypeDesc> instanceTypes, InterfaceDesc resolved) {
        super(src);
        this.generic = generic;
        this.instanceTypes = unmodifiableList(instanceTypes);
        this.resolved = resolved;
    }

    public String kind() { return "interfaceInst"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("generic",       index(this.generic));
        obj.put("instanceTypes", indexList(this.instanceTypes));
        obj.put("resolved",      index(this.resolved));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.generic,           () -> ((InterfaceInst)c).generic),
            Cmp.deferList(this.instanceTypes, () -> ((InterfaceInst)c).instanceTypes),
            Cmp.defer(this.resolved,          () -> ((InterfaceInst)c).resolved)
        );
    }   
}
