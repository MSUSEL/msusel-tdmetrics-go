package abstractor.core;

import spoon.Launcher;
import spoon.MavenLauncher;
import spoon.reflect.*;
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.*;
import spoon.reflect.reference.*;

import java.io.File;
import java.util.ArrayList;

import abstractor.core.constructs.*;
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
    public void addClassFromSource(String ...sourceLines) {
        String source = String.join("\n", sourceLines);
        this.addClass(Launcher.parseClass(source));
    }

    private void addModel(CtModel model) {
        for(CtPackage pkg : model.getAllPackages())
            this.addPackage(pkg);
    }

    static private String packagePath(CtPackage p) {
        SourcePosition pos = p.getPosition();
        if (!pos.isValidPosition()) return "";
        
        final File file = pos.getFile();
        if (file == null) return "";

        final String path = file.getPath();
        final String tail = "package-info.java";
        if (!path.endsWith(tail)) return path;
        return path.substring(0, path.length()-tail.length());
    }

    private PackageCon addPackage(CtPackage pkg) {
        final PackageCon existing = proj.packages.findWithSource(pkg);
        if (existing != null) return existing;

        this.log.log("Adding package " + pkg);
        this.log.push();

        final String name = pkg.getQualifiedName();
        final String path = packagePath(pkg);
        final PackageCon pkgCon = new PackageCon(pkg, name, path);
        final TryAddResult<PackageCon> prior = this.proj.packages.tryAdd(pkgCon);
        if (prior.existed) return prior.value;

        for (CtType<?> t : pkg.getTypes()) {
            if (t instanceof CtClass<?> c) this.addClass(c);
            else if (t instanceof CtInterface<?> i) this.addInterface(i);
            else this.log.error("Unhandled (" + t.getClass().getName() + ") "+t.getQualifiedName());
        }

        this.log.pop();
        return pkgCon;
    }

    /**
     * Handles adding and processing classes, enums, and records.
     * @param c The class to process.
     */
    private ObjectDecl addClass(CtClass<?> c) {
        final ObjectDecl existing = proj.objectDecls.findWithSource(c);
        if (existing != null) return existing;

        this.log.log("Adding class " + c.getQualifiedName());
        this.log.push();

        final CtPackage pkg = c.getPackage();
        final PackageCon pkgCon = pkg == null ? null : this.addPackage(pkg);
        final Locations.Location loc = proj.locations.create(c.getPosition());
        final String name = c.getSimpleName();
        final StructDesc struct = this.addStruct(c);

        // TODO: Handle enum?
        //if (c instanceof CtEnum<?> e) {}

        final ObjectDecl obj = new ObjectDecl(c, struct, pkgCon, loc, name);
        final TryAddResult<ObjectDecl> prior = this.proj.objectDecls.tryAdd(obj);
        if (prior.existed) return prior.value;
        
        if (pkgCon != null) pkgCon.objectDecls.add(obj);

        for (CtMethod<?> m : c.getAllMethods())
            this.addMethod(obj, m);

        // TODO: Deal with struct

        //for (CtFieldReference<?> fr : c.getAllFields())
        //    this.addField(fr.getFieldDeclaration());

        // TODO: Implement
        
        this.log.pop();
        return obj;
    }

    private MethodDecl addMethod(ObjectDecl receiver, CtMethod<?> m) {
        final MethodDecl existing = proj.methodDecls.findWithSource(m);
        if (existing != null) return existing;

        this.log.log("Adding method " + m.getSimpleName());
        this.log.push();

        final PackageCon pkgCon = receiver.pkg;
        final Locations.Location loc = proj.locations.create(m.getPosition());
        final String name = m.getSimpleName();
        // TODO: src.getBody()
        // TODO: src.getFormalCtTypeParameters()

        MethodDecl md = new MethodDecl(m, pkgCon, receiver, loc, name);
        final TryAddResult<MethodDecl> prior = this.proj.methodDecls.tryAdd(md);
        if (prior.existed) return prior.value;

        // TODO: Implement
        if (pkgCon != null) pkgCon.methodDecls.add(md);
        receiver.methodDecls.add(md);

        this.log.pop();
        return md;
    }

    private StructDesc addStruct(CtClass<?> c) {
        StructDesc existing = proj.structDescs.findWithSource(c);
        if (existing != null) return existing;
        
        this.log.log("Adding struct " + c.getSimpleName());
        this.log.push();

        // TODO: Handle enum?
        //if (c instanceof CtEnum<?> e) {}
        
        ArrayList<Field> fields = new ArrayList<Field>();

        StructDesc sd = new StructDesc(c, fields);
        final TryAddResult<StructDesc> prior = this.proj.structDescs.tryAdd(sd);
        if (prior.existed) return prior.value;

        return sd;
    }

    private Field addField(CtField<?> src) {
        Field existing = proj.fields.findWithSource(src);
        if (existing != null) return existing;

        // TODO: Get initial stuff

        Field f = new Field(src);
        final TryAddResult<Field> prior = this.proj.fields.tryAdd(f);
        if (prior.existed) return prior.value;
        
        // TODO: Finish loading 

        return f;
    }
    
    private InterfaceDecl addInterface(CtInterface<?> i) {
        final InterfaceDecl existing = proj.interfaceDecls.findWithSource(i);
        if (existing != null) return existing;

        this.log.log("Adding interface " + i.getQualifiedName());
        this.log.push();
        
        final CtPackage pkg = i.getPackage();
        final PackageCon pkgCon = pkg == null ? null : this.addPackage(pkg);
        final Locations.Location loc = proj.locations.create(i.getPosition());
        final String name = i.getSimpleName();

        InterfaceDecl id = new InterfaceDecl(i, pkgCon, loc, name);
        final TryAddResult<InterfaceDecl> prior = this.proj.interfaceDecls.tryAdd(id);
        if (prior.existed) return prior.value;

        //for (CtMethod<?> m : i.getAllMethods())
        //    this.addAbstract(m);
        
        // TODO: Implement
        if (pkgCon != null) pkgCon.interfaceDecls.add(id);

        this.log.pop();
        return id;
    }

    /*
    static public Basic Create(Project proj, CtField<?> src) {
        Basic existing = proj.basics.findWithSource(src);
        if (existing != null) return existing;

        String name = "TODO"; // TODO: Get name

        Basic f = new Basic(src, name);
        existing = proj.basics.tryAdd(f);
        if (existing != null) return existing;
        return f;
    }

    static public Argument Create(Project proj, CtField<?> src) {
        Argument existing = proj.arguments.findWithSource(src);
        if (existing != null) return existing;

        // TODO: Get initial stuff

        Argument f = new Argument(src);
        existing = proj.arguments.tryAdd(f);
        if (existing != null) return existing;
        
        // TODO: Finish loading 

        return f;
    }

    static public Abstract Create(Project proj, CtField<?> src) {
        Abstract existing = proj.abstracts.findWithSource(src);
        if (existing != null) return existing;

        // TODO: Get initial stuff

        Abstract f = new Abstract(src);
        existing = proj.abstracts.tryAdd(f);
        if (existing != null) return existing;
        
        // TODO: Finish loading 

        return f;
    }

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
