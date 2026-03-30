package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.List;
import java.util.TreeSet;
import java.util.SortedSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class ObjectInst extends ConstructImp implements TypeDesc {
    public ObjectDecl generic;
    public final List<TypeDesc> instanceTypes;
    public final SortedSet<MethodInst> methods;
    public StructDesc resData;
    public InterfaceDesc resInterface;

    public ObjectInst() {
        this.instanceTypes = new ArrayList<TypeDesc>();
        this.methods       = new TreeSet<MethodInst>();
    }

    public ObjectInst(ObjectDecl generic, List<TypeDesc> instanceTypes,
        StructDesc resData, InterfaceDesc resInterface) {
        this.generic       = generic;
        this.instanceTypes = instanceTypes;
        this.methods       = new TreeSet<MethodInst>();
        this.resData       = resData;
        this.resInterface  = resInterface;
    }

    public ConstructKind kind() { return ConstructKind.OBJECT_INST; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("generic",       index(this.generic));
        obj.put("instanceTypes", keyList(this.instanceTypes));
        obj.put("methods",       indexSet(this.methods));
        obj.put("resData",       index(this.resData));
        obj.put("resInterface",  index(this.resInterface));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.generic,           () -> ((ObjectInst)c).generic),
            Cmp.deferList(this.instanceTypes, () -> ((ObjectInst)c).instanceTypes),
            Cmp.defer(this.resData,           () -> ((ObjectInst)c).resData),
            Cmp.defer(this.resInterface,      () -> ((ObjectInst)c).resInterface)
        );
    }   
}
