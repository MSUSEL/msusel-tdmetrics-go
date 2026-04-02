package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.List;
import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class ObjectInst extends ConstructImp implements TypeDesc {
    public Ref<ObjectDecl> generic;
    public final ArrayList<Ref<? extends TypeDesc>> instanceTypes = new ArrayList<Ref<? extends TypeDesc>>();
    public final TreeSet<Ref<MethodInst>> methods = new TreeSet<Ref<MethodInst>>();
    public Ref<StructDesc>    resData;
    public Ref<InterfaceDesc> resInterface;

    public ObjectInst() {}

    public ObjectInst(Ref<ObjectDecl> generic, List<Ref<? extends TypeDesc>> instanceTypes,
        Ref<StructDesc> resData, Ref<InterfaceDesc> resInterface) {
        this.generic = generic;
        if (instanceTypes != null) this.instanceTypes.addAll(instanceTypes);
        this.resData      = resData;
        this.resInterface = resInterface;
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
            Cmp.defer(    this.generic,       () -> ((ObjectInst)c).generic),
            Cmp.deferList(this.instanceTypes, () -> ((ObjectInst)c).instanceTypes),
            Cmp.defer(    this.resData,       () -> ((ObjectInst)c).resData),
            Cmp.defer(    this.resInterface,  () -> ((ObjectInst)c).resInterface)
        );
    }   
}
