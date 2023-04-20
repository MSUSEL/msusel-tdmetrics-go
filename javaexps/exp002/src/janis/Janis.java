package janis;

import json.JsonList;
import json.JsonMap;
import json.JsonObj;
import json.JsonStr;
import spoon.Launcher;
import spoon.reflect.CtModel;
import spoon.reflect.declaration.*;
import spoon.reflect.reference.CtTypeReference;
import java.util.List;

/**
 * Janis is a Java Analysis tool to produce a JSON description of a module.
 * See <a href="https://spoon.gforge.inria.fr/structural_elements.html">Structural Elements</a>
 * and <a href="https://spoon.gforge.inria.fr/code_elements.html">Code Elements</a>
 */
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

    static private JsonObj readPackageName(CtPackage p) {
        return new JsonStr(p.isUnnamedPackage() ? "<root>" : p.getSimpleName());
    }

    static private JsonObj readType(CtTypeReference ref) {
        return new JsonStr(ref.getSimpleName());
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

    static private JsonObj readParameters(CtMethod<?> m) {
        JsonList results = new JsonList();
        List<CtParameter<?>> params = m.getParameters();
        for (CtParameter<?> param : params) {
            JsonList types = results.addList();
            for (CtTypeReference<?> ref : param.getReferencedTypes())
                types.with(ref.getSimpleName());
        }
        return results;
    }

    static private JsonObj readReturns(CtMethod<?> m) {
        return new JsonStr(m.getType().getSimpleName());
    }

    static private JsonList readMethods(CtModel model) {
        JsonList methodData = new JsonList();
        for (CtType<?> t : model.getAllTypes()) {
            for (CtMethod<?> m : t.getMethods()) {
                methodData.addMap().
                    with("name", m.getSimpleName()).
                    with("receiver", t.getSimpleName()).
                    withOmitOnDefault("parameters", readParameters(m)).
                    withOmitOnDefault("returns", readReturns(m));

                // TODO: Add cyclomatic complexity
            }
        }
        return methodData;
    }
}
