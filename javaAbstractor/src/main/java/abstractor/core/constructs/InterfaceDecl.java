package abstractor.core.constructs;

import java.util.*;

import abstractor.core.cmp.*;
import abstractor.core.iter.*;
import abstractor.core.json.*;

public class InterfaceDecl extends DeclarationImp implements TypeDeclaration {
    public       Ref<InterfaceDesc>          inter;
    public final ArrayList<Ref<TypeParam>>   typeParams  = new ArrayList<>();
    public final TreeSet<Ref<InterfaceInst>> instances   = new TreeSet<>();

    public InterfaceDecl(Ref<PackageCon> pkg, Location loc,
        String name, Ref<InterfaceDesc> inter, List<Ref<TypeParam>> typeParams) throws Exception {
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

    @Override
    public Iterator<Ref<? extends Construct>> subConstructs() {
        return new Bundle<Ref<? extends Construct>>().
            add(this.inter).
            add(this.typeParams).
            add(this.instances);
    }
}
