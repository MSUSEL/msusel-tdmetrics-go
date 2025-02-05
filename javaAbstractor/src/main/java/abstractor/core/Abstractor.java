package abstractor.core;

import spoon.Launcher;
import spoon.MavenLauncher;
import spoon.reflect.CtModel;
import spoon.reflect.declaration.*;

import abstractor.core.json.*;
import abstractor.core.log.*;

public class Abstractor {
    private final Logger out;

    public Abstractor(Logger out) {
        this.out = out;
    }

    /**
     * Reads a project containing a pom.xml maven file.
     * @param mavenProject The path to the project file. 
     */
    public void addMavenProject(String mavenProject) {
        this.out.log("Reading " + mavenProject);
        final MavenLauncher launcher = new MavenLauncher(mavenProject,
                 MavenLauncher.SOURCE_TYPE.APP_SOURCE);
        this.addModel(launcher.buildModel());
    }

    /**
     * Parses the source for a given class and adds it.
     * 
     * This is designed to test classes, records, and enumerators,
     * but will not work for interfaces.
     * @example parseClass("class C { void m() { System.out.println(\"hello\"); } }"); 
     * @param source The class source code.
     */
    public void addClassFromSource(String source) {
        this.addClass(Launcher.parseClass(source));
    }

    private void addModel(CtModel model) {
        for(CtPackage pkg : model.getAllPackages())
            this.addPackage(pkg);
    }

    private void addPackage(CtPackage pkg) {
        this.out.log("Adding package " + pkg);
        this.out.push();

        //TODO: Store package for later

        for (CtType<?> t : pkg.getTypes())
            this.addType(t);
        this.out.pop();
    }

    private void addType(CtType<?> t) {
        if (t instanceof CtEnum<?> e) this.addEnum(e);
        else if (t instanceof CtClass<?> c) this.addClass(c);
        else if (t instanceof CtInterface<?> i) this.addInterface(i);
        else this.out.log("Unhandled (" + t.getClass().getName() + ") "+t.getQualifiedName());
    }
    
    private void addEnum(CtEnum<?> e) {
        this.out.log("Adding enum " + e.getQualifiedName());
        // TODO: Implement
    }

    /**
     * Handles adding and processing classes and records.
     * @param c The class to process.
     */
    private void addClass(CtClass<?> c) {
        this.out.log("Adding class " + c.getQualifiedName());
        // TODO: Implement
    }
    
    private void addInterface(CtInterface<?> i) {
        this.out.log("Adding interface " + i.getQualifiedName());
        // TODO: Implement
    }

    public JsonNode toJson(boolean writeTypes, boolean writeIndices) {
        JsonObject obj = new JsonObject();
        obj.put("language", JsonValue.of("java"));




        return obj;
    }
}
