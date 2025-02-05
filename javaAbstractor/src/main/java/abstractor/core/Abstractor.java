package abstractor.core;

import java.io.PrintStream;

import spoon.Launcher;
import spoon.MavenLauncher;
import spoon.reflect.CtModel;
import spoon.reflect.declaration.*;

import abstractor.core.json.*;

public class Abstractor {
    private final PrintStream out;

    public Abstractor(PrintStream out) {
        this.out = out;
    }

    private void log(String text) {
        if (this.out != null) this.out.println(text);
    }

    /**
     * Reads a project containing a pom.xml maven file.
     * @param mavenProject The path to the project file. 
     */
    public void addMavenProject(String mavenProject) {
        this.log("Reading " + mavenProject + "...");
        final MavenLauncher launcher = new MavenLauncher(mavenProject,
                 MavenLauncher.SOURCE_TYPE.APP_SOURCE);
        CtModel model = launcher.buildModel();
        this.log("Done reading.");

        this.addModel(model);
    }

    private void addModel(CtModel model) {
        for(CtPackage pkg : model.getAllPackages())
            this.addPackage(pkg);
    }

    private void addPackage(CtPackage pkg) {
        this.log("Adding package " + pkg);

        //TODO: Store package for later

        for (CtType<?> t : pkg.getTypes())
            this.addType(t);
    }

    private void addType(CtType<?> t) {
        if (t instanceof CtClass<?> c)
            this.addClass(c);
        else
            this.log("Unhandled "+t);
    }

    /**
     * Parses the source for a given class and adds it.
     * @example parseClass("class C { void m() { System.out.println(\"hello\"); } }"); 
     * @param source The class source code.
     */
    public void addClassFromSource(String source) {
        this.addClass(Launcher.parseClass(source));
    }

    private void addClass(CtClass<?> c) {
        this.log("Adding class " + c.getQualifiedName());
    }

    public JsonNode toJson(boolean writeTypes, boolean writeIndices) {
        JsonObject obj = new JsonObject();
        obj.put("language", JsonValue.of("java"));




        return obj;
    }
}
