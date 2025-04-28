package abstractor.core.constructs;

import java.util.List;
import java.util.TreeSet;
    
import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class MethodDecl extends DeclarationImp implements Method {
    public final ObjectDecl receiver;
    public final Signature signature;
    public final List<TypeParam> typeParams;
    public final TreeSet<MethodInst> instances;
    // TODO: Add a flag to indicate if this method is a constructor or not.
    
    public Metrics metrics;

    public MethodDecl(PackageCon pkg, ObjectDecl receiver, Location loc,
        String name, Signature signature, List<TypeParam> typeParams) {
        super(pkg, loc, name);
        this.receiver = receiver;
        this.signature = signature;
        this.typeParams = unmodifiableList(typeParams);
        this.instances = new TreeSet<MethodInst>();
    }

    public ConstructKind kind() { return ConstructKind.METHOD_DECL; }
    
    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("instances",  indexSet(this.instances));
        obj.putNotEmpty("metrics",    index(this.metrics));
        obj.putNotEmpty("receiver",   index(this.receiver));
        obj.putNotEmpty("signature",  index(this.signature));
        obj.putNotEmpty("typeParams", indexList(this.typeParams));
        return obj;
    }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.receiver,       () -> ((MethodDecl)c).receiver),
            Cmp.defer(this.signature,      () -> ((MethodDecl)c).signature),
            Cmp.deferList(this.typeParams, () -> ((MethodDecl)c).typeParams)
        );
    }
}
