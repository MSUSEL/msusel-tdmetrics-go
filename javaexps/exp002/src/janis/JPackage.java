package janis;

import json.JsonMap;
import json.JsonObj;
import json.Jsonable;
import spoon.reflect.declaration.CtPackage;

import java.util.ArrayList;
import java.util.List;

public final class JPackage implements Jsonable {
    public String name;
    public List<JPackage> packages;
    public List<JClass> classes;

    public JPackage(CtPackage p) {
        this.name = p.getQualifiedName();
        this.packages = new ArrayList<>();
        for (CtPackage sub : p.getPackages())
            this.packages.add(new JPackage(sub));
    }

    @Override
    public String toString() {
        return this.name;
    }

    @Override
    public JsonObj toJson() {
        JsonMap m = new JsonMap();
        m.put("package", this.name);
        m.put("subpackages", this.packages);
        return m;
    }
}
