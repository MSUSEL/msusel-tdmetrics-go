package abstractor.core;

import spoon.Launcher;
import spoon.MavenLauncher;
import spoon.reflect.CtModel;
import spoon.reflect.declaration.*;

import abstractor.core.log.*;

public class Abstractor {
    private final Logger log;
    private final Project proj;

    public Abstractor(Logger log, Project proj) {
        this.log = log;
        this.proj = proj;
    }

    /**
     * Reads a project containing a pom.xml maven file.
     * @param mavenProject The path to the project file. 
     */
    public void addMavenProject(String mavenProject) {
        this.log.log("Reading " + mavenProject);
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
        if (this.proj.packages.contains(pkg)) return;
        this.log.log("Adding package " + pkg);
        this.log.push();
        this.proj.packages.add(pkg);
        for (CtType<?> t : pkg.getTypes())
            this.addType(t);
        this.log.pop();
    }

    private void addType(CtType<?> t) {
        if (t instanceof CtEnum<?> e) this.addEnum(e);
        else if (t instanceof CtClass<?> c) this.addClass(c);
        else if (t instanceof CtInterface<?> i) this.addInterface(i);
        else this.log.error("Unhandled (" + t.getClass().getName() + ") "+t.getQualifiedName());
    }
    
    private void addEnum(CtEnum<?> e) {
        this.log.log("Adding enum " + e.getQualifiedName());
        
        // TODO: Implement

        this.log.error("enum unimplemented");
    }

    /**
     * Handles adding and processing classes and records.
     * @param c The class to process.
     */
    private void addClass(CtClass<?> c) {
        if (this.proj.objects.contains(c)) return;
        this.log.log("Adding class " + c.getQualifiedName());
        this.log.push();
        this.proj.objects.add(c);

        // TODO: Implement
        
        this.log.pop();
    }
    
    private void addInterface(CtInterface<?> i) {
        if (this.proj.interfaceDecls.contains(i)) return;
        this.log.log("Adding interface " + i.getQualifiedName());
        this.log.push();
        this.proj.interfaceDecls.add(i);
        
        // TODO: Implement

        this.log.pop();
    }
}
