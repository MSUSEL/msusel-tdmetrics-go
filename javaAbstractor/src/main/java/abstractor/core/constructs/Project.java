package abstractor.core.constructs;

import abstractor.core.json.*;

public class Project implements Jsonable {
    public final Locations              locations      = new Locations();
    public final Factory<Abstract>      abstracts      = new Factory<Abstract>();
    public final Factory<Argument>      arguments      = new Factory<Argument>();
    public final Factory<Basic>         basics         = new Factory<Basic>();
    public final Factory<Field>         fields         = new Factory<Field>();
    public final Factory<InterfaceDecl> interfaceDecls = new Factory<InterfaceDecl>();
    public final Factory<InterfaceDesc> interfaceDescs = new Factory<InterfaceDesc>();
    public final Factory<InterfaceInst> interfaceInsts = new Factory<InterfaceInst>();
    public final Factory<MethodDecl>    methodDecls    = new Factory<MethodDecl>();
    public final Factory<MethodInst>    methodInsts    = new Factory<MethodInst>();
    public final Factory<Metrics>       metrics        = new Factory<Metrics>();
    public final Factory<ObjectDecl>    objectDecls    = new Factory<ObjectDecl>();
    public final Factory<ObjectInst>    objectInsts    = new Factory<ObjectInst>();
    public final Factory<PackageCon>    packages       = new Factory<PackageCon>();
    public final Factory<Selection>     selections     = new Factory<Selection>();
    public final Factory<Signature>     signatures     = new Factory<Signature>();
    public final Factory<StructDesc>    structDescs    = new Factory<StructDesc>();
    public final Factory<TypeParam>     typeParams     = new Factory<TypeParam>();
    public final Factory<Value>         values         = new Factory<Value>();

    private void setAllIndices() {
        this.abstracts.setIndices();
        this.arguments.setIndices();
        this.basics.setIndices();
        this.fields.setIndices();
        this.interfaceDecls.setIndices();
        this.interfaceDescs.setIndices();
        this.interfaceInsts.setIndices();
        this.methodDecls.setIndices();
        this.methodInsts.setIndices();
        this.metrics.setIndices();
        this.objectDecls.setIndices();
        this.objectInsts.setIndices();
        this.packages.setIndices();
        this.selections.setIndices();
        this.signatures.setIndices();
        this.structDescs.setIndices();
        this.typeParams.setIndices();
        this.values.setIndices();
    }

    public JsonNode toJson(JsonHelper h) {
        this.locations.prepareForOutput();
        this.setAllIndices();

        JsonObject obj = new JsonObject();
        obj.put("language", "java");
        obj.putNotEmpty("locs",           this.locations.toJson(h));
        obj.putNotEmpty("abstracts",      this.abstracts.toJson(h));
        obj.putNotEmpty("arguments",      this.arguments.toJson(h));
        obj.putNotEmpty("basics",         this.basics.toJson(h));
        obj.putNotEmpty("fields",         this.fields.toJson(h));
        obj.putNotEmpty("interfaceDecls", this.interfaceDecls.toJson(h));
        obj.putNotEmpty("interfaceDescs", this.interfaceDescs.toJson(h));
        obj.putNotEmpty("interfaceInsts", this.interfaceInsts.toJson(h));
        obj.putNotEmpty("methods",        this.methodDecls.toJson(h));
        obj.putNotEmpty("methodInsts",    this.methodInsts.toJson(h));
        obj.putNotEmpty("metrics",        this.metrics.toJson(h));
        obj.putNotEmpty("objects",        this.objectDecls.toJson(h));
        obj.putNotEmpty("objectInsts",    this.objectInsts.toJson(h));
        obj.putNotEmpty("packages",       this.packages.toJson(h));
        obj.putNotEmpty("selections",     this.selections.toJson(h));
        obj.putNotEmpty("signatures",     this.signatures.toJson(h));
        obj.putNotEmpty("structDescs",    this.structDescs.toJson(h));
        obj.putNotEmpty("typeParams",     this.typeParams.toJson(h));
        obj.putNotEmpty("values",         this.values.toJson(h));
        return obj;
    }   
}
