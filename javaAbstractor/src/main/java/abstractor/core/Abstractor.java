package abstractor.core;

import java.util.*;

import spoon.Launcher;
import spoon.MavenLauncher;
import spoon.reflect.*;
import spoon.reflect.declaration.*;
import spoon.reflect.path.CtRole;
import spoon.reflect.reference.*;
import spoon.support.compiler.VirtualFile;

import abstractor.core.constructs.*;
import abstractor.core.log.*;
import abstractor.core.validator.*;
import abstractor.core.spoonUtils.*;

public class Abstractor {
    public final Logger  log;
    public final Project proj;

    public final HashSet<CtMethod<?>> pendingMetrics  = new HashSet<>();
    public final HashSet<CtPackage>   pendingPackages = new HashSet<>();

    public Abstractor(Logger log, Project proj) {
        this.log  = log;
        this.proj = proj;
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
        this.pendingPackages.addAll(model.getAllPackages());
    }

    //===[ Construct Adders ]===================================================

    public Ref<PackageCon> addPackage(CtPackage pkg) throws Exception {
        this.pendingPackages.add(pkg);
        return this.proj.packages.addOfGetRefForElem(pkg,
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
        if (elem instanceof CtEnum<?>      e) return this.addObjectDecl(e);
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
        if (decl instanceof CtEnum<?>    e) return this.addMethod(this.addObjectDecl(e), m);
        if (decl instanceof CtClass<?>   c) return this.addMethod(this.addObjectDecl(c), m);
        if (decl instanceof CtInterface<?>) return this.addAbstract(m);

        this.log.error("method has unhandled declaring type: " + SpoonUtils.describeElem(decl));
        return null;
    }

    public Ref<MethodDecl> addMethod(Ref<ObjectDecl> receiver, CtMethod<?> m) throws Exception {
        assert(!SpoonUtils.isObjectMethod(m));
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
        assert(!SpoonUtils.isObjectMethod(m));
        return this.proj.abstracts.create(this.log, m,
            "abstract " + SpoonUtils.describeElem(m),
            () -> {
                final String         name      = m.getSimpleName();
                final Ref<Signature> signature = this.addSignature(m);
                return new Abstract(name, signature);
            });
    }

    public Ref<Signature> addSignature(CtMethod<?> m) throws Exception {
        assert(!SpoonUtils.isObjectMethod(m));
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

    public Ref<InterfaceInst> addArray(CtArrayTypeReference<?> tr) throws Exception {
        final Ref<? extends TypeDesc> elem = this.addTypeDesc(tr.getArrayType());
        Ref<InterfaceInst> ref = this.proj.baker.arrayInst(tr.getSimpleName(), elem);
        this.proj.interfaceInsts.setRefForElem(tr, ref);
        return ref;
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
    
    public Ref<Metrics> addMetrics(CtMethod<?> m) throws Exception {
        return this.proj.metrics.create(this.log, m,
            "metrics " + SpoonUtils.describeElem(m),
            () -> {
                final Location loc = this.proj.locations.create(m.getPosition());
                final Analyzer ana = new Analyzer(this, loc);
                ana.addMethod(m);
                return ana.getMetrics();
            });
    }







    //===[ BELOW NEEDS SOME WORK ]==============================================






    /**
     * Handle Java primitives and object equivalents to primitives (boxed primitives)
     * such that the boxed primitives (e.g. Integer) and String become Basic's.
     * Other types become stub InterfaceDecl.
     */
    public Ref<? extends TypeDesc> addExternalStub(CtTypeReference<?> tr) throws Exception {
        // If a type can not be resolved return an object (kind of like an `any` in Go).
        if (tr == null) return this.proj.baker.anyDesc();

        // Unlike Go's nil that can carry the type, Java's null type has
        // no type associated with it so instead use an object.
        if (SpoonUtils.isNull(tr)) return this.proj.baker.anyDesc();

        // TODO: WHY IS THIS ONLY DOING THE ERASURE?!
        final CtTypeReference<?> erasure = tr.getTypeErasure();
        final Ref<Basic>         boxed   = this.proj.baker.basicForBoxedOrString(erasure);
        if (boxed != null) return boxed;

        final Ref<InterfaceDecl> decl = this.proj.interfaceDecls.create(this.log, erasure,
            "type erasure interface decl " + SpoonUtils.describeElem(erasure),
            () -> {
                final Ref<PackageCon> pkg  = this.addPackageFor(erasure);
                final Location        loc  = this.proj.locations.create(erasure.getPosition());
                final String          name = erasure.getSimpleName();

                // TODO: Should this be "any" since that means it has no methods,
                // or should the stub have abstracts for the methods?
                final Ref<InterfaceDesc> inter = this.proj.baker.anyDesc();

                return new InterfaceDecl(pkg, loc, name, inter, new ArrayList<>());
            });
        
        // TODO: If the stub is not an "any", then the abstract methods may need to have
        //       type parameters and arguments that need to be handled for an instantiation.
        /*
        // Check if the type is a generic instantiation.
        final List<CtTypeReference<?>> typeArgs = tr.getActualTypeArguments();
        if (typeArgs == null || typeArgs.isEmpty()) return decl;

        return this.proj.interfaceInsts.create(this.log, tr,
            "type erasure interface instance " + SpoonUtils.describeElem(tr),
            () -> {
                final ArrayList<Ref<? extends TypeDesc>> instanceTypes = new ArrayList<>(typeArgs.size());
                for (CtTypeReference<?> arg : typeArgs) instanceTypes.add(this.addTypeDesc(arg));
    
                // TODO: decl.getResolved().inter can be null. Figue out another ady to do this!
                return new InterfaceInst(decl, instanceTypes, decl.getResolved().inter);
            });
        */
        return decl;
    }

    /**
     * Handles adding and processing classes, enums, and records.
     * @param c The class to process.
     */
    public Ref<ObjectDecl> addObjectDecl(CtClass<?> c) throws Exception {
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

                // Synthesize the interface description for the class.
                final TreeSet<Ref<Abstract>> abstracts = new TreeSet<Ref<Abstract>>();
                for (CtMethod<?> m : c.getAllMethods()) {
                    if (!m.isStatic() && !SpoonUtils.isObjectMethod(m))
                        abstracts.add(this.addAbstract(m));
                }

                // TODO: FIX BELOW (ref is pin)
                obj.inter = this.proj.interfaceDescs.addOrGetRef(new InterfaceDesc(abstracts, ref), "interface for object");

                // TODO: Finish implementing
                //System.out.println("1) >>> " + c.getSuperInterfaces());

                // Add any nested types.
                for (CtType<?> nt : c.getNestedTypes())
                    this.addTypeDesc(nt.getReference());
            });
    }
    
    public Ref<ObjectDecl> addEnum(CtEnum<?> e) throws Exception {
        return this.proj.objectDecls.create(this.log, e,
            "enum " + SpoonUtils.describeElem(e),
            () -> {
                final String          name = e.getSimpleName();
                final Ref<PackageCon> pkg  = this.addPackageFor(e);
                final Location        loc  = this.proj.locations.create(e.getPosition());

                final CtTypeReference<?> tr = e.getSuperclass();
                // TODO: Move to its own method
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
                // TODO: Finish by adding the "const values" to the package for each enumerator value.
            });
    }

    public Ref<? extends TypeDesc> addTypeDesc(CtTypeReference<?> tr) throws Exception {
        if (tr == null)       return this.proj.baker.anyDesc();
        if (tr.isPrimitive()) return this.addBasic(tr);
        if (tr.isArray())     return this.addArray((CtArrayTypeReference<?>)tr);

        // Handle wildcard types (e.g., ?, ? extends Foo, ? super Bar).
        if (tr instanceof CtWildcardReference wr) return this.addWildcard(wr);

        // Type of the `null` literal in Spoon - not a real external type.
        if (SpoonUtils.isNull(tr)) return this.proj.baker.anyDesc();

        // TODO: NEED TO REEVALUATE ALL OF THIS AI SLOP

        // Use getTypeDeclaration (not getDeclaration) to get shadow types
        // for external/JDK types instead of null.
        CtType<?> ty = null;
        try {
            ty = tr.getTypeDeclaration();
        } catch (Exception ex) {
            this.log.warning("Failed to get type declaration for " + SpoonUtils.describeElem(tr) + ": " + ex.getMessage());
            return this.addExternalStub(tr);
        }

        // If still null, treat as external / unresolvable type.
        if (ty == null) return this.addExternalStub(tr);

        // Annotation types don't participate in data flow. Use an object instead.
        if (ty instanceof CtAnnotationType<?> ann) {
            this.log.notice("Annotation type as type desc (using object): " + SpoonUtils.describeElem(ann));
            return this.proj.baker.anyDesc();
        }

        // Skip anonymous and local types - their code is attributed
        // to the enclosing method (handled in later steps).
        if (tr.isAnonymous() || tr.isLocalType()) {
            // TODO: Fix
            //CtTypeReference<?> superRef = tr.getSuperclass();
            //if (superRef != null && !safeQualifiedName(superRef).equals(objName))
            //    return this.addTypeDesc(superRef);

            //var superIfaces = tr.getSuperInterfaces();
            //if (superIfaces != null && !superIfaces.isEmpty())
            //    return this.addTypeDesc(superIfaces.iterator().next());

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

        this.log.warning("Unhandled type: " + SpoonUtils.describeElem(tr));
        return this.proj.baker.anyDesc();
    }

    public Ref<? extends TypeDesc> addWildcard(CtWildcardReference wr) throws Exception {
        CtTypeReference<?> bound = wr.getBoundingType();
        if (bound == null || bound instanceof CtWildcardReference) return this.proj.baker.anyDesc();

        // Spoon often uses java.lang.Object as the synthetic bound for unbounded "?".
        // Resolving it would pull the entire JDK Object graph into the abstraction.
        try {
            if (SpoonUtils.isObject(bound)) return this.proj.baker.anyDesc();
        } catch (Exception ignored) {
            // The exception here can be ignored since it comes from Spoon
            // failing to get the qualified name and it can be recovered from.
            return this.proj.baker.anyDesc();
        }
        return this.addTypeDesc(bound);
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
                for (CtType<?> t : pkg.getTypes()) pkgCon.add(this.addDeclaration(t));
            });
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

    private void consolidateCons() throws Exception {
        this.proj.setAllIndices();
        while (this.proj.consolidateCons(this.log))
            this.proj.setAllIndices();
        this.proj.setAllIndices();
    }

    private void crossConnectConstructs() throws Exception {
        for (MethodDecl m : this.proj.methodDecls.conSet) {
            final Ref<MethodDecl> decl = this.proj.methodDecls.addOrGetRef(m, "method in package " + m.pkg);
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
