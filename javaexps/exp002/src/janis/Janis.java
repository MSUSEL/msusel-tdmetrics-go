package janis;

import json.JsonList;
import json.JsonMap;
import spoon.Launcher;
import spoon.reflect.CtModel;
import spoon.reflect.declaration.CtMethod;
import spoon.reflect.declaration.CtPackage;
import spoon.reflect.declaration.CtType;

public class Janis {
    private Janis() { }

    static public JsonMap read(String path) {
        Launcher sl = new Launcher();
        sl.addInputResource(path);
        sl.getEnvironment().setNoClasspath(true);
        sl.getEnvironment().setComplianceLevel(7);

        sl.buildModel();
        return readModel(sl.getModel());
    }

    static private String readPackageName(CtPackage p) {
        return p.isUnnamedPackage() ? "<root>" : p.getSimpleName();
    }

    static private JsonMap readModel(CtModel model) {
        return new JsonMap().
            withOmitOnDefault("packages", readPackages(model)).
            withOmitOnDefault("types", readTypes(model)).
            withOmitOnDefault("methods", readMethods(model));
    }

    static private JsonList readPackages(CtModel model) {
        JsonList packages = new JsonList();
        for (CtPackage p : model.getAllPackages()) {
            JsonList subpackageData = new JsonList();
            for (CtPackage sub : p.getPackages())
                subpackageData.with(readPackageName(sub));

            packages.addMap().
                with("name", readPackageName(p)).
                withOmitOnDefault("subpackages", subpackageData);
        }
        return packages;
    }

    static private JsonList readTypes(CtModel model) {
        JsonList typesData = new JsonList();
        for (CtType<?> t : model.getAllTypes()) {
            typesData.addMap().
                with("name", t.getSimpleName()).
                with("package", readPackageName(t.getPackage()));
        }
        return typesData;
    }

    static private JsonList readMethods(CtModel model) {
        JsonList methodData = new JsonList();
        for (CtType<?> t : model.getAllTypes()) {
            for (CtMethod<?> m : t.getMethods()) {
                methodData.addMap().
                    with("name", m.getSimpleName()).
                    with("receiver", t.getSimpleName());
            }
        }
        return methodData;
    }
}
