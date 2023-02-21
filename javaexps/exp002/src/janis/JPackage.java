package janis;

import json.JsonMap;
import json.JsonObj;
import json.Jsonable;
import spoon.reflect.declaration.CtPackage;

import java.util.ArrayList;
import java.util.List;

public final class JPackage implements Jsonable {
    public String name;
    public List<String> subpackages;

    public JPackage(CtPackage p) {
        this.name = p.isUnnamedPackage() ? "<root>" : p.getSimpleName();
        this.subpackages = new ArrayList<>();
        for (CtPackage sub : p.getPackages())
            this.subpackages.add(sub.getSimpleName());
    }

    @Override
    public String toString() {
        return this.name;
    }

    @Override
    public JsonObj toJson() {
        return new JsonMap().
            with("name", this.name).
            withOmitEmpty("subpackages", this.subpackages);
    }
}
