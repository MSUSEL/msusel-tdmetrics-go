package abstractor.core.constructs;

import java.util.List;
import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class ObjectDecl extends DeclarationImp implements TypeDeclaration {
    public final StructDesc struct;
    public final TreeSet<MethodDecl> methodDecls;
    public final List<TypeParam> typeParams;
    public final TreeSet<InterfaceInst> instances;
    
    public InterfaceDesc inter;

    public ObjectDecl(PackageCon pkg, Location loc,
        String name, StructDesc struct, List<TypeParam> typeParams) {
        super(pkg, loc, name);
        this.struct = struct;
        this.methodDecls = new TreeSet<MethodDecl>();
        this.typeParams = unmodifiableList(typeParams);
        this.instances = new TreeSet<InterfaceInst>();
    }

    public ConstructKind kind() { return ConstructKind.OBJECT_DECL; }
    
    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("data",               index(this.struct));
        obj.putNotEmpty("instances",  indexSet(this.instances));
        obj.putNotEmpty("methods",    indexSet(this.methodDecls));
        obj.putNotEmpty("typeParams", indexList(this.typeParams));
        obj.put("interface",          index(this.inter));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.struct,         () -> ((ObjectDecl)c).struct),
            Cmp.deferList(this.typeParams, () -> ((ObjectDecl)c).typeParams)
        );
    }
}
