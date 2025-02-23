package abstractor.core.constructs;

import java.util.Collections;
import java.util.List;
import java.util.TreeSet;

import spoon.reflect.declaration.CtMethod;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class MethodDecl extends Declaration implements Method {
    public final ObjectDecl receiver;
    public final Metrics metrics;
    public final Signature signature;
    public final List<TypeParam> typeParams;
    public final TreeSet<MethodInst> instances;

    public MethodDecl(CtMethod<?> src, PackageCon pkg, ObjectDecl receiver, Location loc,
        String name, Signature signature, List<TypeParam> typeParams, Metrics metrics) {
        super(src, pkg, loc, name);
        this.receiver = receiver;
        this.metrics = metrics;
        this.signature = signature;
        this.typeParams = Collections.unmodifiableList(typeParams);
        this.instances = new TreeSet<MethodInst>();
    }

    public String kind() { return "method"; }
    
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
