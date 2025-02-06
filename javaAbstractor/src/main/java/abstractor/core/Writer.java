package abstractor.core;

import abstractor.core.json.*;
import abstractor.core.log.Logger;
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.*;

public class Writer {
    private final Logger log;
    private final Project proj;
    private final Indexer indexer;
    private final boolean writeKinds;
    private final boolean writeIndices;

    static public JsonNode toJson(Logger log, Project proj, boolean writeKinds, boolean writeIndices) {
        return new Writer(log, proj, writeKinds, writeIndices).projectToJson();
    }

    private Writer(Logger log, Project proj, boolean writeKinds, boolean writeIndices) {
        this.log = log;
        this.proj = proj;
        this.indexer = new Indexer(this.log, this.proj);
        this.writeKinds = writeKinds;
        this.writeIndices = writeIndices;
    }

    private int indexOf(Object o) { return this.indexer.indexOf(o); }

    static private String packagePath(CtPackage p) {
        SourcePosition pos = p.getPosition();
        if (!pos.isValidPosition()) return "";
        
        String path = pos.getFile().getPath();
        final String tail = "package-info.java";
        if (path.endsWith(tail))
            path = path.substring(0, path.length()-tail.length());
        return path;
    }

    private JsonNode packageToJson(CtPackage p) {
        JsonObject obj = new JsonObject();
        if (this.writeKinds) obj.put("kind", "package");
        if (this.writeIndices) obj.put("index", this.indexOf(p));
        obj.putNotEmpty("path", packagePath(p));
        obj.put("name", p.getQualifiedName());

        // TODO: imports
        // TODO: interfaces
        // TODO: methods
        // TODO: objects
        // TODO: values

        return obj;
    }

    private JsonNode packageSetToJson() {
        JsonArray array = new JsonArray();
        for (CtPackage p : this.proj.packages)
            array.add(this.packageToJson(p));
        return array;
    }

    private JsonNode objectToJson(CtClass<?> c) {
        JsonObject obj = new JsonObject();
        if (this.writeKinds) obj.put("kind", "object");
        if (this.writeIndices) obj.put("index", this.indexOf(c));
        obj.put("name", c.getSimpleName());
        obj.put("package", this.indexOf(c.getPackage()));

        // TODO: add more

        return obj;
    }

    private JsonNode objectSetToJson() {
        JsonArray array = new JsonArray();
        for (CtClass<?> c : this.proj.objects)
            array.add(this.objectToJson(c));
        return array;
    }

    private JsonNode projectToJson() {
        JsonObject obj = new JsonObject();

        // TODO: abstracts
        // TODO: arguments
        // TODO: basics
        // TODO: fields
        // TODO: interfaceDecls
        // TODO: interfaceDescs
        // TODO: interfaceInsts

        obj.put("language", "java");

        // TODO: locs
        // TODO: methods
        // TODO: methodInsts
        // TODO: metrics
        
        obj.putNotEmpty("objects", this.objectSetToJson());

        // TODO: objectInsts
        
        obj.putNotEmpty("packages", this.packageSetToJson());

        // TODO: selections
        // TODO: signatures
        // TODO: structDescs
        // TODO: typeParams
        // TODO: values

        return obj;
    }   
}
