package abstractor.core;

import java.io.File;
import java.util.ArrayList;
import java.util.List;
import java.util.SortedSet;
import java.util.TreeSet;

import spoon.Launcher;
import spoon.MavenLauncher;
import spoon.reflect.*;
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.*;
import spoon.reflect.reference.*;

import abstractor.core.constructs.*;
import abstractor.core.log.*;

public class Abstractor {
    public final Logger log;
    public final Project proj;

    public Abstractor(Logger log, Project proj) {
        this.log = log;
        this.proj = proj;
    }

    /**
     * Reads a project containing a pom.xml maven file.
     * @param mavenProject The path to the project file. 
     */
    public void addMavenProject(String mavenProject) throws Exception {
        this.log.log("Reading " + mavenProject);
        final MavenLauncher launcher = new MavenLauncher(mavenProject,
            MavenLauncher.SOURCE_TYPE.APP_SOURCE);
        launcher.addInputResource(mavenProject);
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
    public void addClassFromSource(String ...sourceLines) throws Exception {
        String source = String.join("\n", sourceLines);
        this.addObjectDecl(Launcher.parseClass(source));
    }

    public void finish() throws Exception {
        this.log.log("Finishing");
        this.log.push();
        this.resolveAllReferences();
        this.clearAllReferences();
        this.log.pop();
        this.log.log("Finished");
    }

    private void resolveAllReferences() throws Exception  {
        for (DeclarationRef ref : this.proj.declRefs) {
            if (!ref.isResolved()) ref.setResolved(this.addDeclaration(ref.elem));
        }

        for (TypeDescRef ref : this.proj.typeDescRefs) {
            if (!ref.isResolved()) ref.setResolved(this.addTypeDeclaration(ref.elem));
        }

        for (TypeParamRef ref : this.proj.typeParamRefs) {
            if (!ref.isResolved()) {
            }
        }
    }

    private void clearAllReferences() {
        for (DeclarationRef ref : this.proj.declRefs)
            this.log.errorIf(!ref.isResolved(), "Unresolved decl reference: " + ref.toString());
        for (TypeDescRef ref : this.proj.typeDescRefs)
            this.log.errorIf(!ref.isResolved(), "Unresolved type desc reference: " + ref.toString());
        for (TypeParamRef ref : this.proj.typeParamRefs)
            this.log.errorIf(!ref.isResolved(), "Unresolved type param reference: " + ref.toString());
        this.proj.declRefs.clear();
        this.proj.typeDescRefs.clear();
        this.proj.typeParamRefs.clear();
    }

    private <T extends Construct> void resolveReferencesFor(CtElement elem, T con) {
        if (con instanceof Reference<?>) return;
        if (con instanceof Declaration decl) {
            final DeclarationRef ref = this.proj.declRefs.get(elem);
            if (ref != null) ref.setResolved(decl);
        }
        if (con instanceof TypeDesc td) {
            final TypeDescRef ref = this.proj.typeDescRefs.get(elem);
            if (ref != null) ref.setResolved(td);
        }
        if (con instanceof TypeParam tp) {
            final TypeParamRef ref = this.proj.typeParamRefs.get(elem);
            if (ref != null) ref.setResolved(tp);
        }
    }

    public interface ConstructCreator<T extends Construct> { T create() throws Exception; }
    public interface FinishConstruct<T extends Construct> { void finish(T con) throws Exception; }
    public interface InProgressHandle<T extends Construct> { T handle() throws Exception; }

    public <T extends Construct, U extends T> T create(Factory<U> factory, CtElement elem,
        String title, ConstructCreator<U> c, FinishConstruct<U> f, InProgressHandle<T> h) throws Exception {
        final T existing = factory.get(elem);
        if (existing != null) return existing;

        try {
            log.log("Adding " + title);
            log.push();

            if (factory.inProgress) {
                if (h != null) {
                    log.log("Handling in progress: " + title);
                    return h.handle();
                }
                throw new Exception("Already in progress: " + title);
            }

            factory.inProgress = true;
            final U newCon = c.create();
            if (newCon == null) {
                factory.inProgress = false;
                return null;
            }

            final U other = factory.get(newCon);
            if (other != null) {
                factory.addElemKey(elem, other);
                factory.inProgress = false;
                return other;
            }
            factory.add(elem, newCon);

            this.resolveReferencesFor(elem, newCon);
            if (f != null) f.finish(newCon);
            factory.inProgress = false;
            return newCon;
        } finally {
            log.pop();
        }
    }

    public <T extends Construct, U extends T> T create(Factory<U> factory, CtElement elem,
        String title, ConstructCreator<U> c, FinishConstruct<U> f) throws Exception {
        return this.create(factory, elem, title, c, f, null);
    }

    public <T extends Construct, U extends T> T create(Factory<U> factory, CtElement elem,
        String title, ConstructCreator<U> c) throws Exception{
        return this.create(factory, elem, title, c, null, null);
    }

    private void addModel(CtModel model) throws Exception {
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

    private PackageCon addPackage(CtPackage pkg) throws Exception {
        return this.create(this.proj.packages, pkg,
            "package " + pkg.getQualifiedName(),
            () -> {
                final String name = pkg.getQualifiedName();
                final String path = packagePath(pkg);
                return new PackageCon(name, path);
            },
            (PackageCon pkgCon) -> {
                for (CtType<?> t : pkg.getTypes()) this.addDeclaration(t);
            });
    }

    private Declaration addDeclaration(CtElement elem) throws Exception {
        if (elem instanceof CtTypeReference<?> tr) elem = tr.getTypeDeclaration();

        if (elem instanceof CtClass<?> c) return this.addObjectDecl(c);
        if (elem instanceof CtInterface<?> i) return this.addInterfaceDecl(i);
        this.log.error("Unhandled decl (" + elem.getClass().getName() + ") "+elem.getShortRepresentation());
        return null;
    }

    private TypeDeclaration addTypeDeclaration(CtElement elem) throws Exception {
        if (elem instanceof CtTypeReference<?> tr) elem = tr.getTypeDeclaration();

        if (elem instanceof CtClass<?> c) return this.addObjectDecl(c);
        if (elem instanceof CtInterface<?> i) return this.addInterfaceDecl(i);
        this.log.error("Unhandled type decl (" + elem.getClass().getName() + ") "+elem.getShortRepresentation());
        return null;
    }

    /**
     * Handles adding and processing classes, enums, and records.
     * @param c The class to process.
     */
    private TypeDeclaration addObjectDecl(CtClass<?> c) throws Exception {
        return this.create(this.proj.objectDecls, c,
            "object decl " + c.getQualifiedName(),
            () -> {
                final CtPackage       pkg        = c.getPackage();
                final PackageCon      pkgCon     = pkg == null ? null : this.addPackage(pkg);
                final Location        loc        = proj.locations.create(c.getPosition());
                final String          name       = c.getSimpleName();
                final StructDesc      struct     = this.addStruct(c);
                final List<TypeParam> typeParams = this.addTypeParams(c.getFormalCtTypeParameters());
                return new ObjectDecl(pkgCon, loc, name, struct, typeParams);
            },
            (ObjectDecl obj) -> {
                obj.setVisibility(c);
                if (obj.pkg != null) obj.pkg.objectDecls.add(obj);

                //System.out.println("1) >>> " + c.getSuperclass());
                //System.out.println("2) >>> " + c.getSuperInterfaces());
                //System.out.println("3) >>> " + c.getConstructors());
                //System.out.println("4) >>> " + c.getNestedTypes());
                //System.out.println("5) >>> " + c.getTypeMembers());
                
                for (CtMethod<?> m : c.getAllMethods()) {
                    if (m.getParent() == c) this.addMethod(obj, m);
                }

                SortedSet<Abstract> abstracts = new TreeSet<Abstract>();
                for (CtMethod<?> m : c.getAllMethods()) {
                    abstracts.add(this.addAbstract(m));
                }
                obj.inter = this.proj.interfaceDescs.addOrGet(new InterfaceDesc(abstracts, obj));

                // TODO: Finish implementing
            },
            () -> {
                return this.addTypeDescRef(c.getReference());
            });
    }

    private MethodDecl addMethod(ObjectDecl receiver, CtMethod<?> m) throws Exception {
        return this.create(this.proj.methodDecls, m,
            "method " + m.getSignature(),
            () -> {
                final PackageCon      pkgCon     = receiver.pkg;
                final Location        loc        = proj.locations.create(m.getPosition());
                final String          name       = m.getSimpleName();
                final Signature       signature  = this.addSignature(m);
                final List<TypeParam> typeParams = this.addTypeParams(m.getFormalCtTypeParameters());
                return new MethodDecl(pkgCon, receiver, loc, name, signature, typeParams);
            },
            (MethodDecl md) -> {
                md.setVisibility(m);
                if (receiver.pkg != null) receiver.pkg.methodDecls.add(md);
                receiver.methodDecls.add(md);
                md.metrics = this.addMetrics(m);
            });
    }

    static private boolean isVoid(CtTypeReference<?> tr) {
        return tr.isPrimitive() && tr.getSimpleName().equals("void");
    }

    private Signature addSignature(CtMethod<?> m) throws Exception {
        return this.create(this.proj.signatures, m,
            "signature " + m.getSignature(),
            () -> {
                List<CtParameter<?>> params = m.getParameters();
                boolean variadic = params.size() > 0 && params.get(params.size()-1).isVarArgs();
                List<Argument> inArgs = new ArrayList<Argument>();
                for (CtParameter<?> p : params)
                    inArgs.add(this.addArgument(p));
        
                CtTypeReference<?> res = m.getType();
                List<Argument> outArgs = new ArrayList<Argument>();
                if (!isVoid(res)) outArgs.add(this.addArgument(res));

                return new Signature(variadic, inArgs, outArgs);
            });
    }

    private Argument addArgument(CtParameter<?> p) throws Exception {
        return this.create(this.proj.arguments, p,
            "parameter " + p.getSimpleName(),
            () -> {
                final String   name = p.getSimpleName();
                final TypeDesc type = this.addTypeDesc(p.getType());
                return new Argument(name, type);
            });
    }
    
    private Argument addArgument(CtTypeReference<?> t) throws Exception {
        return this.create(this.proj.arguments, t,
            "parameter <unnamed> " + t.getSimpleName(),
            () -> {
                final TypeDesc type = this.addTypeDesc(t);
                return new Argument("", type);
            });
    }

    private StructDesc addStruct(CtClass<?> c) throws Exception {
        return this.create(this.proj.structDescs, c,
            "struct " + c.getQualifiedName(),
            () -> {
                // Handle enum?
                //if (c instanceof CtEnum<?> e) {}
                ArrayList<Field> fields = new ArrayList<Field>();
                for (CtFieldReference<?> fr : c.getAllFields())
                    fields.add(this.addField(fr.getFieldDeclaration()));
                return new StructDesc(fields);
            });
    }

    private Field addField(CtField<?> f) throws Exception {
        return this.create(this.proj.fields, f,
            "field " + f.getSimpleName(),
            () -> {
                final String   name = f.getSimpleName();
                final TypeDesc type = this.addTypeDesc(f.getType());
                return new Field(name, type);
            },
            (Field field) -> {
                field.setVisibility(f);
            });
    }
    
    private TypeDeclaration addInterfaceDecl(CtInterface<?> i) throws Exception {
        return this.create(this.proj.interfaceDecls, i,
            "interface decl " + i.getQualifiedName(),
            () -> {
                final CtPackage       pkg        = i.getPackage();
                final PackageCon      pkgCon     = pkg == null ? null : this.addPackage(pkg);
                final Location        loc        = proj.locations.create(i.getPosition());
                final String          name       = i.getSimpleName();
                final InterfaceDesc   inter      = this.addInterfaceDesc(i);
                final List<TypeParam> typeParams = this.addTypeParams(i.getFormalCtTypeParameters());
                return new InterfaceDecl(pkgCon, loc, name, inter, typeParams);
            },
            (InterfaceDecl id) -> {
                id.setVisibility(i);
                if (id.pkg != null) id.pkg.interfaceDecls.add(id);                
            },
            () -> {
                return this.addTypeDescRef(i.getReference());
            });
    }

    private TypeDesc addTypeDesc(CtTypeReference<?> tr) throws Exception {
        if (tr.isPrimitive()) return this.addBasic(tr);
        if (tr.isArray())     return this.addArray((CtArrayTypeReference<?>)tr);

        CtType<?> ty = tr.getDeclaration();
        if (ty == null)       return this.proj.baker.objectDesc();
        if (tr.isClass())     return this.addObjectDecl((CtClass<?>)ty);
        if (tr.isInterface()) return this.addInterfaceDecl((CtInterface<?>)ty);
        if (tr.isGenerics())  return this.addTypeParam((CtTypeParameter)ty);

        // TODO: Finish implementing.
        return this.unknownTypeDesc(tr);
    }

    private TypeDesc unknownTypeDesc(CtTypeReference<?> tr) throws Exception {
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

    private TypeDeclaration addTypeDescRef(CtTypeReference<?> tr) throws Exception {
        return this.create(this.proj.typeDescRefs, tr,
            "type decl ref "+ tr.getSimpleName(),
            () -> {
                final String name = tr.getSimpleName();
                final String pkgPath = tr.getPackage().toString();
                final List<TypeDesc> tps = this.addTypeArguments(tr.getActualTypeArguments());
                return new TypeDescRef(tr, pkgPath, name, tps);
            });
    }

    private InterfaceInst addArray(CtArrayTypeReference<?> tr) throws Exception {
        return this.create(this.proj.interfaceInsts, tr,
            "array instance " + tr.getSimpleName(),
            () -> {
                final TypeDesc elem = this.addTypeDesc(tr.getArrayType());
                return this.proj.baker.arrayInst(tr.getQualifiedName(), elem);
            },
            (InterfaceInst inst) -> {
                inst.generic.instances.add(inst);
            });
    }
    
    private Basic addBasic(CtTypeReference<?> tr) throws Exception {
        return this.create(this.proj.basics, tr,
            "basic " + tr.getSimpleName(),
            () -> {
                final String name = tr.getSimpleName();
                this.log.errorIf(name == "void", "A void was added as a basic.");
                return new Basic(name);
            });
    }

    private InterfaceDesc addInterfaceDesc(CtInterface<?> i) throws Exception {
        return this.create(this.proj.interfaceDescs, i,
            "interface description " + i.getSimpleName(),
            () -> {
                final SortedSet<Abstract> abstracts = new TreeSet<Abstract>();
                for (CtMethod<?> m : i.getAllMethods())
                    abstracts.add(this.addAbstract(m));

                // TODO: Determine how to pin this interface.
                return new InterfaceDesc(abstracts);
            },
            (InterfaceDesc id) -> {
                // TODO: Implement Inheritance
            });
    }

    private Abstract addAbstract(CtMethod<?> m) throws Exception {
        return this.create(this.proj.abstracts, m,
            "abstract " + m.getSimpleName(),
            () -> {
                final String name = m.getSimpleName();
                final Signature signature = this.addSignature(m);
                return new Abstract(name, signature);
            });
    }

    private List<TypeDesc> addTypeArguments(List<CtTypeReference<?>> trs) throws Exception {
        List<TypeDesc> result = new ArrayList<TypeDesc>(trs.size());
        for (CtTypeReference<?> tr : trs) result.add(this.addTypeDesc(tr));
        return result;
    }

    private List<TypeParam> addTypeParams(List<CtTypeParameter> tps) throws Exception {
        List<TypeParam> result = new ArrayList<TypeParam>(tps.size());
        for (CtTypeParameter tp : tps) result.add(this.addTypeParam(tp));
        return result;
    }

    private TypeParam addTypeParam(CtTypeParameter tp) throws Exception {
        return this.create(this.proj.typeParams, tp,
            "type params " + tp.getQualifiedName(),
            () -> {
                final String name = tp.getQualifiedName();
                
                // TODO: Remove
                //System.out.println(">> " + name + " >> " + tp.prettyprint());
                //System.out.println(">> >> " + tp.getTypeErasure());
                //for (CtTypeReference<?> tpr : tp.getSuperInterfaces())
                //    System.out.println(">>  >> " + tpr.getSimpleName() + " >> " + tpr.prettyprint());

                CtTypeReference<?> tr = tp.getTypeErasure();
                final TypeDesc type = tr == null ?
                    this.proj.baker.objectDesc() :
                    this.addTypeDesc(tr);

                // TODO: Finish
                return new TypeParam(name, type);
            });
    }

    private Metrics addMetrics(CtMethod<?> m) throws Exception {
        return this.create(this.proj.metrics, m,
            "metrics",
            () -> {
                final Location loc = proj.locations.create(m.getPosition());
                final Analyzer ana = new Analyzer(this, loc);
                ana.addMethod(m);
                return ana.getMetrics();
            });
    }
}
