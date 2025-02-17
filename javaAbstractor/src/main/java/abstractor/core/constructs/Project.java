package abstractor.core.constructs;

import abstractor.core.json.*;

public class Project implements Jsonable {
    public final Locations locations = new Locations();
    public final Factory<Abstract> abstracts = new Factory<Abstract>();
    public final Factory<Argument> arguments = new Factory<Argument>();
    public final Factory<Basic> basics = new Factory<Basic>();
    public final Factory<Field> fields = new Factory<Field>();
    public final Factory<InterfaceDecl> interfaceDecls = new Factory<InterfaceDecl>();
    // TODO: interfaceDescs
    // TODO: interfaceInsts
    public final Factory<MethodDecl> methodDecls = new Factory<MethodDecl>();
    // TODO: methodInsts
    // TODO: metrics
    public final Factory<ObjectDecl> objectDecls = new Factory<ObjectDecl>();
    // TODO: objectInsts
    public final Factory<PackageCon> packages = new Factory<PackageCon>();
    // TODO: selections
    // TODO: signatures
    public final Factory<StructDesc> structDescs = new Factory<StructDesc>();
    // TODO: typeParams
    // TODO: values

    private void prepareForOutput() {
        this.locations.prepareForOutput();
        this.abstracts.setIndices();
        this.arguments.setIndices();
        this.basics.setIndices();
        this.fields.setIndices();
        this.interfaceDecls.setIndices();
        // TODO: interfaceDescs
        // TODO: interfaceInsts
        this.methodDecls.setIndices();
        // TODO: methodInsts
        // TODO: metrics
        this.objectDecls.setIndices();
        // TODO: objectInsts
        this.packages.setIndices();
        // TODO: selections
        // TODO: signatures
        this.structDescs.setIndices();
        // TODO: typeParams
        // TODO: values
    }

    public JsonNode toJson(JsonHelper h) {
        this.prepareForOutput();

        JsonObject obj = new JsonObject();
        obj.put("language", "java");
        obj.putNotEmpty("locs", this.locations.toJson(h));
        obj.putNotEmpty("abstracts", this.abstracts.toJson(h));
        obj.putNotEmpty("arguments", this.arguments.toJson(h));
        obj.putNotEmpty("basics", this.basics.toJson(h));
        obj.putNotEmpty("fields", this.fields.toJson(h));
        obj.putNotEmpty("interfaceDecls", this.interfaceDecls.toJson(h));
        // TODO: interfaceDescs
        // TODO: interfaceInsts
        obj.putNotEmpty("methods", this.methodDecls.toJson(h));
        // TODO: methodInsts
        // TODO: metrics
        obj.putNotEmpty("objects", this.objectDecls.toJson(h));
        // TODO: objectInsts
        obj.putNotEmpty("packages", this.packages.toJson(h));
        // TODO: selections
        // TODO: signatures
        obj.putNotEmpty("structDescs", this.structDescs.toJson(h));
        // TODO: typeParams
        // TODO: values
        return obj;
    }   
}
