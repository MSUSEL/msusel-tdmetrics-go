package abstractor.core;

import abstractor.core.json.*;
import abstractor.core.log.Logger;
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.*;

public class Writer {
    private JsonNode abstractToJson(CtMethod<?> m) {
        JsonObject obj = new JsonObject();
        if (this.writeKinds) obj.put("kind", "abstract");
        if (this.writeIndices) obj.put("index", this.indexer.index(m));
        obj.put("name", m.getSimpleName());
    
        // TODO: Implement
        return obj;
    }

    private JsonNode methodToJson(CtMethod<?> m) {
        JsonObject obj = new JsonObject();
        if (this.writeKinds) obj.put("kind", "method");
        if (this.writeIndices) obj.put("index", this.indexer.index(m));
        obj.put("name", m.getSimpleName());
    
        // TODO: Implement
        return obj;
    }

    private JsonNode objectToJson(CtClass<?> c) {
        JsonObject obj = new JsonObject();
        if (this.writeKinds) obj.put("kind", "object");
        if (this.writeIndices) obj.put("index", this.indexer.index(c));
        obj.put("name", c.getSimpleName());
        obj.put("package", this.indexer.index(c.getPackage()));

        // TODO: Handle enum
        //if (c instanceof CtEnum<?> e) {}

        // TODO: data
        // TODO: exported
        // TODO: index
        // TODO: instances
        // TODO: loc
        // TODO: methods
        // TODO: typeParams
        // TODO: interface

        return obj;
    }
}
