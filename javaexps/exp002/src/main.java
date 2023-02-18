import spoon.*;
import spoon.reflect.*;
import spoon.reflect.declaration.*;

class Main {

    public static void main(String args[]) {
        Launcher sl = new Launcher();
        sl.addInputResource("./testProjects/bcel/src/main/java/org/apache/bcel");
        sl.getEnvironment().setNoClasspath(true);
        sl.getEnvironment().setComplianceLevel(7);

        sl.buildModel();
        CtModel model = sl.getModel();

        for (CtType t : model.getAllTypes()) {
            System.out.println(t.getQualifiedName());
        }
    }
}
