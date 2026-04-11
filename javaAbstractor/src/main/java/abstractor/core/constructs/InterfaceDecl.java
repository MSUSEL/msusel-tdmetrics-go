package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.List;
import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;

public class InterfaceDecl extends DeclarationImp implements TypeDeclaration {
    public       Ref<InterfaceDesc>          inter;
    public final ArrayList<Ref<TypeParam>>   typeParams = new ArrayList<Ref<TypeParam>>();
    public final TreeSet<Ref<InterfaceInst>> instances  = new TreeSet<Ref<InterfaceInst>>();

    public InterfaceDecl() {}

    public InterfaceDecl(Ref<PackageCon> pkg, Location loc,
        String name, Ref<InterfaceDesc> inter, List<Ref<TypeParam>> typeParams) {
        super(pkg, loc, name);
        this.inter = inter;
        if (typeParams != null) this.typeParams.addAll(typeParams);
    }

    public ConstructKind kind() { return ConstructKind.INTERFACE_DECL; }
    
    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("instances",  indexSet(this.instances));
        obj.put(        "interface",  index(this.inter));
        obj.putNotEmpty("typeParams", indexList(this.typeParams));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.defer(    this.inter,      () -> ((InterfaceDecl)c).inter),
            Cmp.deferList(this.typeParams, () -> ((InterfaceDecl)c).typeParams)
        );
    }
}
