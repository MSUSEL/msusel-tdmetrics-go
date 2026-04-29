package abstractor.core;

import java.io.File;
import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.TreeSet;
import java.util.HashSet;

import spoon.Launcher;
import spoon.MavenLauncher;
import spoon.reflect.*;
import spoon.reflect.code.*;
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.*;
import spoon.reflect.path.CtRole;
import spoon.reflect.reference.*;
import spoon.reflect.visitor.filter.TypeFilter;
import abstractor.core.constructs.*;
import abstractor.core.log.*;
import abstractor.core.validator.Validator;

public class Abstractor {
    public final Logger log;
    public final Project proj;

    public final HashSet<CtMethod<?>> pendingMetrics = new HashSet<CtMethod<?>>();

    public Abstractor(Logger log, Project proj) {
        this.log  = log;
        this.proj = proj;
    }

    /**
     * Reads a project containing a pom.xml maven file.
     * @param mavenProject The path to the project file. 
     */
    public void addMavenProject(String mavenProject) throws Exception {
        this.log.log("Reading " + mavenProject);
        MavenLauncher launcher = new MavenLauncher(mavenProject, MavenLauncher.SOURCE_TYPE.APP_SOURCE);
        CtModel model = launcher.buildModel();
        if (model.getAllTypes().size() > 0) {
            this.addModel(model);
            return;
        }

        // If the model couldn't be loaded (it has no types in it) from the app
        // source alone then try again with the maven project path as an input
        // resource. We can't always add the input resource otherwise it will
        // cause duplicate identifiers in some projects. For the integration
        // tests in testData/java, we do need the input resource. I have no clue
        // what the difference is between the maven models to require this
        // but if it works, I'm not going to fix it right now.
        launcher = new MavenLauncher(mavenProject, MavenLauncher.SOURCE_TYPE.APP_SOURCE);
        launcher.addInputResource(mavenProject);
        model = launcher.buildModel();
        this.addModel(model);
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

    private void addModel(CtModel model) throws Exception {
        for (CtPackage pkg : model.getAllPackages())
            this.addPackage(pkg);
    }

    static private String normalizePath(String path) {
        return path.replaceAll("\\\\", "/");
    }

    static private String packageName(CtPackage pkg) {
        if (pkg == null) return "<java.lang>";
        String name = pkg.getQualifiedName();
        return name.isBlank() ? "<unnamed>" : name;
    }

    static private String packagePath(CtPackage pkg) {
        if (pkg == null) return "";
        SourcePosition pos = pkg.getPosition();
        if (!pos.isValidPosition()) return "";
        
        final File file = pos.getFile();
        if (file == null) return "";

        final String path = normalizePath(file.getPath());
        final String tail = "/package-info.java";
        if (!path.endsWith(tail)) return path;
        return path.substring(0, path.length()-tail.length());
    }

    private List<Ref<PackageCon>> getImports(CtPackage pkg) throws Exception {

        // TODO: REWORK and FINISH
        System.out.println("> imports for " + pkg); // TODO: REMOVE
        
        for (CtType<?> type : pkg.getTypes()) {
            CtCompilationUnit cu = type.getPosition().getCompilationUnit();
            if (cu != null) {
                for (CtImport imp : cu.getImports()) {
                    System.out.println("  a> "+imp);
                }
            }
        }

        // also check package-info.java if there is one
        CtCompilationUnit pkgCu = pkg.getPosition().getCompilationUnit();
        if (pkgCu != null) {
            for (CtImport imp : pkgCu.getImports()) {
                System.out.println("  b> "+imp);
            }
        }

        // TODO: Printing call/executable/type information
        class Temp {
            void printTyp(CtType<?> typ) {
                for (CtMethod<?> method : typ.getMethods()) {
                    for (CtTypeAccess<?> typeAccess : method.getElements(new TypeFilter<>(CtTypeAccess.class))) {
                        CtTypeReference<?> tr = typeAccess.getType();
                        System.out.println("  c> type access: " + tr.getQualifiedName());
                    }

                    for (CtFieldAccess<?> fa : method.getElements(new TypeFilter<>(CtFieldAccess.class))) {
                        CtFieldReference<?> fr = fa.getVariable();
                        System.out.println("  d> field access: " + fr.getDeclaringType().getQualifiedName() + "." + fr.getSimpleName());
                    }

                    for (CtInvocation<?> inv : method.getElements(new TypeFilter<>(CtInvocation.class))) {
                        CtExecutableReference<?> exe = inv.getExecutable();
                        System.out.println("  e> invocation: " + exe.getDeclaringType().getQualifiedName() + "." + exe.getSimpleName());
                    }
                }
                
                for (CtType<?> nestedTyp : typ.getNestedTypes()) printTyp(nestedTyp);
            }
        }

        Temp tmp = new Temp();
        for (CtType<?> typ : pkg.getTypes()) tmp.printTyp(typ);
        System.out.println();
        return null;
    }

    public Ref<PackageCon> addPackage(CtPackage pkg) throws Exception {
        final String name = packageName(pkg);
        return this.proj.packages.create(this.log, pkg,
            "package " + name,
            () -> {
                final String path = packagePath(pkg);
                return new PackageCon(name, path);
            },
            (Ref<PackageCon> ref, PackageCon pkgCon) ->{
                for (CtType<?> t : pkg.getTypes())
                    this.addDeclarationToPackage(pkgCon, this.addDeclaration(t));

                // TODO: add Imports
                this.getImports(pkg);

                // TODO: Finish
            });
    }

    static private <T extends Construct> boolean tryToAdd(Set<Ref<T>> set, Ref<? extends Construct> e, ConstructKind kind) {
        if (e.kind() == kind) {
            @SuppressWarnings("unchecked")
            Ref<T> cast = (Ref<T>)e;
            set.add(cast);
            return true;
        }
        return false;
    }

    public void addDeclarationToPackage(PackageCon pkg, Ref<? extends Construct> decl) {
        if (tryToAdd(pkg.objectDecls, decl, ConstructKind.OBJECT_DECL)) return;
        if (tryToAdd(pkg.interfaceDecls, decl, ConstructKind.INTERFACE_DECL)) return;
        if (tryToAdd(pkg.methodDecls, decl, ConstructKind.METHOD_DECL)) return;
        if (tryToAdd(pkg.values, decl, ConstructKind.VALUE)) return;
        
        this.log.error("Unhandled declaration type: " + decl.kind());
    }

    public Ref<? extends Construct> addDeclaration(CtElement elem) throws Exception {
        if (elem instanceof CtTypeReference<?> tr) elem = tr.getTypeDeclaration();

        if (elem instanceof CtClass<?> c)     return this.addObjectDecl(c);
        if (elem instanceof CtInterface<?> i) return this.addInterfaceDecl(i);
        if (elem instanceof CtMethod<?> m)    return this.addGeneralMethod(m);
        this.log.error("Unhandled decl (" + elem.getClass().getName() + ") "+elem.toStringDebug());
        return null;
    }

    public Ref<? extends TypeDeclaration> addTypeDeclaration(CtElement elem) throws Exception {
        if (elem instanceof CtTypeReference<?> tr) elem = tr.getTypeDeclaration();

        if (elem instanceof CtClass<?> c)     return this.addObjectDecl(c);
        if (elem instanceof CtInterface<?> i) return this.addInterfaceDecl(i);
        this.log.error("Unhandled type decl (" + elem.getClass().getName() + ") "+elem.toStringDebug());
        return null;
    }

    static public boolean isObjectMethod(CtMethod<?> m) {
        if (m == null) return false;

        final CtTypeReference<?> objectRef = m.getFactory().Type().objectType();
        final CtType<?> objectDecl = objectRef.getTypeDeclaration();
        if (objectDecl == null) return false;

        final String sig = m.getSignature();
        for (CtMethod<?> objectMethod : objectDecl.getMethods()) {
            if (sig.equals(objectMethod.getSignature())) return true;
        }
        return false;
    }

    /**
     * Handles adding and processing classes, enums, and records.
     * @param c The class to process.
     */
    public Ref<ObjectDecl> addObjectDecl(CtClass<?> c) throws Exception {
        return this.proj.objectDecls.create(this.log, c,
            "object decl " + c.getQualifiedName(),
            () -> {
                final Ref<PackageCon>      pkg        = this.addPackage(c.getPackage());
                final Location             loc        = proj.locations.create(c.getPosition());
                final String               name       = c.getSimpleName();
                final Ref<StructDesc>      struct     = this.addStruct(c);
                final List<Ref<TypeParam>> typeParams = this.addTypeParams(c.getFormalCtTypeParameters());
                return new ObjectDecl(pkg, loc, name, struct, typeParams);
            },
            (Ref<ObjectDecl> ref, ObjectDecl obj) -> {
                obj.setVisibility(c);
                
                // Add constructors as (static) methods.
                for (CtConstructor<?> ctor : c.getConstructors()) {
                    if (ctor.getParent().equals(c)) this.addConstructorMethod(ref, ctor);
                }

                // Add methods for the class.
                for (CtMethod<?> m : c.getAllMethods()) {
                    if (m.getParent().equals(c) && !isObjectMethod(m))
                        this.addMethod(ref, m);
                }

                // Synthesize the interface description for the class.
                TreeSet<Ref<Abstract>> abstracts = new TreeSet<Ref<Abstract>>();
                for (CtMethod<?> m : c.getAllMethods()) {
                    if (!m.isStatic() && !isObjectMethod(m))
                        abstracts.add(this.addAbstract(m));
                }
                obj.inter = this.proj.interfaceDescs.addOrGetRef(new InterfaceDesc(abstracts, ref));

                // TODO: Finish implementing
                //System.out.println("1) >>> " + c.getSuperInterfaces());

                // Add any nested types.
                for (CtType<?> nt : c.getNestedTypes()) this.addTypeDesc(nt.getReference());
            });
    }

    public Ref<MethodDecl> addConstructorMethod(Ref<ObjectDecl> receiver, CtConstructor<?> ctor) throws Exception {
        if (!receiver.isResolved())
            throw new Exception("Expected the receiver for a constructor method to be resolved: " + receiver.toString());
        ObjectDecl recv = receiver.getResolved();

        return this.proj.methodDecls.create(log, ctor,
            "constructor " + ctor.getSignature(),
            () -> {
                final Ref<PackageCon>      pkg        = recv.pkg;
                final Location             loc        = proj.locations.create(ctor.getPosition());
                final String               name       = recv.name;
                final Ref<Signature>       signature  = this.addConstructSignature(ctor);
                final List<Ref<TypeParam>> typeParams = this.addTypeParams(ctor.getFormalCtTypeParameters());
                return new MethodDecl(pkg, receiver, loc, name, signature, typeParams);
            },
            (Ref<MethodDecl> ref, MethodDecl md) -> {
                md.setVisibility(ctor);
                recv.methodDecls.add(ref);
            });
    }

    public Ref<? extends Construct> addGeneralMethod(CtMethod<?> m) throws Exception {
        CtType<?> decl = m.getDeclaringType();
        if (decl instanceof CtClass<?> c) {
            Ref<ObjectDecl> obj = this.addObjectDecl(c);
            return this.addMethod(obj, m);
        }
        if (decl instanceof CtInterface<?> i) {
            //Ref<InterfaceDecl> it = this.addInterfaceDecl(i);
            Ref<Abstract> ab = this.addAbstract(m);

            // TODO: Finish interface
            return ab;
        }
        throw new Exception("Unhandled general method declaring type (" + decl.getClass().getName() + ") "+decl.getQualifiedName());
    }

    public Ref<MethodDecl> addMethod(Ref<ObjectDecl> receiver, CtMethod<?> m) throws Exception {
        if (!receiver.isResolved())
            throw new Exception("Expected the object receiver for a method to be resolved: " + receiver.toString());
        ObjectDecl recv = receiver.getResolved();

        return this.proj.methodDecls.create(this.log, m,
            "method " + m.getSignature(),
            () -> {
                final Ref<PackageCon>      pkg        = recv.pkg;
                final Location             loc        = proj.locations.create(m.getPosition());
                final String               name       = m.getSimpleName();
                final Ref<Signature>       signature  = this.addSignature(m);
                final List<Ref<TypeParam>> typeParams = this.addTypeParams(m.getFormalCtTypeParameters());
                return new MethodDecl(pkg, receiver, loc, name, signature, typeParams);
            },
            (Ref<MethodDecl> ref, MethodDecl md) -> {
                md.setVisibility(m);
                //if (pkg != null) pkg.methodDecls.add(md); // TODO: Move to a follow up when we know the package is done.
                recv.methodDecls.add(ref);
                this.pendingMetrics.add(m);
            });
    }

    static public boolean isVoid(CtTypeReference<?> tr) {
        return tr.isPrimitive() && tr.getSimpleName().equals("void");
    }

    public Ref<Signature> addConstructSignature(CtConstructor<?> m) throws Exception {
        return this.proj.signatures.create(this.log, m,
            "constructor signature " + m.getSignature(),
            () -> {
                final List<CtParameter<?>> ps = m.getParameters();
                final boolean variadic = ps.size() > 0 && ps.get(ps.size()-1).isVarArgs();

                final ArrayList<Ref<Argument>> params = new ArrayList<Ref<Argument>>();
                for (CtParameter<?> p : ps) params.add(this.addArgument(p));
                
                final ArrayList<Ref<Argument>> results = new ArrayList<Ref<Argument>>();
                results.add(this.addArgument(m.getType()));

                return new Signature(variadic, params, results);
            });
    }

    public Ref<Signature> addSignature(CtMethod<?> m) throws Exception {
        assert(!isObjectMethod(m));
        return this.proj.signatures.create(this.log, m,
            "signature " + m.getSignature(),
            () -> {
                final List<CtParameter<?>> ps = m.getParameters();
                final boolean variadic = ps.size() > 0 && ps.get(ps.size()-1).isVarArgs();
                
                final ArrayList<Ref<Argument>> params = new ArrayList<Ref<Argument>>();
                for (CtParameter<?> p : ps) params.add(this.addArgument(p));
                
                final ArrayList<Ref<Argument>> results = new ArrayList<Ref<Argument>>();
                final CtTypeReference<?> res = m.getType();
                if (!isVoid(res)) results.add(this.addArgument(res));
                
                return new Signature(variadic, params, results);
            });
    }

    public Ref<Argument> addArgument(CtParameter<?> p) throws Exception {
        return this.proj.arguments.create(this.log, p,
            "parameter " + p.getSimpleName(),
            () -> {
                final String name = p.getSimpleName();
                final Ref<? extends TypeDesc> type = this.addTypeDesc(p.getType());
                return new Argument(name, type);
            });
    }
    
    public Ref<Argument> addArgument(CtTypeReference<?> p) throws Exception {
        return this.proj.arguments.create(this.log, p,
            "parameter <unnamed> " + p.getSimpleName(),
            () -> {
                final Ref<? extends TypeDesc> type = this.addTypeDesc(p);
                return new Argument("", type);
            });
    }
    
    public Ref<StructDesc> addStruct(CtClass<?> c) throws Exception {
        return this.proj.structDescs.create(this.log, c,
            "struct " + c.getQualifiedName(),
            () -> {
                // TODO: Handle enum?
                //if (c instanceof CtEnum<?> e) {}

                // Collect all fields.
                final ArrayList<Ref<Field>> fields = new ArrayList<Ref<Field>>();
                for (CtFieldReference<?> fr : c.getAllFields())
                    fields.add(this.addField(fr.getFieldDeclaration()));

                // Add extended class as a "$super" field.
                CtTypeReference<?> superFr = c.getSuperclass();
                if (superFr != null) fields.add(this.addField("$super", superFr));

                // Add access to nesting class as a "$nest" field.
                if (c.getRoleInParent() == CtRole.NESTED_TYPE) {
                    if (c.getParent() instanceof CtTypeReference<?> nest && nest != null) {
                        fields.add(this.addField("$nest", nest));
                    } else {
                        this.log.error("Unhandled nested object decl "+ c.getQualifiedName() + " in " + c.getParent());
                    }
                }

                return new StructDesc(fields);
            });
    }

    private Ref<Field> addField(CtField<?> f) throws Exception {
        return this.proj.fields.create(this.log, f,
            "field " + f.getSimpleName(),
            () -> {
                final String name = f.getSimpleName();
                final Ref<? extends TypeDesc> type = this.addTypeDesc(f.getType());
                return new Field(name, type);
            },
            (Ref<Field> ref, Field field) -> {
                field.setVisibility(f);
            });
    }

    private Ref<Field> addField(String name, CtTypeReference<?> f) throws Exception {
        return this.proj.fields.create(this.log, f,
            "field " + name,
            () -> {
                final Ref<? extends TypeDesc> type = this.addTypeDesc(f);
                return new Field(name, type);
            });
    }

    public Ref<Selection> addSelection(CtField<?> field) throws Exception {
        final String name = field.getSimpleName();
        return this.proj.selections.create(this.log, field,
            "select field " + name,
            () -> {
                final Ref<? extends Construct> decl = this.addDeclaration(field.getDeclaringType());
                return new Selection(name, decl);
            });
    }
    
    public Ref<InterfaceDecl> addInterfaceDecl(CtInterface<?> i) throws Exception {
        return this.proj.interfaceDecls.create(this.log, i,
            "interface decl " + i.getQualifiedName(),
            () -> {
                final Ref<PackageCon>    pkg   = this.addPackage(i.getPackage());
                final Location           loc   = proj.locations.create(i.getPosition());
                final String             name  = i.getSimpleName();
                final Ref<InterfaceDesc> inter = this.addInterfaceDesc(i);
                final ArrayList<Ref<TypeParam>> typeParams = this.addTypeParams(i.getFormalCtTypeParameters());

                if (i.getRoleInParent() == CtRole.NESTED_TYPE) {
                    // TODO: Need to differentiate this from an interface by
                    //       the same name nested in a different class or not nested in any class.
                    this.log.error("Unhandled nested interface decl "+ i.getQualifiedName());
                }

                return new InterfaceDecl(pkg, loc, name, inter, typeParams);
            },
            (Ref<InterfaceDecl> ref, InterfaceDecl id) -> {
                id.setVisibility(i);
                //if (id.pkg != null) id.pkg.interfaceDecls.add(id); // TODO: Move to a follow up when we know the package is done.
            });
    }

    public Ref<ObjectDecl> addEnum(CtEnum<?> e) throws Exception {
        return this.proj.objectDecls.create(this.log, e,
            "enum " + e.getQualifiedName(),
            () -> {
                final Ref<PackageCon> pkg  = this.addPackage(e.getPackage());
                final Location        loc  = proj.locations.create(e.getPosition());
                final String          name = e.getQualifiedName();

                final CtTypeReference<?> tr = e.getSuperclass();
                Ref<StructDesc> struct = this.proj.structDescs.create(this.log, tr,
                    "enum struct " + e.getQualifiedName(),
                    () -> {
                        final ArrayList<Ref<Field>> fields = new ArrayList<Ref<Field>>();
                        fields.add(this.addField("$value", tr));
                        return new StructDesc(fields);
                    });

                return new ObjectDecl(pkg, loc, name, struct, null);
            },
            (Ref<ObjectDecl> ref, ObjectDecl od) -> {
                // TODO: Finish by adding the "const values" to the package for each enumerator value.
            });
    }

    public Ref<? extends TypeDesc> addTypeDesc(CtTypeReference<?> tr) throws Exception {
        if (tr == null)       return this.proj.baker.objectDesc();
        if (tr.isPrimitive()) return this.addBasic(tr);
        if (tr.isArray())     return this.addArray((CtArrayTypeReference<?>)tr);

        CtType<?> ty = tr.getDeclaration();
        if (ty == null)       return this.proj.baker.objectDesc();
        if (tr.isClass())     return this.addObjectDecl((CtClass<?>)ty);
        if (tr.isInterface()) return this.addInterfaceDecl((CtInterface<?>)ty);
        if (tr.isGenerics())  return this.addTypeParam((CtTypeParameter)ty);
        if (tr.isEnum())      return this.addEnum((CtEnum<?>)ty);

        this.unknownTypeDesc(tr);
        return null;
    }

    private void unknownTypeDesc(CtTypeReference<?> tr) throws Exception {
        this.log.error("Unhandled (" + tr.getClass().getName() + "): "+tr.prettyprint());
        this.log.push();
        this.log.notice("isAnnotationType:    " + tr.isAnnotationType());
        this.log.notice("isAnonymous:. . . . ." + tr.isAnonymous());
        this.log.notice("isArray:             " + tr.isArray());
        this.log.notice("isClass:. . . . . . ." + tr.isClass());
        this.log.notice("isEnum:              " + tr.isEnum());
        this.log.notice("isGenerics: . . . . ." + tr.isGenerics());
        this.log.notice("isImplicit:          " + tr.isImplicit());
        this.log.notice("isInterface:. . . . ." + tr.isInterface());
        this.log.notice("isLocalType:         " + tr.isLocalType());
        this.log.notice("isParameterized:. . ." + tr.isParameterized());
        this.log.notice("isParentInitialized: " + tr.isParentInitialized());
        this.log.notice("isPrimitive:. . . . ." + tr.isPrimitive());
        this.log.notice("isShadow:            " + tr.isShadow());
        this.log.notice("isSimplyQualified:. ." + tr.isSimplyQualified());
        this.log.pop();
        throw new Exception("Unhandled Type");
    }

    public Ref<InterfaceInst> addArray(CtArrayTypeReference<?> tr) throws Exception {
        final Ref<? extends TypeDesc> elem = this.addTypeDesc(tr.getArrayType());
        Ref<InterfaceInst> ref = this.proj.baker.arrayInst(tr.getQualifiedName(), elem);
        this.proj.interfaceInsts.setRefForElem(tr, ref);
        //inst.generic.instances.add(inst); // TODO: Move to a follow up when we know the package is done.
        return ref;
    }
    
    public Ref<Basic> addBasic(CtTypeReference<?> tr) throws Exception {
        return this.proj.basics.create(this.log, tr,
            "basic " + tr.getSimpleName(),
            () -> {
                final String name = tr.getSimpleName();
                this.log.errorIf(name == "void", "A void was added as a basic.");
                return new Basic(name);
            });
    }

    public Ref<InterfaceDesc> addInterfaceDesc(CtInterface<?> i) throws Exception {
        return this.proj.interfaceDescs.create(this.log, i,
            "interface description " + i.getSimpleName(),
            () -> {
                final TreeSet<Ref<Abstract>> abstracts = new TreeSet<Ref<Abstract>>();
                for (CtMethod<?> m : i.getAllMethods()) {
                    if (!isObjectMethod(m)) abstracts.add(this.addAbstract(m));
                }

                // TODO: Determine how to pin this interface.
                return new InterfaceDesc(abstracts);
            },
            (Ref<InterfaceDesc> ref, InterfaceDesc id) -> {
                // TODO: Implement Inheritance
            });
    }

    public Ref<Abstract> addAbstract(CtMethod<?> m) throws Exception {
        assert(!isObjectMethod(m));
        return this.proj.abstracts.create(this.log, m,
            "abstract " + m.getSimpleName(),
            () -> {
                final String name = m.getSimpleName();
                final Ref<Signature> signature = this.addSignature(m);
                return new Abstract(name, signature);
            });
    }

    public ArrayList<Ref<? extends TypeDesc>> addTypeArguments(List<CtTypeReference<?>> trs) throws Exception {
        ArrayList<Ref<? extends TypeDesc>> result = new ArrayList<Ref<? extends TypeDesc>>(trs.size());
        for (CtTypeReference<?> tr : trs) result.add(this.addTypeDesc(tr));
        return result;
    }

    public ArrayList<Ref<TypeParam>> addTypeParams(List<CtTypeParameter> tps) throws Exception {
        ArrayList<Ref<TypeParam>> result = new ArrayList<Ref<TypeParam>>(tps.size());
        for (CtTypeParameter tp : tps) result.add(this.addTypeParam(tp));
        return result;
    }

    public Ref<TypeParam> addTypeParam(CtTypeParameter tp) throws Exception {
        return this.proj.typeParams.create(this.log, tp,
            "type params " + tp.getQualifiedName(),
            () -> {
                final String name = tp.getQualifiedName();
                
                // TODO: Remove
                //System.out.println(">> " + name + " >> " + tp.prettyprint());
                //System.out.println(">> >> " + tp.getTypeErasure());
                //for (CtTypeReference<?> tpr : tp.getSuperInterfaces())
                //    System.out.println(">>  >> " + tpr.getSimpleName() + " >> " + tpr.prettyprint());

                CtTypeReference<?> tr = tp.getTypeErasure();
                final Ref<? extends TypeDesc> type = this.addTypeDesc(tr);

                // TODO: Finish
                return new TypeParam(name, type);
            });
    }

    public void finish() throws Exception {
       this.processPendingMetrics();
       this.consolidateCons();
       this.crossConnectConstructs();
       this.validate();
    }

    private void processPendingMetrics() throws Exception {
        for(CtMethod<?> m : this.pendingMetrics) {
            if (m.getBody() == null) continue;
            if (m.getBody().getStatements().isEmpty()) continue;

            Ref<MethodDecl> ref = this.proj.methodDecls.getRef(m);
            if (!ref.isResolved())
                throw new Exception("Expected " + ref + " to be resolved before processing pending metrics.");

            MethodDecl md = ref.getResolved();
            if (md.metrics != null)
                throw new Exception("The metrics for " + md + " have already been processed before " + m.getSimpleName() + ".");

            Ref<Metrics> metRef = this.addMetrics(m);
            Metrics met = metRef.getResolved();
            if (met.hasBody()) md.metrics = metRef;
            else {
                // remove the reference and metrics from factory since bodiless methods can be ignored.
                this.proj.metrics.removeElem(this.log, m, "metrics " + m.getSimpleName());
            }
        }
    }

    public Ref<Metrics> addMetrics(CtMethod<?> m) throws Exception {
        return this.proj.metrics.create(this.log, m,
            "metrics " + m.getSimpleName(),
            () -> {
                final Location loc = proj.locations.create(m.getPosition());
                final Analyzer ana = new Analyzer(this, loc);
                ana.addMethod(m);
                return ana.getMetrics();
            });
    }

    private void consolidateCons() throws Exception {
        this.proj.setAllIndices();
        while (this.proj.consolidateCons(this.log))
            this.proj.setAllIndices();
        this.proj.setAllIndices();
    }

    private void crossConnectConstructs() throws Exception {
        for (MethodDecl m : this.proj.methodDecls.conSet)
            m.pkg.getResolved().methodDecls.add(this.proj.methodDecls.addOrGetRef(m));

        for (ObjectDecl obj : this.proj.objectDecls.conSet) {
            final PackageCon pkg = obj.pkg.getResolved();
            for (Ref<MethodDecl> met : obj.methodDecls)
                pkg.methodDecls.add(met);
        }

        // TODO: Add more to packages
    }

    private void validate() throws Exception {
        new Validator(this.log, this.proj).validate();
        if (this.log.errorCount() > 0)
            throw new Exception("Errors logged before or during validation.");
    }
}
