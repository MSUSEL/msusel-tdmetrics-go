package abstractor.core;

import spoon.Launcher;
import spoon.MavenLauncher;
import spoon.reflect.*;
import spoon.reflect.declaration.*;
import spoon.reflect.reference.CtFieldReference;
import abstractor.core.constructs.*;
import abstractor.core.constructs.Package;
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
        if (this.proj.packages.containsSource(pkg)) return;
        this.log.log("Adding package " + pkg);
        this.log.push();
        Package.Create(this.proj, pkg);
        for (CtType<?> t : pkg.getTypes())
            this.addType(t);
        this.log.pop();
    }

    private void addType(CtType<?> t) {
        if (t instanceof CtClass<?> c) this.addClass(c);
        else if (t instanceof CtInterface<?> i) this.addInterface(i);
        else this.log.error("Unhandled (" + t.getClass().getName() + ") "+t.getQualifiedName());
    }

    /**
     * Handles adding and processing classes, enums, and records.
     * @param c The class to process.
     */
    private void addClass(CtClass<?> c) {
        if (this.proj.objectDecls.containsSource(c)) return;
        this.log.log("Adding class " + c.getQualifiedName());
        this.log.push();
        ObjectDecl obj = ObjectDecl.Create(this.proj, c);
        for (CtMethod<?> m : c.getAllMethods())
            this.addMethod(obj, m);

        //for (CtFieldReference<?> fr : c.getAllFields()) {
	    //    CtField<?> f = fr.getFieldDeclaration();
        //    f.
        //}

        // TODO: Implement
        
        this.log.pop();
    }

    private void addMethod(ObjectDecl receiver, CtMethod<?> m) {
        if (this.proj.methodDecls.containsSource(m)) return;
        this.log.log("Adding method " + m.prettyprint());
        this.log.push();
        MethodDecl.Create(proj, receiver, m);

        // TODO: Implement
        
        this.log.pop();
    }
    
    private void addInterface(CtInterface<?> i) {
        if (this.proj.interfaceDecls.containsSource(i)) return;
        this.log.log("Adding interface " + i.getQualifiedName());
        this.log.push();
        InterfaceDecl.Create(proj, i);
        //for (CtMethod<?> m : i.getAllMethods())
        //    this.addAbstract(m);
        
        // TODO: Implement

        this.log.pop();
    }

    /*
    private void addAbstract(CtMethod<?> m) {
        if (this.proj.abstracts.containsSource(m)) return;
        this.log.log("Adding abstract " + m.prettyprint());
        this.log.push();
        this.proj.abstracts.add(m);

        // TODO: Implement
        
        this.log.pop();
    }
    */
}
