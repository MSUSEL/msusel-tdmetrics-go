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
import abstractor.core.require.Require;
import abstractor.core.validator.*;
import abstractor.core.spoonUtils.*;

public class Abstractor {
    public final Logger  log;
    public final Project proj;
    public final Instantiator instantiator;

    public final HashSet<CtExecutable<?>> pendingMetrics  = new HashSet<>();
    public final HashSet<CtPackage>       pendingPackages = new HashSet<>();

    public CtModel model;

    public Abstractor(Logger log, Project proj) {
        this.log          = log;
        this.proj         = proj;
        this.instantiator = new Instantiator();
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

        final ElementKey elemKey = new ElementKey(pkg);
        final Ref<PackageCon> pkgRef = this.proj.packages.getRefByElem(elemKey);
        if (pkgRef != null) return pkgRef;

        this.log.log("Pending package " + SpoonUtils.describeElem(pkg));
        this.pendingPackages.add(pkg);
        return this.proj.packages.addOrGetRefForElem(elemKey,
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
        if (elem instanceof CtReference ref) {
            final CtElement decl = ref.getDeclaration();
            if (decl == null) return null;
            if (decl instanceof CtClass<?>     c) return this.addObjectInst((CtTypeReference<?>)ref, c);
            if (decl instanceof CtInterface<?> i) return this.addInterfaceInst((CtTypeReference<?>)ref, i);

            this.log.notice("addDeclaration with CtReference and no reference handler: using element");
            elem = decl;
        }

        // Skip annotation types — they don't participate in data flow.
        if (elem instanceof CtAnnotationType<?>) return null;

        // Check CtEnum before CtClass since CtEnum extends CtClass.
        if (elem instanceof CtEnum<?>      e) return this.addEnum(e);
        if (elem instanceof CtClass<?>     c) return this.addObjectDecl(c);
        if (elem instanceof CtInterface<?> i) return this.addInterfaceDecl(i);
        if (elem instanceof CtMethod<?>    m) return this.addMethodDeclOrAbstract(m);

        this.log.error("Unhandled decl: " + SpoonUtils.describeElem(elem));
        return null;
    }

    private Ref<? extends Construct> addMethodDeclOrAbstract(CtMethod<?> m) throws Exception {
        final CtType<?> decl = m.getDeclaringType();
        if (decl instanceof CtEnum<?>    e) return this.addMethodDecl(this.addEnum(e), m);
        if (decl instanceof CtClass<?>   c) return this.addMethodDecl(this.addObjectDecl(c), m);
        if (decl instanceof CtInterface<?>) return this.addAbstract(m);

        this.log.error("Method has unhandled declaring type: " + SpoonUtils.describeElem(decl));
        return null;
    }

    private Ref<? extends Construct> getParent(CtElement elem) throws Exception {
        if (elem.getParent() instanceof CtType<?> parent && parent != null) {
            this.log.log("getting parent type for " + SpoonUtils.describeElem(elem));
            this.log.push();
            try { return this.addDeclaration(parent); }
            finally { this.log.pop(); }
        }
        return null;
    }

    public Ref<InterfaceDecl> addInterfaceDecl(CtInterface<?> i) throws Exception {
        try {
            // All declarations must be added without type arguments.
            this.instantiator.pushCleanFrame();
            return this.proj.interfaceDecls.create(this.log, new ElementKey(i),
                "interface decl " + SpoonUtils.describeElem(i),
                () -> {
                    final String               name       = i.getSimpleName();
                    final Ref<PackageCon>      pkg        = this.addPackageFor(i);
                    final Location             loc        = this.proj.locations.create(i.getPosition());
                    final Ref<InterfaceDesc>   inter      = this.addInterfaceDesc(i);
                    final List<Ref<TypeParam>> typeParams = this.addTypeParams(i);
                    return new InterfaceDecl(pkg, loc, name, inter, typeParams);
                },
                (Ref<InterfaceDecl> ref, InterfaceDecl id) -> {
                    id.setVisibility(i);
                    final Ref<? extends Construct> parent = this.getParent(i);
                    if (parent != null) {
                        this.log.notice(parent.toString()+" => "+SpoonUtils.describeElem(i)); // TODO: Finish
                    }
                    // Add any nested types.
                    for (CtType<?> nt : i.getNestedTypes())
                        this.addTypeDesc(nt.getReference()); // TODO: Finish
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    public Ref<? extends TypeDesc> addInterfaceInst(CtTypeReference<?> tr, CtInterface<?> i) throws Exception {
        final Ref<InterfaceDecl> decl = this.addInterfaceDecl(i);
        if (!i.isGenerics()) return decl;

        final List<Ref<TypeParam>> typeParams = this.addTypeParams(i);
        final ArrayList<Ref<? extends TypeDesc>> typeArgs = this.addTypeArguments(tr, typeParams);
        if (typeArgs == null) return decl;

        try {
            this.instantiator.pushFrame();
            for (int j = 0; j < typeParams.size(); j++)
                this.instantiator.add(typeParams.get(j), typeArgs.get(j));

            return this.proj.interfaceInsts.create(this.log, new ElementKey(tr, this.instantiator.typeArgs()),
                "interface instantiation "+SpoonUtils.describeGeneric(tr),
                () -> {
                    final Ref<InterfaceDesc> resolved = this.addInterfaceDesc(i);
                    return new InterfaceInst(decl, this.instantiator.typeArgs(), resolved);
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    public Ref<InterfaceDesc> addInterfaceDesc(CtInterface<?> i) throws Exception {
        return this.proj.interfaceDescs.create(this.log, new ElementKey(i, this.instantiator.typeArgs()),
            "interface description " + SpoonUtils.describeElem(i),
            () -> {
                final TreeSet<Ref<Abstract>> abstracts = new TreeSet<Ref<Abstract>>();
                for (CtMethod<?> m : i.getAllMethods()) {
                    if (!SpoonUtils.isObjectMethod(m)) abstracts.add(this.addAbstract(m));
                }

                Ref<? extends Construct> pin = null;
                if (i.getRoleInParent() == CtRole.NESTED_TYPE) {
                    final CtElement parent = i.getParent();
                    if (parent instanceof CtTypeReference<?> nest && nest != null) {
                        pin = this.addTypeDesc(nest);
                    } else {
                        this.log.error("Unhandled nested interface decl " + SpoonUtils.describeElem(i) + " in " + parent);
                    }
                }

                return new InterfaceDesc(abstracts, pin);
            },
            (Ref<InterfaceDesc> ref, InterfaceDesc id) -> {
                // Add direct super-interfaces this interface extends
                for (CtTypeReference<?> supRef : i.getSuperInterfaces()) {
                    final CtType<?> supDecl = supRef.getTypeDeclaration(); // may be null for shadow/unresolved
                    if (supDecl != null && supDecl instanceof CtInterface<?> supId && supId != null) {
                        id.inherits.add(this.addInterfaceDesc(supId));
                    } else {
                        this.log.error("Unhandled super-interface " + SpoonUtils.describeElem(supDecl) + " for " + id);
                    }
                }
            });
    }

    public Ref<MethodDecl> addMethodDecl(Ref<ObjectDecl> receiver, CtMethod<?> m) throws Exception {
        Require.notObjectMethod(m);
        final ObjectDecl recv = receiver.mustGetResolved();
        try {
            // All declarations must be added without type arguments.
            this.instantiator.pushCleanFrame();
            return this.proj.methodDecls.create(this.log, new ElementKey(m),
                "method " + SpoonUtils.describeElem(m),
                () -> {
                    final Ref<PackageCon>      pkg        = recv.pkg;
                    final Location             loc        = this.proj.locations.create(m.getPosition());
                    final String               name       = m.getSimpleName();
                    final Ref<Signature>       signature  = this.addSignature(m);
                    final List<Ref<TypeParam>> typeParams = this.addTypeParams(m);
                    final MethodDecl md = new MethodDecl(pkg, receiver, loc, name, signature, typeParams);
                    md.isStatic = m.isStatic();
                    return md;
                },
                (Ref<MethodDecl> ref, MethodDecl md) -> {
                    md.setVisibility(m);
                    recv.methodDecls.add(ref);
                    this.pendingMetrics.add(m);
                    final Ref<? extends Construct> parent = this.getParent(m);
                    if (parent != null) {
                        this.log.notice(parent.toString()+" => "+SpoonUtils.describeElem(m)); // TODO: Finish
                    }
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    public Ref<MethodInst> addMethodInst(Ref<ObjectInst> receiver, CtMethod<?> m) throws Exception {
        Require.notObjectMethod(m);
        final ObjectInst recv = receiver.mustGetResolved();
        return this.proj.methodInsts.create(this.log, new ElementKey(m, this.instantiator.typeArgs()),
            "method instantiation " + SpoonUtils.describeElem(m),
            () -> {
                final Ref<MethodDecl>               generic       = this.addMethodDecl(recv.generic, m);
                final List<Ref<? extends TypeDesc>> instanceTypes = this.instantiator.typeArgs();
                final Ref<Signature>                resolved      = this.addSignature(m);
                return new MethodInst(generic, receiver, instanceTypes, resolved);
            },
            (Ref<MethodInst> ref, MethodInst mi) -> {
                recv.methods.add(ref);
            });
    }

    public Ref<Abstract> addAbstract(CtMethod<?> m) throws Exception {
        Require.notObjectMethod(m);
        return this.proj.abstracts.create(this.log, new ElementKey(m, this.instantiator.typeArgs()),
            "abstract " + SpoonUtils.describeElem(m),
            () -> {
                final String         name      = m.getSimpleName();
                final Ref<Signature> signature = this.addSignature(m);
                return new Abstract(name, signature);
            });
    }

    public Ref<Signature> addSignature(CtMethod<?> m) throws Exception {
        Require.notObjectMethod(m);
        return this.proj.signatures.create(this.log, new ElementKey(m, this.instantiator.typeArgs()),
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

    public Ref<MethodDecl> addMethodDeclForConstructor(Ref<ObjectDecl> receiver, CtConstructor<?> ctor) throws Exception {
        try {
            // All declarations must be added without type arguments.
            this.instantiator.pushCleanFrame();
            return this.proj.methodDecls.create(log, new ElementKey(ctor),
                "constructor " + ctor.getSignature(),
                () -> {
                    final ObjectDecl           recv       = receiver.mustGetResolved();
                    final Ref<PackageCon>      pkg        = recv.pkg;
                    final Location             loc        = this.proj.locations.create(ctor.getPosition());
                    final String               name       = recv.name;
                    final Ref<Signature>       signature  = this.addSignatureForConstructor(ctor);
                    final List<Ref<TypeParam>> typeParams = this.addTypeParams(ctor);
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
        } finally {
            this.instantiator.popFrame();
        }
    }

    public Ref<Signature> addSignatureForConstructor(CtConstructor<?> m) throws Exception {
        return this.proj.signatures.create(this.log, new ElementKey(m, this.instantiator.typeArgs()),
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
        return this.proj.arguments.create(this.log, new ElementKey(p, this.instantiator.typeArgs()),
            "parameter " + SpoonUtils.describeElem(p),
            () -> {
                final String                  name = p.getSimpleName();
                final Ref<? extends TypeDesc> type = this.addTypeDesc(p.getType());
                return new Argument(name, type);
            });
    }
    
    public Ref<Argument> addArgument(CtTypeReference<?> p) throws Exception {
        return this.proj.arguments.create(this.log, new ElementKey(p, this.instantiator.typeArgs()),
            "parameter <unnamed> " + SpoonUtils.describeGeneric(p),
            () -> {
                final Ref<? extends TypeDesc> type = this.addTypeDesc(p);
                return new Argument("", type);
            });
    }
    
    public Ref<StructDesc> addStructDesc(CtType<?> c) throws Exception {
        return this.proj.structDescs.create(this.log, new ElementKey(c, this.instantiator.typeArgs()),
            "struct " + SpoonUtils.describeElem(c),
            () -> {
                // Collect all fields.
                final ArrayList<Ref<Field>> fields = new ArrayList<>();
                for (CtFieldReference<?> fr : c.getAllFields())
                    fields.add(this.addField(fr.getFieldDeclaration()));

                // Add extended class as a "$super" field.
                final CtTypeReference<?> superFr = c.getSuperclass();
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
        return this.proj.fields.create(this.log, new ElementKey(f, this.instantiator.typeArgs()),
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
        return this.proj.fields.create(this.log, new ElementKey(f, this.instantiator.typeArgs()),
            "field " + name,
            () -> {
                final Ref<? extends TypeDesc> type = this.addTypeDesc(f);
                return new Field(name, type);
            });
    }

    public Ref<Selection> addSelection(CtField<?> field) throws Exception {
        return this.proj.selections.create(this.log, new ElementKey(field, this.instantiator.typeArgs()),
            "select field " + SpoonUtils.describeElem(field),
            () -> {
                final String name = field.getSimpleName();

                // TODO: Is this the correct way to get the decl? Does it need to be the instantiated type?

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
        final ElementKey elemKey = new ElementKey(tr);
        return this.proj.interfaceInsts.setRefForElem(elemKey, ref);
    }
    
    public Ref<Basic> addBasic(CtTypeReference<?> tr) throws Exception {
        return this.proj.basics.create(this.log, new ElementKey(tr, this.instantiator.typeArgs()),
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

    private List<Ref<TypeParam>> addTypeParams(CtElement elem) throws Exception {
        final List<Ref<TypeParam>> result =
            (elem.getParent() instanceof CtType<?> parent && parent != null)
            ? this.addTypeParams(parent)
            : new ArrayList<>();

        if (elem instanceof CtFormalTypeDeclarer td) {
            for (CtTypeParameter tp : td.getFormalCtTypeParameters()) {
                result.add(this.addTypeParam(tp));
            }
        }
        return result;
    }

    public Ref<TypeParam> addTypeParam(CtTypeParameter tp) throws Exception {
        // Do not use type arguments in the ElementKey for typeParams.
        // The typeParams will be replaced by the instantiator later.
        return this.proj.typeParams.create(this.log, new ElementKey(tp, null),
            "type params " + SpoonUtils.describeElem(tp),
            () -> {
                final String                  name = tp.getSimpleName();
                final CtTypeReference<?>      tr   = tp.getTypeErasure();
                final Ref<? extends TypeDesc> type = this.addTypeDesc(tr);
                return new TypeParam(name, type);
            });
    }
    
    public Ref<Metrics> addMetrics(CtExecutable<?> m) throws Exception {
        return this.proj.metrics.create(this.log, new ElementKey(m, this.instantiator.typeArgs()),
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
        try {
            // All declarations must be added without type arguments.
            this.instantiator.pushCleanFrame();
            return this.proj.objectDecls.create(this.log, new ElementKey(c),
                "object decl " + SpoonUtils.describeElem(c),
                () -> {
                    final Ref<PackageCon>      pkg        = this.addPackageFor(c);
                    final Location             loc        = this.proj.locations.create(c.getPosition());
                    final String               name       = c.getSimpleName();
                    final Ref<StructDesc>      struct     = this.addStructDesc(c);
                    final List<Ref<TypeParam>> typeParams = this.addTypeParams(c);
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
                            this.addMethodDeclForConstructor(ref, ctor);
                        }
                    }

                    // Add methods for the class.
                    for (CtMethod<?> m : c.getAllMethods()) {
                        if (m.getParent().equals(c) && !SpoonUtils.isObjectMethod(m))
                            this.addMethodDecl(ref, m);
                    }

                    obj.inter = this.synthesizeObjectInterface(c, ref);

                    final Ref<? extends Construct> parent = this.getParent(c);
                    if (parent != null) {
                        this.log.notice(parent.toString()+" => "+SpoonUtils.describeElem(c)); // TODO: Finish
                    }

                    // Add any nested types.
                    //for (CtType<?> nt : c.getNestedTypes()) // TODO: Do we need more for nested types?
                    //    this.addTypeDesc(nt.getReference());
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    private Ref<InterfaceDesc> synthesizeObjectInterface(CtClass<?> c, Ref<? extends Construct> ref) throws Exception {
        // Synthesize the interface abstractions for the class.
        final TreeSet<Ref<Abstract>> abstracts = new TreeSet<Ref<Abstract>>();
        for (CtMethod<?> m : c.getAllMethods()) {
            if (!m.isStatic() && !SpoonUtils.isObjectMethod(m))
                abstracts.add(this.addAbstract(m));
        }

        // Synthesize the interface description for the class.
        if (abstracts.size() > 0 || c.getSuperInterfaces().size() > 0) {
            final InterfaceDesc it = new InterfaceDesc(abstracts, ref);
            final Ref<InterfaceDesc> inter = this.proj.interfaceDescs.addOrGetRef(it, this.instantiator.typeArgs(), "interface for object");

            // Add direct super-interfaces this object extends.
            for (CtTypeReference<?> supRef : c.getSuperInterfaces()) {
                final CtType<?> supDecl = supRef.getTypeDeclaration(); // may be null for shadow/unresolved
                if (supDecl != null && supDecl instanceof CtInterface<?> supId && supId != null) {
                    it.inherits.add(this.addInterfaceDesc(supId));
                } else {
                    this.log.error("Unhandled super-interface " + SpoonUtils.describeElem(supDecl) + " for " + ref);
                }
            }
            return inter;
        }
        return this.proj.baker.anyDesc();
    }

    public Ref<? extends TypeDesc> addObjectInst(CtTypeReference<?> tr, CtClass<?> c) throws Exception {
        final Ref<ObjectDecl> decl = this.addObjectDecl(c);
        if (!c.isGenerics()) return decl;

        final List<Ref<TypeParam>> typeParams = this.addTypeParams(c);
        final ArrayList<Ref<? extends TypeDesc>> typeArgs = this.addTypeArguments(tr, typeParams);
        if (typeArgs == null) return decl;

        this.instantiator.pushFrame();
        for (int i = 0; i < typeParams.size(); i++)
            this.instantiator.add(typeParams.get(i), typeArgs.get(i));

        try {
            return this.proj.objectInsts.create(this.log, new ElementKey(tr, this.instantiator.typeArgs()),
                "object instantiation "+SpoonUtils.describeGeneric(tr),
                () -> {                    
                    final Ref<StructDesc> resData = this.addStructDesc(c);
                    final Ref<InterfaceDesc> resInterface = this.synthesizeObjectInterface(c, null);
                    return new ObjectInst(decl, this.instantiator.typeArgs(), resData, resInterface);
                },
                (Ref<ObjectInst> ref, ObjectInst obj) -> {
                    // Add methods for the class instantiation.
                    for (CtMethod<?> m : c.getAllMethods()) {
                        if (m.getParent().equals(c) && !SpoonUtils.isObjectMethod(m))
                            this.addMethodInst(ref, m);
                    }

                    // Add any nested types.
                    for (CtType<?> nt : c.getNestedTypes()) // TODO: Do we need more for nested types?
                        this.addTypeDesc(nt.getReference());
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    /**
     * This adds the type arguments from the type reference.
     * @param tr The type reference for the possible instantiation
     * @param typeParams The type parameters from the interface, method, or object.
     * @return The list of type arguments or null if there is no instantiation.
     */
    private ArrayList<Ref<? extends TypeDesc>> addTypeArguments(CtTypeReference<?> tr, List<Ref<TypeParam>> typeParams) throws Exception {
        final List<CtTypeReference<?>> ctTypeArgs = tr.getActualTypeArguments();
        if (ctTypeArgs == null) return null;

        final int count = ctTypeArgs.size();
        if (count <= 0) return null;
        if (count != typeParams.size()) return null;

        final ArrayList<Ref<? extends TypeDesc>> typeArgs = new ArrayList<>();
        for (CtTypeReference<?> ctTypeArg : ctTypeArgs)
            typeArgs.add(this.addTypeDesc(ctTypeArg));

        for (int i = 0; i < count; i++) {
            if (!typeArgs.get(i).equals(typeParams.get(i))) {
                // There was a difference so there is an instantiation
                return typeArgs;
            }
        }
        // There was no difference so the instantiation is not useful.
        return null;
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
        final CtType<?> ty = tr.getTypeDeclaration();
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

        // Handle type parameters by checking if there is a type argument replacement
        // when defining an instantiation instead of a generic. 
        if (ty instanceof CtTypeParameter tp)
            return this.instantiator.replace(this.addTypeParam(tp));

        // Check CtEnum before CtClass since CtEnum extends CtClass.
        if (ty instanceof CtEnum<?>      e) return this.addEnum(e);
        if (ty instanceof CtClass<?>     c) return this.addObjectInst(tr, c);
        if (ty instanceof CtInterface<?> i) return this.addInterfaceInst(tr, i);

        this.log.warning("Unhandled type description: " + SpoonUtils.describeElem(ty));
        return null;
    }

    public Ref<? extends TypeDesc> addWildcard(CtWildcardReference wr) throws Exception {
        final CtTypeReference<?> bound = wr.getBoundingType();
        // Spoon often uses java.lang.Object as the synthetic bound for unbounded "?".
        // Resolving it would pull the entire JDK Object graph into the abstraction.
        if (bound == null || bound instanceof CtWildcardReference || SpoonUtils.isObject(bound))
            return this.proj.baker.anyDesc();
        return this.addTypeDesc(bound);
    }

    public Ref<ObjectDecl> addEnum(CtEnum<?> e) throws Exception {
        try {
            // All declarations must be added without type arguments.
            this.instantiator.pushCleanFrame();
            return this.proj.objectDecls.create(this.log, new ElementKey(e),
                "enum " + SpoonUtils.describeElem(e),
                () -> {
                    final String             name   = e.getSimpleName();
                    final Ref<PackageCon>    pkg    = this.addPackageFor(e);
                    final Location           loc    = this.proj.locations.create(e.getPosition());
                    final CtTypeReference<?> tr     = e.getSuperclass();
                    final Ref<StructDesc>    struct = this.proj.structDescs.create(this.log, new ElementKey(tr),
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
                        this.proj.values.create(this.log, new ElementKey(e),
                            "enum value "+ SpoonUtils.describeElem(ev),
                            () -> {
                                final String   name = ev.getSimpleName();
                                final Location loc  = this.proj.locations.create(ev.getPosition());
                                return new Value(od.pkg, loc, name, true, null, ref);
                            });
                    }

                    // TODO: Add methods that have been added to the enum?
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    public Ref<? extends TypeDesc> addShadowTypeDesc(CtTypeReference<?> tr) throws Exception {
        // from isShadow() method:
        // > When an element isn't present in the factory (created in another factory),
        // > this element is considered as "shadow". e.g., a shadow element can be a
        // > CtType of java.lang.Class built when we call CtTypeReference.getTypeDeclaration()
        // > on a reference of java.lang.Class."
        return this.proj.baker.anyDesc();
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
        return this.proj.packages.create(this.log, new ElementKey(pkg),
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
                    // Use the following for debugging, but it is commented out since it can be noisy.
                    //this.log.log("skipping metrics for " + SpoonUtils.describeElem(m) + ": null body");
                    continue;
                }
                if (m.getBody().getStatements().isEmpty()) {
                    // Use the following for debugging, but it is commented out since it can be noisy.
                    //this.log.log("skipping metrics for " + SpoonUtils.describeElem(m) + ": empty statement list");
                    continue;
                }

                final ElementKey elemKey = new ElementKey(m);
                final Ref<MethodDecl> ref = this.proj.methodDecls.getRefByElem(elemKey);
                if (!ref.isResolved())
                    throw new AbstractorException("Expected " + ref + " to be resolved before processing pending metrics.");

                final MethodDecl md = ref.getResolved();
                if (md.metrics != null)
                    throw new AbstractorException("The metrics for " + md + " have already been processed before " + m.getSimpleName() + ".");

                final Ref<Metrics> metRef = this.addMetrics(m);
                final Metrics met = metRef.getResolved();
                if (met.hasBody()) md.metrics = metRef;
                else {
                    // remove the reference and metrics from factory since bodiless methods can be ignored.
                    this.proj.metrics.removeElem(this.log, elemKey, "metrics " + m.getSimpleName());
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
            final Ref<MethodDecl> decl = this.proj.methodDecls.addOrGetRef(m, null, "method in package " + pkg);
            pkg.methodDecls.add(decl);
        }

        for (ObjectDecl obj : this.proj.objectDecls.getConSet()) {
            final PackageCon pkg = obj.pkg.mustGetResolved();
            if (pkg == null) this.log.error("package for object is null: " + obj);
            pkg.objectDecls.add(this.proj.objectDecls.addOrGetRef(obj, null, "object in package " + pkg));
            for (Ref<MethodDecl> met : obj.methodDecls)
                pkg.methodDecls.add(met);
        }
        
        for (InterfaceDecl it : this.proj.interfaceDecls.getConSet()) {
            final PackageCon pkg = it.pkg.mustGetResolved();
            if (pkg == null) this.log.error("package for interface is null: " + it);
            pkg.interfaceDecls.add(this.proj.interfaceDecls.addOrGetRef(it, null, "interface in package " + pkg));
        }

        for (Value v : this.proj.values.getConSet()) {
            final PackageCon pkg = v.pkg.mustGetResolved();
            if (pkg == null) this.log.error("package for value is null: " + v);
            pkg.values.add(this.proj.values.addOrGetRef(v, null, "value in package " + pkg));
        }
    }

    private void validate() throws Exception {
        new Validator(this.log, this.proj).validate();
        if (this.log.errorCount() > 0)
            throw new AbstractorException("Errors logged before or during validation.");
    }
}
