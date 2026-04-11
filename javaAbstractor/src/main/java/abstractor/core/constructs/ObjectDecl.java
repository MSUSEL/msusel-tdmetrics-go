package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.List;
import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;

public class ObjectDecl extends DeclarationImp implements TypeDeclaration {
    public Ref<StructDesc>    struct;
    public Ref<InterfaceDesc> inter;
    public final TreeSet<Ref<MethodDecl>>    methodDecls = new TreeSet<Ref<MethodDecl>>();
    public final ArrayList<Ref<TypeParam>>   typeParams  = new ArrayList<Ref<TypeParam>>();
    public final TreeSet<Ref<InterfaceInst>> instances   = new TreeSet<Ref<InterfaceInst>>();

    public ObjectDecl() {}

    public ObjectDecl(Ref<PackageCon> pkg, Location loc,
        String name, Ref<StructDesc> struct, List<Ref<TypeParam>> typeParams) {
        super(pkg, loc, name);
        this.struct = struct;
        if (typeParams != null) this.typeParams.addAll(typeParams);
    }

    public ConstructKind kind() { return ConstructKind.OBJECT_DECL; }
    
    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put(        "data",       index(this.struct));
        obj.put(        "interface",  index(this.inter));
        obj.putNotEmpty("instances",  indexSet(this.instances));
        obj.putNotEmpty("methods",    indexSet(this.methodDecls));
        obj.putNotEmpty("typeParams", indexList(this.typeParams));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.defer(    this.struct,     () -> ((ObjectDecl)c).struct),
            Cmp.deferList(this.typeParams, () -> ((ObjectDecl)c).typeParams)
        );
    }
}
