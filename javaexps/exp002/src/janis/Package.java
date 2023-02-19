package janis;

import json.JsonMap;
import json.JsonObj;
import json.Jsonable;
import spoon.reflect.declaration.CtPackage;

import java.util.ArrayList;
import java.util.List;

public class Package implements Jsonable {
    public String name;
    public List<Package> packages;

    public Package(CtPackage p) {
        this.name = p.getQualifiedName();
        this.packages = new ArrayList<>();
        for (CtPackage sub : p.getPackages())
            this.packages.add(new Package(sub));
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
