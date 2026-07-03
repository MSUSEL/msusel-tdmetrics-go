package abstractor.core.constructs;

import java.util.*;
    
import abstractor.core.cmp.*;
import abstractor.core.iter.*;
import abstractor.core.json.*;

public class MethodDecl extends DeclarationImp implements Method {
    public Ref<ObjectDecl> receiver;
    public Ref<Signature>  signature;
    public final ArrayList<Ref<TypeParam>> typeParams = new ArrayList<>();
    public final TreeSet<Ref<MethodInst>>  instances  = new TreeSet<>();
    public boolean constructor;
    public Ref<Metrics> metrics;

    public MethodDecl(Ref<PackageCon> pkg, Ref<ObjectDecl> receiver, Location loc,
        String name, Ref<Signature> signature, List<Ref<TypeParam>> typeParams) throws Exception {
        super(pkg, loc, name);
        this.receiver  = receiver;
        this.signature = signature;
        if (typeParams != null) this.typeParams.addAll(typeParams);
    }

    public ConstructKind kind() { return ConstructKind.METHOD_DECL; }
    
    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("instances",   indexSet(this.instances));
        obj.putNotEmpty("metrics",     index(this.metrics));
        obj.putNotEmpty("receiver",    index(this.receiver));
        obj.putNotEmpty("signature",   index(this.signature));
        obj.putNotEmpty("typeParams",  indexList(this.typeParams));
        obj.putNotEmpty("constructor", this.constructor);
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.defer(    this.receiver,    () -> ((MethodDecl)c).receiver),
            Cmp.defer(    this.signature,   () -> ((MethodDecl)c).signature),
            Cmp.deferList(this.typeParams,  () -> ((MethodDecl)c).typeParams)
        );
    }

    @Override
    public Iterator<Ref<? extends Construct>> subConstructs() {
        return new Bundle<Ref<? extends Construct>>().
            add(this.receiver).
            add(this.signature).
            add(this.typeParams).
            add(this.instances).
            add(this.metrics);
    }
}
