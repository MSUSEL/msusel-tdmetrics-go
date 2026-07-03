package abstractor.core.constructs;

import java.util.*;

import abstractor.core.cmp.*;
import abstractor.core.iter.*;
import abstractor.core.json.*;

public class ObjectInst extends ConstructImp implements TypeDesc {
    public Ref<ObjectDecl> generic;
    public final ArrayList<Ref<? extends TypeDesc>> instanceTypes = new ArrayList<>();
    public final TreeSet<Ref<MethodInst>>           methods       = new TreeSet<>();
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
        obj.put("resData",       index(this.resData));
        obj.put("resInterface",  index(this.resInterface));
        obj.putNotEmpty("methods", indexSet(this.methods));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.defer(    this.generic,       () -> ((ObjectInst)c).generic),
            Cmp.deferList(this.instanceTypes, () -> ((ObjectInst)c).instanceTypes),
            Cmp.defer(    this.resData,       () -> ((ObjectInst)c).resData),
            Cmp.defer(    this.resInterface,  () -> ((ObjectInst)c).resInterface)
        );
    }

    @Override
    public Iterator<Ref<? extends Construct>> subConstructs() {
        return new Bundle<Ref<? extends Construct>>().
            add(this.generic).
            add(this.instanceTypes).
            add(this.methods).
            add(this.resData).
            add(this.resInterface);
    }
}
