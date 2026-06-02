package abstractor.core;

import java.util.*;

import spoon.Launcher;
import spoon.MavenLauncher;
import spoon.reflect.*;
import spoon.reflect.code.CtConstructorCall;
import spoon.reflect.declaration.*;
import spoon.reflect.path.CtRole;
import spoon.reflect.reference.*;
import spoon.reflect.visitor.filter.TypeFilter;
import spoon.support.compiler.VirtualFile;

import abstractor.core.constructs.*;
import abstractor.core.log.*;
import abstractor.core.require.Require;
import abstractor.core.validator.*;
import abstractor.core.spoonUtils.*;

public class Abstractor {
    public final Logger  log;
    public final Project proj;

    public final HashSet<CtExecutable<?>> pendingMetrics  = new HashSet<>();
    public final HashSet<CtPackage>       pendingPackages = new HashSet<>();

    public CtModel model;

    public Abstractor(Logger log, Project proj) {
        this.log   = log;
        this.proj  = proj;
    }

    /**
     * Reads a project containing a pom.xml maven file.
     * @param mavenProject The path to the project file. 
     */
    public void prepareMavenProject(String mavenProject) throws Exception {
        this.log.log("Reading " + mavenProject);
        MavenLauncher launcher = new MavenLauncher(mavenProject, MavenLauncher.SOURCE_TYPE.APP_SOURCE);
        CtModel model = launcher.buildModel();
        if (model.getAllTypes().size() > 0) {
            this.prepareModel(model);
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
        this.prepareModel(model);
    }

    /**
     * Parses the source for one or more classes and adds it.
     * 
     * This is designed to test classes, records, and enumerators,
     * but will not work for interfaces.
     * @example parseClass("class C { void m() { System.out.println(\"hello\"); } }"); 
     * @param source The source code containing one or more classes.
     */
    public void prepareClassesFromSource(String ...sourceLines) throws Exception {
        final String   filename = "ClassesFromSource.java";
        final String   source   = String.join("\n", sourceLines);
        final Launcher launcher = new Launcher();
        launcher.addInputResource(new VirtualFile(source, filename));
        launcher.buildModel();
        this.prepareModel(launcher.getModel());
    }

    private void prepareModel(CtModel model) throws Exception {
        Require.isNull(this.model, "currently this can only handle one model at a time");
        this.model = model;
        for (CtPackage pkg: model.getAllPackages()) {
            this.log.log("Init pending package " + SpoonUtils.describeElem(pkg));
            this.pendingPackages.add(pkg);
        }
    }

    //===[ Construct Adders ]===================================================

    public Ref<PackageCon> addPackage(CtPackage pkg) throws Exception {
        if (pkg == null) return this.proj.baker.builtinPackage();

        Ref<PackageCon> pkgRef = this.proj.packages.getRefByElem(pkg);
        if (pkgRef != null) return pkgRef;

        this.log.log("Pending package " + SpoonUtils.describeElem(pkg));
        this.pendingPackages.add(pkg);
        return this.proj.packages.addOrGetRefForElem(pkg,
            "for pending package " + SpoonUtils.describeElem(pkg));
    }
    
    public Ref<PackageCon> addPackageFor(CtType<?> t) throws Exception {
        return this.addPackage(t.getTopLevelType().getPackage());
    }

    public Ref<PackageCon> addPackageFor(CtTypeReference<?> tr) throws Exception {
        return this.addPackageFor(tr.getTypeDeclaration());
    }

    public Ref<? extends Construct> addDeclaration(CtElement elem) throws Exception {
        if (elem == null) return null;

        // If a reference, get the actual element.
        if (elem instanceof CtReference ref) elem = ref.getDeclaration();
        if (elem == null) return null;

        // Skip annotation types — they don't participate in data flow.
        if (elem instanceof CtAnnotationType<?>) return null;

        // Check CtEnum before CtClass since CtEnum extends CtClass.
        if (elem instanceof CtEnum<?>      e) return this.addEnum(e);
        if (elem instanceof CtClass<?>     c) return this.addObjectDecl(c);
        if (elem instanceof CtInterface<?> i) return this.addInterfaceDecl(i);
        if (elem instanceof CtMethod<?>    m) return this.addMethodOrAbstract(m);

        this.log.error("Unhandled decl: " + SpoonUtils.describeElem(elem));
        return null;
    }

    public Ref<InterfaceDecl> addInterfaceDecl(CtInterface<?> i) throws Exception {
        return this.proj.interfaceDecls.create(this.log, i,
            "interface decl " + SpoonUtils.describeElem(i),
            () -> {
                final String             name  = i.getSimpleName();
                final Ref<PackageCon>    pkg   = this.addPackageFor(i);
                final Location           loc   = this.proj.locations.create(i.getPosition());
                final Ref<InterfaceDesc> inter = this.addInterfaceDesc(i);
                final ArrayList<Ref<TypeParam>> typeParams = this.addTypeParams(i.getFormalCtTypeParameters());
                return new InterfaceDecl(pkg, loc, name, inter, typeParams);
            },
            (Ref<InterfaceDecl> ref, InterfaceDecl id) -> {
                id.setVisibility(i);
            });
    }

    public Ref<InterfaceDesc> addInterfaceDesc(CtInterface<?> i) throws Exception {
        return this.proj.interfaceDescs.create(this.log, i,
            "interface description " + SpoonUtils.describeElem(i),
            () -> {
                final TreeSet<Ref<Abstract>> abstracts = new TreeSet<Ref<Abstract>>();
                for (CtMethod<?> m : i.getAllMethods()) {
                    if (!SpoonUtils.isObjectMethod(m)) abstracts.add(this.addAbstract(m));
                }

                Ref<? extends Construct> pin = null;
                if (i.getRoleInParent() == CtRole.NESTED_TYPE) {
                    CtElement parent = i.getParent();
                    if (parent instanceof CtTypeReference<?> nest && nest != null) {
                        pin = this.addDeclaration(nest);
                    } else {
                        this.log.error("Unhandled nested interface decl " + SpoonUtils.describeElem(i) + " in " + parent);
                    }
                }

                return new InterfaceDesc(abstracts, pin);
            },
            (Ref<InterfaceDesc> ref, InterfaceDesc id) -> {
                // TODO: Handle interface object declaration?
                //if (tr.isGenerics()) ...;

                // Add direct super-interfaces this interface extends
                for (CtTypeReference<?> supRef : i.getSuperInterfaces()) {
                    CtType<?> supDecl = supRef.getTypeDeclaration(); // may be null for shadow/unresolved
                    if (supDecl != null && supDecl instanceof CtInterface<?> supId && supId != null) {
                        id.inherits.add(this.addInterfaceDesc(supId));
                    } else {
                        this.log.error("Unhandled super-interface " + SpoonUtils.describeElem(supDecl) + " for " + id);
                    }
                }
            });
    }

    public Ref<? extends Construct> addMethodOrAbstract(CtMethod<?> m) throws Exception {
        final CtType<?> decl = m.getDeclaringType();
        if (decl instanceof CtEnum<?>    e) return this.addMethod(this.addEnum(e), m);
        if (decl instanceof CtClass<?>   c) return this.addMethod(this.addObjectDecl(c), m);
        if (decl instanceof CtInterface<?>) return this.addAbstract(m);

        this.log.error("method has unhandled declaring type: " + SpoonUtils.describeElem(decl));
        return null;
    }

    public Ref<MethodDecl> addMethod(Ref<ObjectDecl> receiver, CtMethod<?> m) throws Exception {
        Require.notObjectMethod(m);
        final ObjectDecl recv = receiver.mustGetResolved();
        return this.proj.methodDecls.create(this.log, m,
            "method " + SpoonUtils.describeElem(m),
            () -> {
                final Ref<PackageCon>      pkg        = recv.pkg;
                final Location             loc        = this.proj.locations.create(m.getPosition());
                final String               name       = m.getSimpleName();
                final Ref<Signature>       signature  = this.addSignature(m);
                final List<Ref<TypeParam>> typeParams = this.addTypeParams(m.getFormalCtTypeParameters());
                final MethodDecl md = new MethodDecl(pkg, receiver, loc, name, signature, typeParams);
                md.isStatic = m.isStatic();
                return md;
            },
            (Ref<MethodDecl> ref, MethodDecl md) -> {
                md.setVisibility(m);
                recv.methodDecls.add(ref);
                this.pendingMetrics.add(m);
            });
    }

    public Ref<Abstract> addAbstract(CtMethod<?> m) throws Exception {
        Require.notObjectMethod(m);
        return this.proj.abstracts.create(this.log, m,
            "abstract " + SpoonUtils.describeElem(m),
            () -> {
                final String         name      = m.getSimpleName();
                final Ref<Signature> signature = this.addSignature(m);
                return new Abstract(name, signature);
            });
    }

    public Ref<Signature> addSignature(CtMethod<?> m) throws Exception {
        Require.notObjectMethod(m);
        return this.proj.signatures.create(this.log, m,
            "signature " + SpoonUtils.describeElem(m),
            () -> {
                final List<CtParameter<?>> ps = m.getParameters();
                final boolean variadic = ps.size() > 0 && ps.get(ps.size()-1).isVarArgs();
                
                final ArrayList<Ref<Argument>> params = new ArrayList<>();
                for (CtParameter<?> p : ps) params.add(this.addArgument(p));
                
                final ArrayList<Ref<Argument>> results = new ArrayList<>();
                final CtTypeReference<?> res = m.getType();
                if (!SpoonUtils.isVoid(res)) results.add(this.addArgument(res));
                
                return new Signature(variadic, params, results);
            });
    }

    public Ref<MethodDecl> addConstructorMethod(Ref<ObjectDecl> receiver, CtConstructor<?> ctor) throws Exception {
        return this.proj.methodDecls.create(log, ctor,
            "constructor " + ctor.getSignature(),
            () -> {
                final ObjectDecl           recv       = receiver.mustGetResolved();
                final Ref<PackageCon>      pkg        = recv.pkg;
                final Location             loc        = this.proj.locations.create(ctor.getPosition());
                final String               name       = recv.name;
                final Ref<Signature>       signature  = this.addConstructorSignature(ctor);
                final List<Ref<TypeParam>> typeParams = this.addTypeParams(ctor.getFormalCtTypeParameters());
                final MethodDecl md = new MethodDecl(pkg, receiver, loc, name, signature, typeParams);
                md.constructor = true;
                md.isStatic = true;
                return md;
            },
            (Ref<MethodDecl> ref, MethodDecl md) -> {
                md.setVisibility(ctor);
                final ObjectDecl recv = receiver.mustGetResolved();
                recv.methodDecls.add(ref);
                this.pendingMetrics.add(ctor);
            });
    }

    public Ref<Signature> addConstructorSignature(CtConstructor<?> m) throws Exception {
        return this.proj.signatures.create(this.log, m,
            "constructor signature " + SpoonUtils.describeElem(m),
            () -> {
                final List<CtParameter<?>> ps = m.getParameters();
                final boolean variadic = ps.size() > 0 && ps.get(ps.size()-1).isVarArgs();

                final ArrayList<Ref<Argument>> params = new ArrayList<>();
                for (CtParameter<?> p : ps) params.add(this.addArgument(p));
                
                final ArrayList<Ref<Argument>> results = new ArrayList<>();
                results.add(this.addArgument(m.getType()));

                return new Signature(variadic, params, results);
            });
    }

    public Ref<Argument> addArgument(CtParameter<?> p) throws Exception {
        return this.proj.arguments.create(this.log, p,
            "parameter " + SpoonUtils.describeElem(p),
            () -> {
                final String                  name = p.getSimpleName();
                final Ref<? extends TypeDesc> type = this.addTypeDesc(p.getType());
                return new Argument(name, type);
            });
    }
    
    public Ref<Argument> addArgument(CtTypeReference<?> p) throws Exception {
        return this.proj.arguments.create(this.log, p,
            "parameter <unnamed> " + SpoonUtils.describeElem(p),
            () -> {
                final Ref<? extends TypeDesc> type = this.addTypeDesc(p);
                return new Argument("", type);
            });
    }
    
    public Ref<StructDesc> addStruct(CtClass<?> c) throws Exception {
        return this.proj.structDescs.create(this.log, c,
            "struct " + SpoonUtils.describeElem(c),
            () -> {
                // Collect all fields.
                final ArrayList<Ref<Field>> fields = new ArrayList<>();
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
                        this.log.error("Unhandled nested object decl " + SpoonUtils.describeElem(c) + " in " + c.getParent());
                    }
                }

                return new StructDesc(fields);
            });
    }

    public Ref<Field> addField(CtField<?> f) throws Exception {
        return this.proj.fields.create(this.log, f,
            "field " + SpoonUtils.describeElem(f),
            () -> {
                final String                  name = f.getSimpleName();
                final Ref<? extends TypeDesc> type = this.addTypeDesc(f.getType());
                return new Field(name, type);
            },
            (Ref<Field> ref, Field field) -> {
                field.setVisibility(f);
            });
    }

    public Ref<Field> addField(String name, CtTypeReference<?> f) throws Exception {
        return this.proj.fields.create(this.log, f,
            "field " + name,
            () -> {
                final Ref<? extends TypeDesc> type = this.addTypeDesc(f);
                return new Field(name, type);
            });
    }

    public Ref<Selection> addSelection(CtField<?> field) throws Exception {
        return this.proj.selections.create(this.log, field,
            "select field " + SpoonUtils.describeElem(field),
            () -> {
                final String                   name = field.getSimpleName();
                final Ref<? extends Construct> decl = this.addDeclaration(field.getDeclaringType());
                return new Selection(name, decl);
            });
    }

    public Ref<? extends TypeDesc> addArray(CtArrayTypeReference<?> tr) throws Exception {
        final Ref<? extends TypeDesc> td = this.addTypeDesc(tr.getArrayType());

        // Check that `td` is not `T` to prevent $Array<T> being instantiated with T.
        if (td.isResolved()) {
            final Ref<TypeParam> tdT = this.proj.baker.genT();
            if (td.getResolved().equals(tdT.getResolved()))
                return this.proj.baker.arrayDecl();
        }

        final Ref<InterfaceInst> ref = this.proj.baker.arrayInst(tr.getSimpleName(), td);
        return this.proj.interfaceInsts.setRefForElem(tr, ref);
    }
    
    public Ref<Basic> addBasic(CtTypeReference<?> tr) throws Exception {
        return this.proj.basics.create(this.log, tr,
            "basic " + SpoonUtils.describeElem(tr),
            () -> {
                if (SpoonUtils.isVoid(tr))
                    throw new AbstractorException("A void was added as a basic");
                return new Basic(tr.getSimpleName());
            });
    }

    public ArrayList<Ref<? extends TypeDesc>> addTypeArguments(List<CtTypeReference<?>> trs) throws Exception {
        final ArrayList<Ref<? extends TypeDesc>> result = new ArrayList<>(trs.size());
        for (CtTypeReference<?> tr : trs) result.add(this.addTypeDesc(tr));
        return result;
    }

    public ArrayList<Ref<TypeParam>> addTypeParams(List<CtTypeParameter> tps) throws Exception {
        final ArrayList<Ref<TypeParam>> result = new ArrayList<>(tps.size());
        for (CtTypeParameter tp : tps) result.add(this.addTypeParam(tp));
        return result;
    }

    public Ref<TypeParam> addTypeParam(CtTypeParameter tp) throws Exception {
        return this.proj.typeParams.create(this.log, tp,
            "type params " + SpoonUtils.describeElem(tp),
            () -> {
                final String                  name = tp.getSimpleName();
                final CtTypeReference<?>      tr   = tp.getTypeErasure();
                final Ref<? extends TypeDesc> type = this.addTypeDesc(tr);
                return new TypeParam(name, type);
            });
    }
    
    public Ref<Metrics> addMetrics(CtExecutable<?> m) throws Exception {
        return this.proj.metrics.create(this.log, m,
            "metrics " + SpoonUtils.describeElem(m),
            () -> {
                final Location loc = this.proj.locations.create(m.getPosition());
                final Analyzer ana = new Analyzer(this, loc);
                ana.addMethod(m);
                return ana.getMetrics();
            });
    }

    public Ref<ObjectDecl> addObjectDecl(CtClass<?> c) throws Exception {
        Require.notObject(c.getReference());
        return this.proj.objectDecls.create(this.log, c,
            "object decl " + SpoonUtils.describeElem(c),
            () -> {
                final Ref<PackageCon>      pkg        = this.addPackageFor(c);
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
                        if (ctor.isImplicit()) {
                            this.log.notice("skipping default constructor: " + ctor.getSignature());
                            continue;
                        }
                        this.addConstructorMethod(ref, ctor);
                    }
                }

                // Add methods for the class.
                for (CtMethod<?> m : c.getAllMethods()) {
                    if (m.getParent().equals(c) && !SpoonUtils.isObjectMethod(m))
                        this.addMethod(ref, m);
                }

                // Synthesize the interface abstractions for the class.
                final TreeSet<Ref<Abstract>> abstracts = new TreeSet<Ref<Abstract>>();
                for (CtMethod<?> m : c.getAllMethods()) {
                    if (!m.isStatic() && !SpoonUtils.isObjectMethod(m))
                        abstracts.add(this.addAbstract(m));
                }

                // Synthesize the interface description for the class.
                if (abstracts.size() > 0 || c.getSuperInterfaces().size() > 0) {
                    final InterfaceDesc it = new InterfaceDesc(abstracts, ref);
                    obj.inter = this.proj.interfaceDescs.addOrGetRef(it, "interface for object");

                    // Add direct super-interfaces this object extends.
                    for (CtTypeReference<?> supRef : c.getSuperInterfaces()) {
                        CtType<?> supDecl = supRef.getTypeDeclaration(); // may be null for shadow/unresolved
                        if (supDecl != null && supDecl instanceof CtInterface<?> supId && supId != null) {
                            it.inherits.add(this.addInterfaceDesc(supId));
                        } else {
                            this.log.error("Unhandled super-interface " + SpoonUtils.describeElem(supDecl) + " for " + obj);
                        }
                    }
                } else {
                    obj.inter = this.proj.baker.anyDesc();
                }

                // Add any nested types.
                for (CtType<?> nt : c.getNestedTypes())
                    this.addTypeDesc(nt.getReference());

                // Add any generic instances.
                if (c.isGenerics())
                    this.addObjectInstances(c, ref, obj);
            });
    }

    public Ref<? extends TypeDesc> addTypeDesc(CtTypeReference<?> tr) throws Exception {
        if (tr == null) return null;

        // Skip anonymous and local types since they can not escape the enclosing method.
        // (They still will contribute to metrics via super-interfaces and extends).
        if (tr.isAnonymous()) {
            this.log.notice("Ignoring anonymous type: " + SpoonUtils.describeElem(tr));
            return null;
        }
        if (tr.isLocalType()) {
            this.log.notice("Ignoring local type: " + SpoonUtils.describeElem(tr));
            return null;
        }

        // Handle primitive types (i.e. `int` but not `String` nor `Integer`).
        if (tr.isPrimitive()) return this.addBasic(tr);
        
        // Handle an array (i.e. `T[]` not `List<T>`) type.
        if (tr.isArray()) return this.addArray((CtArrayTypeReference<?>)tr);

        // Handle wildcard types (e.g., `?`, `? extends Foo`, `? super Bar`).
        if (tr instanceof CtWildcardReference wr) return this.addWildcard(wr);

        // Type of the `null` literal in Spoon and not a real external type.
        if (SpoonUtils.isNull(tr)) return this.proj.baker.anyDesc();

        // A boxed type (e.g. Integer, String) that we alias as a basic.
        final Ref<Basic> boxed = this.proj.baker.basicForBoxedOrString(tr);
        if (boxed != null) return boxed;

        // Shadow types are external (JDK / third-party) without a type declaration.
        if (tr.isShadow()) return this.addShadowTypeDesc(tr);
        
        // If the type is an Object, return an any for the Object.
        if (SpoonUtils.isObject(tr)) return this.proj.baker.anyDesc();

        // Use getTypeDeclaration (not getDeclaration) to get shadow types
        // for external/JDK types instead of null.
        CtType<?> ty = tr.getTypeDeclaration();
        if (ty == null) {
            this.log.error("Type description did not have a declaration but "+
                "was not labelled a anonymous: " + SpoonUtils.describeElem(tr));
            return null;
        }

        // Annotation types don't participate in data flow. Use an object instead.
        if (ty instanceof CtAnnotationType<?> ann) {
            this.log.notice("Ignoring annotation type: " + SpoonUtils.describeElem(ann));
            return null;
        }

        // Check CtEnum before CtClass since CtEnum extends CtClass.
        if (ty instanceof CtEnum<?>        e) return this.addEnum(e);
        if (ty instanceof CtClass<?>       c) return this.addObjectDecl(c);
        if (ty instanceof CtInterface<?>   i) return this.addInterfaceDecl(i);
        if (ty instanceof CtTypeParameter tp) return this.addTypeParam(tp);

        this.log.warning("Unhandled type description: " + SpoonUtils.describeElem(ty));
        return null;
    }

    public Ref<? extends TypeDesc> addWildcard(CtWildcardReference wr) throws Exception {
        CtTypeReference<?> bound = wr.getBoundingType();
        // Spoon often uses java.lang.Object as the synthetic bound for unbounded "?".
        // Resolving it would pull the entire JDK Object graph into the abstraction.
        if (bound == null || bound instanceof CtWildcardReference || SpoonUtils.isObject(bound))
            return this.proj.baker.anyDesc();
        return this.addTypeDesc(bound);
    }

    public Ref<ObjectDecl> addEnum(CtEnum<?> e) throws Exception {
        return this.proj.objectDecls.create(this.log, e,
            "enum " + SpoonUtils.describeElem(e),
            () -> {
                final String          name = e.getSimpleName();
                final Ref<PackageCon> pkg  = this.addPackageFor(e);
                final Location        loc  = this.proj.locations.create(e.getPosition());

                final CtTypeReference<?> tr = e.getSuperclass();
                final Ref<StructDesc> struct = this.proj.structDescs.create(this.log, tr,
                    "enum struct " + SpoonUtils.describeElem(tr),
                    () -> {
                        final ArrayList<Ref<Field>> fields = new ArrayList<>();
                        fields.add(this.addField("$value", tr));
                        return new StructDesc(fields);
                    });

                return new ObjectDecl(pkg, loc, name, struct, null);
            },
            (Ref<ObjectDecl> ref, ObjectDecl od) -> {
                // Finish by adding the "const values" to the package for each enumerator value.
                for (CtEnumValue<?> ev: e.getEnumValues()) {
                    this.proj.values.create(this.log, e,
                        "enum value "+ SpoonUtils.describeElem(ev),
                        () -> {
                            final String   name = ev.getSimpleName();
                            final Location loc  = this.proj.locations.create(ev.getPosition());
                            return new Value(od.pkg, loc, name, true, null, ref);
                        });
                }
            });
    }

    public Ref<? extends TypeDesc> addShadowTypeDesc(CtTypeReference<?> tr) throws Exception {
        // from isShadow() method:
        // > When an element isn't present in the factory (created in another factory),
        // > this element is considered as "shadow". e.g., a shadow element can be a
        // > CtType of java.lang.Class built when we call CtTypeReference.getTypeDeclaration()
        // > on a reference of java.lang.Class."

        return this.proj.baker.anyDesc();
    }

    public void addObjectInstances(CtClass<?> c, Ref<ObjectDecl> ref, ObjectDecl obj) {
        final List<CtTypeReference<?>> refs = model.getElements(new TypeFilter<>(CtTypeReference.class));
        for (CtTypeReference<?> tr : refs) {
            if (Objects.equals(tr.getTypeDeclaration(), c)) {
                // tr is a use/instantiation of class 'c'
                this.addObjectInst(tr, ref, obj);
            }
        }

        final List<CtConstructorCall<?>> news = model.getElements(new TypeFilter<>(CtConstructorCall.class));
        for (CtConstructorCall<?> cc : news) {
            final CtTypeReference<?> tr = cc.getType();
            if (Objects.equals(tr.getTypeDeclaration(), c)) {
                // constructor instantiation
                this.addObjectInst(tr, ref, obj);
            }
        }
    }

    public Ref<ObjectInst> addObjectInst(CtTypeReference<?> tr, Ref<ObjectDecl> ref, ObjectDecl obj) {
        // Check it tr is an instantiation or skip.
        final List<CtTypeReference<?>> typeArgs = tr.getActualTypeArguments(); 
        if (typeArgs.size() <= 0) return null;

        // TODO: Implement

        return null;
    }

    //===[ Processors ]=========================================================

    /**
     * performAbstraction will process all the packages and metrics, and
     * resolve anything else that needs to be done to finish the abstraction.
     */
    public void performAbstraction() throws Exception {
        while (!this.pendingPackages.isEmpty()) {
            final CtPackage pkg = this.pendingPackages.iterator().next();
            this.pendingPackages.remove(pkg);
            this.processPackage(pkg);
            this.processPendingMetrics();
        }
        this.consolidateCons();
        this.crossConnectConstructs();
        this.validate();
    }

    // TODO: Implement package imports by deriving from actual type usage
    //       rather than import statements. This will be done in a later step
    //       when the Resolver pipeline is created.

    public Ref<PackageCon> processPackage(CtPackage pkg) throws Exception {
        return this.proj.packages.create(this.log, pkg,
           "package " + SpoonUtils.describeElem(pkg),
            () -> {
                final String name = SpoonUtils.packageName(pkg);
                final String path = SpoonUtils.packagePath(pkg);
                return new PackageCon(name, path);
            },
            (Ref<PackageCon> ref, PackageCon pkgCon) -> {
                for (CtType<?> t : pkg.getTypes()) {
                    if (!SpoonUtils.isObject(t)) pkgCon.add(this.addDeclaration(t));
                }
            });
    }

    private void processPendingMetrics() throws Exception {
        // `addMetrics` may register more methods on `pendingMetrics`
        // so add the current methods, then check if more are pending.
        while (!this.pendingMetrics.isEmpty()) {
            final ArrayList<CtExecutable<?>> methods = new ArrayList<>(this.pendingMetrics);
            this.pendingMetrics.clear();
            for (CtExecutable<?> m : methods) {
                if (m.getBody() == null) {
                    this.log.log("skipping metrics for " + SpoonUtils.describeElem(m) + ": null body");
                    continue;
                }
                if (m.getBody().getStatements().isEmpty()) {
                    this.log.log("skipping metrics for " + SpoonUtils.describeElem(m) + ": empty statement list");
                    continue;
                }

                Ref<MethodDecl> ref = this.proj.methodDecls.getRefByElem(m);
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

    private void consolidateCons() throws Exception {
        this.log.log("Consolidating all constructs");
        this.proj.setToCompareResolved();
        this.proj.setAllIndices();
        while (this.proj.consolidateCons(this.log))
            this.proj.setAllIndices();
        this.proj.setAllIndices();
    }

    private void crossConnectConstructs() throws Exception {
        for (MethodDecl m : this.proj.methodDecls.getConSet()) {
            final PackageCon pkg = m.pkg.mustGetResolved();
            if (pkg == null) this.log.error("package for method is null: " + m);
            final Ref<MethodDecl> decl = this.proj.methodDecls.addOrGetRef(m, "method in package " + pkg);
            pkg.methodDecls.add(decl);
        }

        for (ObjectDecl obj : this.proj.objectDecls.getConSet()) {
            final PackageCon pkg = obj.pkg.mustGetResolved();
            if (pkg == null) this.log.error("package for object is null: " + obj);
            pkg.objectDecls.add(this.proj.objectDecls.addOrGetRef(obj, "object in package " + pkg));
            for (Ref<MethodDecl> met : obj.methodDecls)
                pkg.methodDecls.add(met);
        }
        
        for (InterfaceDecl it : this.proj.interfaceDecls.getConSet()) {
            final PackageCon pkg = it.pkg.mustGetResolved();
            if (pkg == null) this.log.error("package for interface is null: " + it);
            pkg.interfaceDecls.add(this.proj.interfaceDecls.addOrGetRef(it, "interface in package " + pkg));
        }

        for (Value v : this.proj.values.getConSet()) {
            final PackageCon pkg = v.pkg.mustGetResolved();
            if (pkg == null) this.log.error("package for value is null: " + v);
            pkg.values.add(this.proj.values.addOrGetRef(v, "value in package " + pkg));
        }
    }

    private void validate() throws Exception {
        new Validator(this.log, this.proj).validate();
        if (this.log.errorCount() > 0)
            throw new AbstractorException("Errors logged before or during validation.");
    }
}
