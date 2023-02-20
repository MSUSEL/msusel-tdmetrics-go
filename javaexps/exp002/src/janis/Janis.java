package janis;

import json.JsonObj;
import spoon.Launcher;
import spoon.reflect.CtModel;

public class Janis {
    private JPackage rootPackage;

    public Janis() {
        this.rootPackage = null;
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
        this.rootPackage = new JPackage(model.getRootPackage());
    }

    public void write() {
        String results = JsonObj.toString(this.rootPackage);
        System.out.println(results);
    }
}
