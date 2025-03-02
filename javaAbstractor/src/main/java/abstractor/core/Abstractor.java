package abstractor.core;

import java.io.File;
import java.util.ArrayList;
import java.util.List;

import spoon.Launcher;
import spoon.MavenLauncher;
import spoon.reflect.*;
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.*;
import spoon.reflect.reference.*;

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
        for (CtPackage pkg : model.getAllPackages())
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
        final Location loc = proj.locations.create(c.getPosition());
        final String name = c.getSimpleName();
        final StructDesc struct = this.addStruct(c);
        final List<TypeParam> typeParams = this.addTypeParams(c.getFormalCtTypeParameters());

        final ObjectDecl obj = new ObjectDecl(c, pkgCon, loc, name, struct, typeParams);
        final TryAddResult<ObjectDecl> prior = this.proj.objectDecls.tryAdd(obj);
        if (prior.existed) return prior.value;
        
        if (pkgCon != null) pkgCon.objectDecls.add(obj);

        //System.out.println("1) >>> " + c.getSuperclass());
        //System.out.println("2) >>> " + c.getSuperInterfaces());
        //System.out.println("3) >>> " + c.getConstructors());
        //System.out.println("4) >>> " + c.getNestedTypes());
        //System.out.println("5) >>> " + c.getTypeMembers());

        for (CtMethod<?> m : c.getAllMethods()) {
            if (m.getParent() == c) this.addMethod(obj, m);
        }

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
        final Location loc = proj.locations.create(m.getPosition());
        final String name = m.getSimpleName();

        // TODO: src.getBody()

        final Signature signature = null; // TODO: Finish
        final List<TypeParam> typeParams = this.addTypeParams(m.getFormalCtTypeParameters());

        MethodDecl md = new MethodDecl(m, pkgCon, receiver, loc, name, signature, typeParams);
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
        for (CtFieldReference<?> fr : c.getAllFields())
            fields.add(this.addField(fr.getFieldDeclaration()));

        StructDesc sd = new StructDesc(c, fields);
        final TryAddResult<StructDesc> prior = this.proj.structDescs.tryAdd(sd);
        if (prior.existed) return prior.value;

        return sd;
    }

    private Field addField(CtField<?> f) {
        Field existing = proj.fields.findWithSource(f);
        if (existing != null) return existing;

        final String name = f.getSimpleName();
        final TypeDesc type = this.addTypeDesc(f.getType());

        Field field = new Field(f, name, type);
        final TryAddResult<Field> prior = this.proj.fields.tryAdd(field);
        if (prior.existed) return prior.value;
        return field;
    }
    
    private InterfaceDecl addInterface(CtInterface<?> i) {
        final InterfaceDecl existing = proj.interfaceDecls.findWithSource(i);
        if (existing != null) return existing;

        this.log.log("Adding interface " + i.getQualifiedName());
        this.log.push();
        
        final CtPackage pkg = i.getPackage();
        final PackageCon pkgCon = pkg == null ? null : this.addPackage(pkg);
        final Location loc = proj.locations.create(i.getPosition());
        final String name = i.getSimpleName();

        final InterfaceDesc inter = null; // TODO: Finish
        final List<TypeParam> typeParams = this.addTypeParams(i.getFormalCtTypeParameters());

        InterfaceDecl id = new InterfaceDecl(i, pkgCon, loc, name, inter, typeParams);
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

    private TypeDesc addTypeDesc(CtTypeReference<?> tr) {
        if (tr.isPrimitive()) return this.addBasic(tr);
        //if (tr.isArray())     return this.addArray(tr); // TODO: Add once we have interfaceDesc.

        this.log.error("Unhandled (" + tr.getClass().getName() + "): "+tr.prettyprint());
        this.log.push();
        this.log.log("isAnnotationType:    " + tr.isAnnotationType());
        this.log.log("isAnonymous:. . . . ." + tr.isAnonymous());
        this.log.log("isArray:             " + tr.isArray());
        this.log.log("isClass:. . . . . . ." + tr.isClass());
        this.log.log("isEnum:              " + tr.isEnum());
        this.log.log("isGenerics: . . . . ." + tr.isGenerics());
        this.log.log("isImplicit:          " + tr.isImplicit());
        this.log.log("isInterface:. . . . ." + tr.isInterface());
        this.log.log("isLocalType:         " + tr.isLocalType());
        this.log.log("isParameterized:. . ." + tr.isParameterized());
        this.log.log("isParentInitialized: " + tr.isParentInitialized());
        this.log.log("isPrimitive:. . . . ." + tr.isPrimitive());
        this.log.log("isShadow:            " + tr.isShadow());
        this.log.log("isSimplyQualified:. ." + tr.isSimplyQualified());
        this.log.pop();
        return null;
    }

    private Basic addBasic(CtTypeReference<?> tr) {
        Basic existing = proj.basics.findWithSource(tr);
        if (existing != null) return existing;
        String name = tr.getSimpleName();
        Basic b = new Basic(tr, name);
        final TryAddResult<Basic> prior = this.proj.basics.tryAdd(b);
        if (prior.existed) return prior.value;
        return b;
    }

    private TypeParam addTypeParam(CtTypeParameter tp) {
        final String name = tp.getQualifiedName();
        
        //System.out.println(">> " + name + " >> " + tp.prettyprint());

        final TypeDesc type = null;

        // TODO: Finish
        return new TypeParam(tp, name, type);
    }

    private List<TypeParam> addTypeParams(List<CtTypeParameter> tps) {
        List<TypeParam> result = new ArrayList<TypeParam>();
        for (CtTypeParameter tp : tps) result.add(this.addTypeParam(tp));
        return result;
    }
}
