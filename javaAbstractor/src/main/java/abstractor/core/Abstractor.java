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
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.*;
import spoon.reflect.path.CtRole;
import spoon.reflect.reference.*;
import spoon.support.compiler.VirtualFile;
import abstractor.core.constructs.*;
import abstractor.core.log.*;
import abstractor.core.validator.Validator;

public class Abstractor {
    static private final String nullName = "<nulltype>";
    static private final String objName = "java.lang.Object";
    static private final boolean doNotCatch = true; // TODO: false;

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
     * Parses the source for one or more classes and adds it.
     * 
     * This is designed to test classes, records, and enumerators,
     * but will not work for interfaces.
     * @example parseClass("class C { void m() { System.out.println(\"hello\"); } }"); 
     * @param source The source code containing one or more classes.
     */
    public void addClassesFromSource(String ...sourceLines) throws Exception {
        final String filename = "ClassesFromSource.java";
        final String source = String.join("\n", sourceLines);
        Launcher launcher = new Launcher();
        launcher.addInputResource(new VirtualFile(source, filename));
        launcher.buildModel();
        this.addModel(launcher.getModel());
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

    // TODO: Implement package imports by deriving from actual type usage
    //       rather than import statements. This will be done in a later step
    //       when the Resolver pipeline is created.

    public Ref<PackageCon> addPackage(CtPackage pkg) throws Exception {
        final String name = packageName(pkg);
        return this.proj.packages.create(this.log, pkg,
           "package " + name,
            () -> {
                final String path = packagePath(pkg);
                return new PackageCon(name, path);
            },
            (Ref<PackageCon> ref, PackageCon pkgCon) ->{
                for (CtType<?> t : pkg.getTypes()) {
                    Ref<? extends Construct> decl = this.addDeclaration(t);
                    if (decl != null)
                        this.addDeclarationToPackage(pkgCon, decl);
                }
            });
    }
    
    // TODO: Use addPackageFor for more places.
    public Ref<PackageCon> addPackageFor(CtType<?> t) throws Exception {
        return this.addPackage(t.getTopLevelType().getPackage());
    }

    public Ref<PackageCon> addPackageFor(CtTypeReference<?> tr) throws Exception {
        return this.addPackageFor(tr.getTypeDeclaration());
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
        if (tryToAdd(pkg.objectDecls,    decl, ConstructKind.OBJECT_DECL))    return;
        if (tryToAdd(pkg.interfaceDecls, decl, ConstructKind.INTERFACE_DECL)) return;
        if (tryToAdd(pkg.methodDecls,    decl, ConstructKind.METHOD_DECL))    return;
        if (tryToAdd(pkg.values,         decl, ConstructKind.VALUE))          return;
        
        this.log.error("Unhandled declaration type: " + decl.kind());
    }

    public Ref<? extends Construct> addDeclaration(CtElement elem) throws Exception {
        if (elem == null) return null;
        if (doNotCatch) return this.addDeclarationImpl(elem);
        try {
            return this.addDeclarationImpl(elem);
        } catch (Exception ex) {
            this.log.warning("addDeclaration failed for " + this.describeElement(elem) + ": " + ex.getMessage());
            return null;
        }
    }

    private Ref<? extends Construct> addDeclarationImpl(CtElement elem) throws Exception {
        if (elem instanceof CtTypeReference<?> tr) elem = tr.getTypeDeclaration();
        if (elem == null) return null;

        // Skip annotation types — they don't participate in data flow.
        if (elem instanceof CtAnnotationType<?>) return null;

        // Check CtEnum before CtClass since CtEnum extends CtClass.
        if (elem instanceof CtEnum<?>      e) return this.addObjectDecl(e);
        if (elem instanceof CtClass<?>     c) return this.addObjectDecl(c);
        if (elem instanceof CtInterface<?> i) return this.addInterfaceDecl(i);
        if (elem instanceof CtMethod<?>    m) return this.addGeneralMethod(m);

        this.log.warning("Skipping unhandled decl (" + elem.getClass().getName() + ")");
        return null;
    }

    public Ref<? extends TypeDeclaration> addTypeDeclaration(CtElement elem) throws Exception {
        if (elem == null) return null;
        if (doNotCatch) return this.addTypeDeclarationImpl(elem);
        try {
            return this.addTypeDeclarationImpl(elem);
        } catch (Exception ex) {
            this.log.warning("addTypeDeclaration failed for " + this.describeElement(elem) + ": " + ex.getMessage());
            return null;
        }
    }

    private Ref<? extends TypeDeclaration> addTypeDeclarationImpl(CtElement elem) throws Exception {
        if (elem instanceof CtTypeReference<?> tr) elem = tr.getTypeDeclaration();
        if (elem == null) return null;

        if (elem instanceof CtEnum<?>      e) return this.addObjectDecl(e);
        if (elem instanceof CtClass<?>     c) return this.addObjectDecl(c);
        if (elem instanceof CtInterface<?> i) return this.addInterfaceDecl(i);

        this.log.warning("Skipping unhandled type decl (" + elem.getClass().getName() + ")");
        return null;
    }

    /** Short description for logs (does not throw). */
    private String describeElement(CtElement elem) {
        if (elem == null) return "(null)";
        try {
            if (elem instanceof CtTypeReference<?> tr) return tr.getQualifiedName();
            if (elem instanceof CtType<?>          ty) return ty.getQualifiedName();
            if (elem instanceof CtExecutable<?>    ex) return ex.getSignature();
        } catch (Exception ex) {
            this.log.error("describe element failed: " + ex.getMessage());
        }
        return elem.getClass().getName();
    }

    /**
     * Handle Java primitives and object equivalents to primitives (boxed primitives)
     * such that the boxed primitives (e.g. Integer) and String become Basic's.
     * Other types become stub InterfaceDecl's, with InterfaceInst when parameterized.
     */
    public Ref<? extends TypeDesc> addExternalStub(CtTypeReference<?> tr) throws Exception {

        // TODO: Why even pass in null?
        //
        // If a type can not be resolved return an object (kind of like an `any` in Go).
        if (tr == null) return this.proj.baker.anyDesc();
        try {
            // Unlike Go's nil that can carry the type, Java's null type has
            // no type associated with it so instead use an object.
            if (nullName.equals(tr.getQualifiedName())) return this.proj.baker.anyDesc();
        } catch (Exception ex) {
            this.log.error("addExternalStub failed resolving qualified name: " + ex.getMessage());
            return this.proj.baker.anyDesc();
        }

        final CtTypeReference<?> erasure = tr.getTypeErasure();
        final Ref<Basic> boxed = this.proj.baker.basicForBoxedOrString(erasure.getQualifiedName());
        if (boxed != null) return boxed;

        final Ref<InterfaceDecl> decl = this.addErasureInterfaceDecl(erasure);
        // Check if the type is a generic instantiation.
        final List<CtTypeReference<?>> typeArgs = tr.getActualTypeArguments();
        if (typeArgs == null || typeArgs.isEmpty()) return decl;

        return this.proj.interfaceInsts.create(this.log, tr,
            "type erasure interface instance " + tr.getQualifiedName(),
            () -> {
                final ArrayList<Ref<? extends TypeDesc>> instanceTypes = new ArrayList<>(typeArgs.size());
                for (CtTypeReference<?> arg : typeArgs) instanceTypes.add(this.addTypeDesc(arg));
    
                // TODO: decl.getResolved().inter can be null. Figue out another ady to do this!
                return new InterfaceInst(decl, instanceTypes, decl.getResolved().inter);
            });
    }

    private Ref<InterfaceDecl> addErasureInterfaceDecl(CtTypeReference<?> typeErasure) throws Exception {
        return this.proj.interfaceDecls.create(this.log, typeErasure,
            "type erasure interface decl " + typeErasure.getQualifiedName(),
            () -> {
                final Ref<PackageCon>    pkg    = this.addPackageFor(typeErasure);
                final Location           loc    = this.proj.locations.create(typeErasure.getPosition());
                final String             name   = typeErasure.getSimpleName();
                final Ref<InterfaceDesc> inter  = this.proj.baker.anyDesc();
                return new InterfaceDecl(pkg, loc, name, inter, new ArrayList<>());
            });
    }

    /**
     * This determines if the given method is a method on the base Object.
     * Since all Objects inherits the base Object, adding those methods are
     * just additional unneeded noise in the abstraction.
     */
    static public boolean isObjectMethod(CtMethod<?> m) {
        if (m == null) return false;

        final CtTypeReference<?> objectRef  = m.getFactory().Type().objectType();
        final CtType<?>          objectDecl = objectRef.getTypeDeclaration();
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
                final Location             loc        = this.proj.locations.create(c.getPosition());
                final String               name       = c.getSimpleName();
                final Ref<StructDesc>      struct     = this.addStruct(c);
                final List<Ref<TypeParam>> typeParams = this.addTypeParams(c.getFormalCtTypeParameters());
                return new ObjectDecl(pkg, loc, name, struct, typeParams);
            },
            (Ref<ObjectDecl> ref, ObjectDecl obj) -> {
                obj.setVisibility(c);
                
                // Add constructors as (static) methods.
                for (CtConstructor<?> ctor : c.getConstructors()) {
                    if (ctor.getParent().equals(c)) {
                        // Skip default constructors
                        if (this.isDefaultConstructor(ctor)) continue;
                        this.addConstructorMethod(ref, ctor);
                    }
                }

                // Add methods for the class.
                for (CtMethod<?> m : c.getAllMethods()) {
                    if (m.getParent().equals(c) && !isObjectMethod(m)) this.addMethod(ref, m);
                }

                // Synthesize the interface description for the class.
                TreeSet<Ref<Abstract>> abstracts = new TreeSet<Ref<Abstract>>();
                for (CtMethod<?> m : c.getAllMethods()) {
                    if (!m.isStatic() && !isObjectMethod(m)) abstracts.add(this.addAbstract(m));
                }
                // TODO: FIX BELOW
                obj.inter = this.proj.interfaceDescs.addOrGetRef(new InterfaceDesc(abstracts, ref), "interface for object");

                // TODO: Finish implementing
                //System.out.println("1) >>> " + c.getSuperInterfaces());

                // Add any nested types.
                for (CtType<?> nt : c.getNestedTypes()) this.addTypeDesc(nt.getReference());
            });
    }

    public boolean isDefaultConstructor(CtConstructor<?> ctor) throws Exception {
        final boolean skip = ctor.isImplicit();
        if (skip) this.log.notice("skipping default constructor: " + ctor.getSignature());
        return skip;
    }

    public Ref<MethodDecl> addConstructorMethod(Ref<ObjectDecl> receiver, CtConstructor<?> ctor) throws Exception {
        if (!receiver.isResolved())
            throw new AbstractorException("Expected the receiver for a constructor method to be resolved: " + receiver.toString());
        ObjectDecl recv = receiver.getResolved();
        return this.proj.methodDecls.create(log, ctor,
            "constructor " + ctor.getSignature(),
            () -> {
                final Ref<PackageCon>      pkg        = recv.pkg;
                final Location             loc        = this.proj.locations.create(ctor.getPosition());
                final String               name       = recv.name;
                final Ref<Signature>       signature  = this.addConstructSignature(ctor);
                final List<Ref<TypeParam>> typeParams = this.addTypeParams(ctor.getFormalCtTypeParameters());
                MethodDecl md = new MethodDecl(pkg, receiver, loc, name, signature, typeParams);
                md.constructor = true;
                md.isStatic = true;
                return md;
            },
            (Ref<MethodDecl> ref, MethodDecl md) -> {
                md.setVisibility(ctor);
                recv.methodDecls.add(ref);
            });
    }

    public Ref<? extends Construct> addGeneralMethod(CtMethod<?> m) throws Exception {
        CtType<?> decl = m.getDeclaringType();
        if (decl instanceof CtEnum<?> e) {
            Ref<ObjectDecl> obj = this.addObjectDecl(e);
            return this.addMethod(obj, m);
        }
        if (decl instanceof CtClass<?> c) {
            Ref<ObjectDecl> obj = this.addObjectDecl(c);
            return this.addMethod(obj, m);
        }
        if (decl instanceof CtInterface<?> i) {
            Ref<Abstract> ab = this.addAbstract(m);
            // TODO: Connect abstract to interface declaration
            return ab;
        }
        this.log.warning("Skipping method with unhandled declaring type (" + decl.getClass().getName() + ") " + decl.getQualifiedName());
        return null;
    }

    public Ref<MethodDecl> addMethod(Ref<ObjectDecl> receiver, CtMethod<?> m) throws Exception {
        if (!receiver.isResolved())
            throw new AbstractorException("Expected the object receiver for a method to be resolved: " + receiver.toString());
        ObjectDecl recv = receiver.getResolved();

        return this.proj.methodDecls.create(this.log, m,
            "method " + m.getSignature(),
            () -> {
                final Ref<PackageCon>      pkg        = recv.pkg;
                final Location             loc        = this.proj.locations.create(m.getPosition());
                final String               name       = m.getSimpleName();
                final Ref<Signature>       signature  = this.addSignature(m);
                final List<Ref<TypeParam>> typeParams = this.addTypeParams(m.getFormalCtTypeParameters());
                MethodDecl md = new MethodDecl(pkg, receiver, loc, name, signature, typeParams);
                md.isStatic = m.isStatic();
                return md;
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
                final Location           loc   = this.proj.locations.create(i.getPosition());
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
                final Location        loc  = this.proj.locations.create(e.getPosition());
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
        if (doNotCatch) return this.addTypeDescImpl(tr);
        try {
            return this.addTypeDescImpl(tr);
        } catch (Exception ex) {
            this.log.warning("addTypeDesc failed for " + this.safeTypeRefName(tr) + ": " + ex.getMessage());
            return this.proj.baker.anyDesc();
        }
    }

    private String safeTypeRefName(CtTypeReference<?> tr) {
        if (tr == null) return "(null)";
        try {
            return tr.getQualifiedName();
        } catch (Exception ignored) {
            // Do not worry about the exception here because it is only from Spoon
            // getting a qualified name and we can recover by using the toString via valueOf.
            return String.valueOf(tr);
        }
    }

    private Ref<? extends TypeDesc> addTypeDescImpl(CtTypeReference<?> tr) throws Exception {
        if (tr == null)       return this.proj.baker.anyDesc();
        if (tr.isPrimitive()) return this.addBasic(tr);
        if (tr.isArray())     return this.addArray((CtArrayTypeReference<?>)tr);

        // Handle wildcard types (e.g., ?, ? extends Foo, ? super Bar).
        if (tr instanceof CtWildcardReference wr)
            return this.addWildcard(wr);

        // TODO: NEED TO REEVALUATE ALL OF THIS AI SLOP

        // Type of the `null` literal in Spoon - not a real external type.
        try {
            if (nullName.equals(tr.getQualifiedName())) return this.proj.baker.anyDesc();
        } catch (Exception ex) {
            this.log.error("addTypeDescImpl failed resolving qualified name: " + ex.getMessage());
            return this.proj.baker.anyDesc();
        }

        // Use getTypeDeclaration (not getDeclaration) to get shadow types
        // for external/JDK types instead of null.
        CtType<?> ty = null;
        try {
            ty = tr.getTypeDeclaration();
        } catch (Exception ex) {
            this.log.warning("Failed to get type declaration for " + tr.getQualifiedName() + ": " + ex.getMessage());
            return this.addExternalStub(tr);
        }

        // If still null, treat as external / unresolvable type.
        if (ty == null) return this.addExternalStub(tr);

        // Annotation types don't participate in data flow. Use an object instead.
        if (ty instanceof CtAnnotationType<?> ann) {
            this.log.notice("Annotation type as type desc (using object): " + ann.getQualifiedName());
            return this.proj.baker.anyDesc();
        }

        // Skip anonymous and local types - their code is attributed
        // to the enclosing method (handled in later steps).
        if (tr.isAnonymous() || tr.isLocalType()) {
            CtTypeReference<?> superRef = tr.getSuperclass();
            if (superRef != null && !superRef.getQualifiedName().equals(objName))
                return this.addTypeDesc(superRef);

            var superIfaces = tr.getSuperInterfaces();
            if (superIfaces != null && !superIfaces.isEmpty())
                return this.addTypeDesc(superIfaces.iterator().next());

            return this.proj.baker.anyDesc();
        }

        // Shadow types are external (JDK / third-party)
        if (ty.isShadow()) return this.addExternalStub(tr);

        // Check type parameter first since it's the most specific.
        if (tr.isGenerics()) return this.addTypeParam((CtTypeParameter)ty);

        // Check CtEnum before CtClass since CtEnum extends CtClass.
        if (ty instanceof CtEnum<?> e)    return this.addObjectDecl(e);
        if (ty instanceof CtClass<?>)     return this.addObjectDecl((CtClass<?>)ty);
        if (ty instanceof CtInterface<?>) return this.addInterfaceDecl((CtInterface<?>)ty);

        this.log.warning("Unhandled type (" + tr.getClass().getName() + "): " + tr.getQualifiedName());
        return this.proj.baker.anyDesc();
    }

    private Ref<? extends TypeDesc> addWildcard(CtWildcardReference wr) throws Exception {
        CtTypeReference<?> bound = wr.getBoundingType();
        if (bound == null || bound instanceof CtWildcardReference) return this.proj.baker.anyDesc();

        // Spoon often uses java.lang.Object as the synthetic bound for unbounded "?".
        // Resolving it would pull the entire JDK Object graph into the abstraction.
        try {
            if (objName.equals(bound.getQualifiedName())) return this.proj.baker.anyDesc();
        } catch (Exception ignored) {
            // The exception here can be ignored since it comes from Spoon
            // failing to get the qualified name and it can be recovered from.
            return this.proj.baker.anyDesc();
        }
        return this.addTypeDesc(bound);
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
                if (name == "void") {
                    this.log.error("A void was added as a basic.");
                    throw new AbstractorException("A void was added as a basic");
                }
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

                CtTypeReference<?> tr = tp.getTypeErasure();
                final Ref<? extends TypeDesc> type = this.addTypeDesc(tr);

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
        // `addMetrics` may register more methods on `pendingMetrics`
        // so add the current methods, then check if more are pending.
        while (!this.pendingMetrics.isEmpty()) {
            final ArrayList<CtMethod<?>> methods = new ArrayList<>(this.pendingMetrics);
            this.pendingMetrics.clear();
            for (CtMethod<?> m : methods) {
                if (m.getBody() == null) continue;
                if (m.getBody().getStatements().isEmpty()) continue;

                Ref<MethodDecl> ref = this.proj.methodDecls.getRef(m);
                if (!ref.isResolved())
                    throw new AbstractorException("Expected " + ref + " to be resolved before processing pending metrics.");

                MethodDecl md = ref.getResolved();
                if (md.metrics != null)
                    throw new AbstractorException("The metrics for " + md + " have already been processed before " + m.getSimpleName() + ".");

                Ref<Metrics> metRef = this.addMetrics(m);
                Metrics met = metRef.getResolved();
                if (met.hasBody()) md.metrics = metRef;
                else {
                    // remove the reference and metrics from factory since bodiless methods can be ignored.
                    this.proj.metrics.removeElem(this.log, m, "metrics " + m.getSimpleName());
                }
            }
        }
    }

    public Ref<Metrics> addMetrics(CtMethod<?> m) throws Exception {
        return this.proj.metrics.create(this.log, m,
            "metrics " + m.getSimpleName(),
            () -> {
                final Location loc = this.proj.locations.create(m.getPosition());
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
        for (MethodDecl m : this.proj.methodDecls.conSet) {
            Ref<MethodDecl> decl = this.proj.methodDecls.addOrGetRef(m, "method in package " + m.pkg);
            m.pkg.getResolved().methodDecls.add(decl);
        }

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
            throw new AbstractorException("Errors logged before or during validation.");
    }
}
