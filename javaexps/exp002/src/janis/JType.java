package janis;

import json.JsonMap;
import json.JsonObj;
import json.Jsonable;
import spoon.reflect.declaration.CtType;

public final class JType implements Jsonable {
    public String name;
    public String packageName;



    public JType(CtType t) {
        this.name = t.getSimpleName();
        this.packageName = t.getPackage().getQualifiedName();
    }

/*
    static bool;
    public bool;
    external bool;
    buildIn bool;

    supers []class.name;
    fields []class.name;
    members []method.name;
 */


    @Override
    public JsonObj toJson() {
        return new JsonMap().
            with("name", this.name).
            with("package", this.packageName);
    }
}
