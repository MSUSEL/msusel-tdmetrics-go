package janis;

import json.JsonObj;
import spoon.Launcher;
import spoon.reflect.CtModel;

public class Janis {
    private Package packages;

    public Janis() {
        this.packages = null;
    }

    public void read(String path) {
        Launcher sl = new Launcher();
        sl.addInputResource(path);
        sl.getEnvironment().setNoClasspath(true);
        sl.getEnvironment().setComplianceLevel(7);

        sl.buildModel();
        CtModel model = sl.getModel();
        this.readPackages(model);
    }

    private void readPackages(CtModel model) {
        this.packages = new Package(model.getRootPackage());
    }

    public void write() {
        String str = JsonObj.toString(this.packages.toJson());
        System.out.println(str);
    }
}
