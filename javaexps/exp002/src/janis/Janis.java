package janis;

import json.JsonMap;
import json.JsonObj;
import json.Jsonable;
import spoon.Launcher;
import spoon.reflect.CtModel;
import spoon.reflect.declaration.CtPackage;
import spoon.reflect.declaration.CtType;

import java.util.ArrayList;

public class Janis implements Jsonable {
    private ArrayList<JPackage> packages;
    private ArrayList<JType> types;
    private ArrayList<JMethod> methods;

    public Janis() {
        this.packages = new ArrayList<>();
        this.types = new ArrayList<>();
        this.methods = new ArrayList<>();
    }

    public void read(String path) {
        Launcher sl = new Launcher();
        sl.addInputResource(path);
        sl.getEnvironment().setNoClasspath(true);
        sl.getEnvironment().setComplianceLevel(7);

        sl.buildModel();
        CtModel model = sl.getModel();
        this.readPackages(model);
        this.readTypes(model);
        this.readMethods(model);
    }

    private void readPackages(CtModel model) {
        this.packages.clear();
        for (CtPackage p: model.getAllPackages())
            this.packages.add(new JPackage(p));
    }

    private void readTypes(CtModel model) {
        this.types.clear();
        for (CtType t: model.getAllTypes()) {
            this.types.add(new JType(t));
        }

    }

    private void readMethods(CtModel model) {
        this.methods.clear();

    }

    @Override
    public JsonObj toJson() {
        return new JsonMap().
            withOmitEmpty("packages", this.packages).
            withOmitEmpty("types", this.types).
            withOmitEmpty("methods", this.methods);
    }

    public void write() {
        String results = JsonObj.toString(this.toJson());
        System.out.println(results);
    }
}
