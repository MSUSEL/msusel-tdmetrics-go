package abstractor.core.constructs;

import java.util.List;
import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class InterfaceDecl extends DeclarationImp implements TypeDeclaration {
    public final InterfaceDesc inter;
    public final List<TypeParam> typeParams;
    public final TreeSet<InterfaceInst> instances;

    public InterfaceDecl(PackageCon pkg, Location loc,
        String name, InterfaceDesc inter, List<TypeParam> typeParams) {
        super(pkg, loc, name);
        this.inter = inter;
        this.typeParams = unmodifiableList(typeParams);
        this.instances = new TreeSet<InterfaceInst>();
    }

    public ConstructKind kind() { return ConstructKind.INTERFACE_DECL; }
    
    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("instances",  indexSet(this.instances));
        obj.put("interface",          index(this.inter));
        obj.putNotEmpty("typeParams", indexList(this.typeParams));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.inter,          () -> ((InterfaceDecl)c).inter),
            Cmp.deferList(this.typeParams, () -> ((InterfaceDecl)c).typeParams)
        );
    }
}
