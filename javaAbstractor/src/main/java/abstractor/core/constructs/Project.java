package abstractor.core.constructs;

import abstractor.core.json.*;

public class Project implements Jsonable {
    // TODO: abstracts
    // TODO: arguments
    // TODO: basics
    // TODO: fields
    public final Factory<InterfaceDecl> interfaceDecls = new Factory<InterfaceDecl>();
    // TODO: interfaceDescs
    // TODO: interfaceInsts
    // TODO: methods
    // TODO: methodInsts
    // TODO: metrics
    public final Factory<ObjectDecl> objectDecls = new Factory<ObjectDecl>();
    // TODO: objectInsts
    public final Factory<Package> packages = new Factory<Package>();
    // TODO: selections
    // TODO: signatures
    // TODO: structDescs
    // TODO: typeParams
    // TODO: values

    public void setIndices() {
        this.interfaceDecls.setIndices();
        this.objectDecls.setIndices();
        this.packages.setIndices();
    }

    public JsonNode toJson(JsonHelper h) {
        this.setIndices();

        JsonObject obj = new JsonObject();
        obj.put("language", "java");
        // TODO: locs

        // TODO: abstracts
        // TODO: arguments
        // TODO: basics
        // TODO: fields
        obj.putNotEmpty("interfaceDecls", this.interfaceDecls.toJson(h));
        // TODO: interfaceDescs
        // TODO: interfaceInsts
        // TODO: methods
        // TODO: methodInsts
        // TODO: metrics
        obj.putNotEmpty("objects", this.objectDecls.toJson(h));
        // TODO: objectInsts
        obj.putNotEmpty("packages", this.packages.toJson(h));
        // TODO: selections
        // TODO: signatures
        // TODO: structDescs
        // TODO: typeParams
        // TODO: values

        return obj;
    }   
}
