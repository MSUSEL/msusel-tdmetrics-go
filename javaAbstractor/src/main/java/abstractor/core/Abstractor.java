package abstractor.core;

import java.io.*;
import java.util.*;

import org.apache.maven.model.Model;

import spoon.Launcher;
import spoon.MavenLauncher;
import spoon.reflect.*;
import spoon.reflect.code.*;
import spoon.reflect.declaration.*;
import spoon.reflect.reference.*;
import spoon.support.compiler.SpoonPom;
import spoon.support.compiler.VirtualFile;

import abstractor.core.cmp.*;
import abstractor.core.constructs.*;
import abstractor.core.iter.Bundle;
import abstractor.core.log.*;
import abstractor.core.json.*;
import abstractor.core.require.Require;
import abstractor.core.spoonUtils.*;
import abstractor.core.tools.*;

public class Abstractor {
    public final Logger  log;
    public final Project proj;
    public final boolean skipProjInfo;
    public final Instantiator instantiator;

    public final Set<CtExecutable<?>> pendingMetrics  = Collections.newSetFromMap(new IdentityHashMap<>());
    public final Set<CtPackage>       pendingPackages = Collections.newSetFromMap(new IdentityHashMap<>());

    public CtModel model;

    public Abstractor(Logger log, Project proj, boolean skipProjInfo) {
        this.log          = log;
        this.proj         = proj;
        this.skipProjInfo = skipProjInfo;
        this.instantiator = new Instantiator();
    }

    /**
     * Reads a project containing a pom.xml maven file.
     * @param mavenProject The path to the project folder containing a pom.xml
     */
    public void prepareMavenProject(String mavenProject) throws Exception {
        this.log.log("Reading " + mavenProject);
        SpoonUtils.addKnownPathRoot(mavenProject);
        this.setCommitHash(mavenProject);

        MavenLauncher launcher = new MavenLauncher(mavenProject, MavenLauncher.SOURCE_TYPE.APP_SOURCE);
        launcher.getEnvironment().setComplianceLevel(17);
        launcher.getEnvironment().setNoClasspath(true);
        CtModel model = launcher.buildModel();
        if (model.getAllTypes().size() > 0) {
            this.setProjectInfo(launcher);
            this.prepareModel(model);
            return;
        }

        // If the model couldn't be loaded (it has no types in it) from the app
        // source alone then try again with the maven project path as an input
        // resource. The input resource can't always be add, otherwise it will
        // cause duplicate identifiers in some projects. For the integration
        // tests in testData/java, the input resource is not needed. I have no clue
        // what the difference is between the maven models to require this
        // but if it works, I'm not going to fix it right now.
        launcher = new MavenLauncher(mavenProject, MavenLauncher.SOURCE_TYPE.APP_SOURCE);
        launcher.addInputResource(mavenProject);
        model = launcher.buildModel();
        this.setProjectInfo(launcher);
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

    private void setProjectInfo(MavenLauncher launcher) {
        if (this.skipProjInfo) return;
        SpoonPom pom = launcher.getPomFile();
        if (pom != null) {
            Model m = pom.getModel();
            this.proj.groupId    = m.getGroupId();
            this.proj.artifactId = m.getArtifactId();
            this.proj.version    = m.getVersion();
            this.proj.name       = m.getName();
        }
    }

    private void setCommitHash(String repoDir) {
        if (this.skipProjInfo) return;
        try {
            Process p = new ProcessBuilder("git", "rev-parse", "HEAD")
                .directory(new File(repoDir))
                .redirectErrorStream(true)
                .start();
            try (BufferedReader r = new BufferedReader(new InputStreamReader(p.getInputStream()))) {
                p.waitFor();
                if (p.exitValue() != 0) return;
                String line = r.readLine();
                if (line != null)
                    this.proj.commitHash = line.trim();
            }
        } catch (Exception e) {
            return;
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
        if (elem instanceof CtEnum<?>        e) return this.addEnum(e);
        if (elem instanceof CtClass<?>       c) return this.addObjectDecl(c);
        if (elem instanceof CtInterface<?>   i) return this.addInterfaceDecl(i);
        if (elem instanceof CtMethod<?>      m) return this.addMethodDeclOrAbstract(m);
        if (elem instanceof CtConstructor<?> c) return this.addMethodDeclForConstructor(c);

        this.log.error("Unhandled decl: " + SpoonUtils.describeElem(elem));
        return null;
    }

    private Ref<? extends Construct> addMethodDeclOrAbstract(CtMethod<?> m) throws Exception {        
        if (SpoonUtils.isObjectMethod(m)) return null;

        final CtType<?> decl = m.getDeclaringType();
        if (decl.isAnonymous()) {
            this.log.notice("Ignoring method of an anonymous declaring type: " + SpoonUtils.describeElem(m) + " in " + SpoonUtils.describeElem(decl));
            return null;
        }
        if (decl.isLocalType()) {
            this.log.notice("Ignoring method of a local declaring type: " + SpoonUtils.describeElem(m) + " in " + SpoonUtils.describeElem(decl));
            return null;
        }
        if (m.isImplicit()) {
            this.log.notice("Ignoring implicit method: " + SpoonUtils.describeElem(m) + " in " + SpoonUtils.describeElem(decl));
            return null;
        }

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
            return this.proj.interfaceDecls.create(this.log, this.instantiator,
                new ElementKey(i),
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
                    id.setNest(this.getParent(i));

                    if (i.isShadow()) {
                        // Shadow interface: only expose members that are actually used;
                        // nested types are pulled in on demand via addTypeDesc.
                        id.isShadow = true;
                        return;
                    }

                    for (CtType<?> nt : i.getNestedTypes())
                        id.nestedTypes.add(this.addTypeDesc(nt.getReference()));
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    public Ref<? extends TypeDesc> addInterfaceInst(CtTypeReference<?> tr, CtInterface<?> i) throws Exception {
        final Ref<InterfaceDecl> decl = this.addInterfaceDecl(i);
        if (!SpoonUtils.isGenerics(i)) return decl;

        final List<Ref<TypeParam>> typeParams = this.addTypeParams(i);
        final List<Ref<? extends TypeDesc>> typeArgs = this.addTypeArguments(tr, typeParams);
        if (typeArgs == null) return decl;

        try {
            final boolean definedInNest = SpoonUtils.inSameNested(tr, i);
            if (definedInNest) this.instantiator.pushFrame();
            else this.instantiator.pushCleanFrame();
            for (int j = 0; j < typeParams.size(); j++)
                this.instantiator.add(typeParams.get(j), typeArgs.get(j));

            return this.proj.interfaceInsts.create(this.log, this.instantiator,
                new ElementKey(tr, this.instantiator.typeArgs()),
                "interface instantiation "+SpoonUtils.describeGeneric(tr),
                () -> {
                    final Ref<InterfaceDesc> resolved = this.addInterfaceDesc(i);
                    final List<Ref<? extends TypeDesc>> argTypes = this.instantiator.typeArgs();
                    return new InterfaceInst(decl, argTypes, resolved);
                }, 
                (Ref<InterfaceInst> ref, InterfaceInst it) -> {
                    if (i.isShadow()) return; // shadows: nested types are pulled in on demand

                    // Create instances for all nested types too.
                    for (CtType<?> nt : i.getNestedTypes())
                        this.addTypeDesc(nt.getReference());
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    public Ref<InterfaceDesc> addInterfaceDesc(CtInterface<?> i) throws Exception {
        return this.proj.interfaceDescs.create(this.log, this.instantiator,
            new ElementKey(i, this.instantiator.typeArgs()),
            "interface description " + SpoonUtils.describeElem(i),
            (Ref<InterfaceDesc> ref) -> {
                final TreeSet<Ref<Abstract>> abstracts = new TreeSet<Ref<Abstract>>();
                if (!i.isShadow()) {
                    // Shadow interfaces: abstracts are attached on demand from addAbstract.
                    for (CtMethod<?> m : i.getAllMethods()) {
                        if (!m.isStatic() && !SpoonUtils.isObjectMethod(m))
                            abstracts.add(this.addAbstract(m));
                    }
                }

                Ref<? extends Construct> pin = null;
                if (SpoonUtils.isNested(i)) {
                    final CtElement parent = i.getParent();
                    if (parent instanceof CtTypeReference<?> nest) {
                        pin = this.addTypeDesc(nest);
                    } else if (parent instanceof CtType<?> nest) {
                        pin = this.addTypeDesc(nest.getReference());
                    } else {
                        this.log.warning("Unhandled nested interface decl " + SpoonUtils.describeElem(i) + " in " + SpoonUtils.describeElem(parent));
                    }
                }

                // If pin is still null for a shadow, then pin it to itself.
                if (i.isShadow() && pin == null) pin = ref;

                return new InterfaceDesc(abstracts, pin);
            },
            (Ref<InterfaceDesc> ref, InterfaceDesc id) -> {
                if (i.isShadow()) return; // shadows: super-interfaces are not proactively walked
                // Add direct super-interfaces this interface extends.
                for (CtTypeReference<?> supRef : i.getSuperInterfaces()) {
                    final CtType<?> supDecl = supRef.getTypeDeclaration(); // may be null for shadow/unresolved
                    if (supDecl == null) {
                        this.log.warning("Unhandled null super-interface for " + id);
                        continue;
                    }
                    if (!(supDecl instanceof CtInterface<?> supIt)) {
                        this.log.error("Unhandled super-interface " + SpoonUtils.describeElem(supDecl) + " for " + id);
                        continue;
                    }
                    id.inherits.add(this.addInterfaceDesc(supIt));
                }
            });
    }

    public Ref<MethodDecl> addMethodDecl(Ref<ObjectDecl> receiver, CtMethod<?> m) throws Exception {
        Require.notObjectMethod(m);
        Require.require(!m.isImplicit(), SpoonUtils.describeElem(m) + " is implicit");
        final ObjectDecl recv = receiver.mustGetResolved();
        try {
            // All declarations must be added without type arguments.
            this.instantiator.pushCleanFrame();
            return this.proj.methodDecls.create(this.log, this.instantiator,
                new ElementKey(m),
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
                    //md.setNest(this.getParent(m)); // Not needed because of receiver
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    public Ref<MethodInst> addMethodInstForObjectInst(Ref<ObjectInst> receiver, CtMethod<?> m) throws Exception {
        Require.notObjectMethod(m);
        try {
            this.instantiator.pushFrame();
            final ObjectInst recv = receiver.mustGetResolved();
            final Ref<MethodDecl> generic = this.addMethodDecl(recv.generic, m);
            final List<Ref<TypeParam>> typeParams = this.addTypeParams(m);
            for (int i = 0; i < typeParams.size(); i++)
                this.instantiator.add(typeParams.get(i), typeParams.get(i));

            return this.proj.methodInsts.create(this.log, this.instantiator,
                new ElementKey(m, this.instantiator.typeArgs()),
                "method for object instantiation " + SpoonUtils.describeElem(m),
                () -> {
                    final List<Ref<? extends TypeDesc>> instanceTypes = this.instantiator.typeArgs();
                    final Ref<Signature>                resolved      = this.addSignature(m);
                    return new MethodInst(generic, receiver, instanceTypes, resolved);
                },
                (Ref<MethodInst> ref, MethodInst mi) -> {
                    recv.methods.add(ref);
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    public Ref<MethodInst> addMethodInstForObjectInst(Ref<ObjectInst> receiver, CtConstructor<?> ctor) throws Exception {
        if (ctor.isImplicit()) return null;
        try {
            this.instantiator.pushFrame();
            final ObjectInst recv = receiver.mustGetResolved();
            final Ref<MethodDecl> generic = this.addMethodDeclForConstructor(recv.generic, ctor);
            final List<Ref<TypeParam>> typeParams = this.addTypeParams(ctor);
            for (int i = 0; i < typeParams.size(); i++)
                this.instantiator.add(typeParams.get(i), typeParams.get(i));

            return this.proj.methodInsts.create(this.log, this.instantiator,
                new ElementKey(ctor, this.instantiator.typeArgs()),
                "constructor for object instantiation " + SpoonUtils.describeElem(ctor),
                () -> {
                    final List<Ref<? extends TypeDesc>> instanceTypes = this.instantiator.typeArgs();
                    final Ref<Signature>                resolved      = this.addSignatureForConstructor(ctor);
                    return new MethodInst(generic, receiver, instanceTypes, resolved);
                },
                (Ref<MethodInst> ref, MethodInst mi) -> {
                    recv.methods.add(ref);
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    /**
     * addMethodInstForCall creates or fetches a MethodInst that captures the
     * actual type arguments at a call site. Falls back to returning the
     * generic MethodDecl (or Abstract) when the call cannot be narrowed to a
     * useful instantiation (no generics at all, missing/inferred type info,
     * interface-declared method, etc).
     *
     * The returned ref is what the caller should use as the invocation edge.
     */
    public Ref<? extends Construct> addMethodInstForCall(CtInvocation<?> in) throws Exception {
        final CtExecutableReference<?> er = in.getExecutable();
        if (er == null) return null;
        if (er.isImplicit()) return null;

        final CtExecutable<?> ex = er.getExecutableDeclaration();
        if (ex == null) return null;

        if (!(ex instanceof CtMethod<?> m)) return null; // caller handles ctors / others
        if (SpoonUtils.isObjectMethod(m)) return null;

        // Fall back to the plain decl path when the receiver is anything other
        // than a tracked class (interfaces produce Abstracts, not method instances).
        final CtType<?> declType = m.getDeclaringType();
        if (!(declType instanceof CtClass<?> declClass))
            return this.addMethodDeclOrAbstract(m);

        // Collect class-level and method-level actual type args from Spoon.
        final CtTypeReference<?> receiverRef = getCallReceiverTypeRef(in, declClass);
        final ArrayList<CtTypeReference<?>> classArgs = new ArrayList<>();
        this.collectActualTypeArgs(receiverRef, classArgs);
        final List<CtTypeReference<?>> methodArgs = in.getActualTypeArguments();

        // If neither the class nor the method are generic, no MethodInst is useful.
        final List<Ref<TypeParam>> classParams  = this.addTypeParams(declClass);
        final List<Ref<TypeParam>> methodParams = new ArrayList<>();
        for (CtTypeParameter tp : m.getFormalCtTypeParameters()) {
            Ref<TypeParam> tpRef = this.addTypeParam(tp);
            Require.notNull(tpRef, "type parameter for " + SpoonUtils.describeElem(tp) +
                " in method " + SpoonUtils.describeElem(m) + " may not be null");
            methodParams.add(tpRef);
        }
        if (classParams.isEmpty() && methodParams.isEmpty()) return this.addDeclaration(m);

        // Spoon didn't hand us enough info to bind every param — fall back.
        if (!classParams.isEmpty()  && classArgs.size()  != classParams.size())  return this.addDeclaration(m);
        if (!methodParams.isEmpty() && methodArgs.size() != methodParams.size()) return this.addDeclaration(m);

        // Resolve args to type descriptors in the ambient frame (so a call-site
        // arg that is itself a type-param `S` resolves to the TypeParam ref for `S`).
        final List<Ref<? extends TypeDesc>> tdClassArgs  = this.addTypeArguments(classArgs);
        final List<Ref<? extends TypeDesc>> tdMethodArgs = this.addTypeArguments(methodArgs);

        // Only build an instantiation if at least one binding differs from its
        // formal type parameter; otherwise the "instance" is just the decl.
        final boolean classUseful  = isUsefulInstantiation(classParams,  tdClassArgs);
        final boolean methodUseful = isUsefulInstantiation(methodParams, tdMethodArgs);
        if (!classUseful && !methodUseful) return this.addDeclaration(m);

        // Look up (or create) the receiver BEFORE pushing the frame.
        // addObjectInst pushes its own frame that copies from prior;
        // if the frame was pushed first, its typeArgs would inherit the method-level
        // bindings and the ObjectInst would end up with the wrong instanceTypes.
        // When classUseful is false, still prefer the generic's ObjectDecl
        // as the receiver so the MethodInst points back at its class.
        final Ref<? extends TypeDesc> receiver = this.getReceiverForCall(receiverRef, declClass, classUseful);

        try {
            this.instantiator.pushFrame();
            for (int i = 0; i < classParams.size();  i++) this.instantiator.add(classParams.get(i),  tdClassArgs.get(i));
            for (int i = 0; i < methodParams.size(); i++) this.instantiator.add(methodParams.get(i), tdMethodArgs.get(i));

            return this.proj.methodInsts.create(this.log, this.instantiator,
                new ElementKey(m, this.instantiator.typeArgs()),
                "method for call site " + SpoonUtils.describeElem(m),
                () -> {
                    final Ref<ObjectDecl>               recvDecl      = this.addObjectDecl(declClass);
                    final Ref<MethodDecl>               generic       = this.addMethodDecl(recvDecl, m);
                    final List<Ref<? extends TypeDesc>> instanceTypes = this.instantiator.typeArgs();
                    final Ref<Signature>                resolved      = this.addSignature(m);
                    return new MethodInst(generic, receiver, instanceTypes, resolved);
                },
                (Ref<MethodInst> ref, MethodInst mi) -> {
                    if (receiver != null && receiver.getResolved() instanceof ObjectInst recvInst)
                        recvInst.methods.add(ref);
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    /**
     * addMethodInstForCall (constructor overload) creates or fetches a
     * MethodInst that captures both the constructed class's type args and
     * the constructor's own type args at the call site. Falls back to the
     * generic MethodDecl when narrowing is not useful or possible.
     */
    public Ref<? extends Construct> addMethodInstForCall(CtConstructorCall<?> cc) throws Exception {
        final CtExecutableReference<?> er = cc.getExecutable();
        if (er == null) return null;
        // Note: er.isImplicit() is usually true for ctor calls (you don't write
        // <init> in source), so do NOT gate on it here. The Analyzer's
        // addConstructorCallUsage already filters synthetic default ctors via
        // ctor.isImplicit() before delegating to us.

        final CtExecutable<?> ex = er.getDeclaration();
        if (!(ex instanceof CtConstructor<?> ctor)) return null;

        // Constructor's declaring class.
        final CtType<?> declType = ctor.getDeclaringType();
        if (!(declType instanceof CtClass<?> declClass))
            return this.addMethodDeclForConstructor(ctor);

        // Class-level actual type args come from the constructed type (Bar<S>).
        // Note: er.getDeclaringType() drops the args (Spoon caveat), so use cc.getType().
        final CtTypeReference<?> receiverRef = cc.getType();
        final ArrayList<CtTypeReference<?>> classArgs = new ArrayList<>();
        this.collectActualTypeArgs(receiverRef, classArgs);
        // Constructor's own type args are on the CtConstructorCall itself (`new <P>Bar<S>()`).
        final List<CtTypeReference<?>> ctorArgs = cc.getActualTypeArguments();

        final List<Ref<TypeParam>> classParams = this.addTypeParams(declClass);
        final List<Ref<TypeParam>> ctorParams  = new ArrayList<>();
        for (CtTypeParameter tp : ctor.getFormalCtTypeParameters()) {
            Ref<TypeParam> tpRef = this.addTypeParam(tp);
            Require.notNull(tpRef, "type parameter for " + SpoonUtils.describeElem(tp) +
                " in constructor " + SpoonUtils.describeElem(ctor) + " may not be null");
            ctorParams.add(tpRef);
        }
        if (classParams.isEmpty() && ctorParams.isEmpty()) return this.addMethodDeclForConstructor(ctor);

        if (!classParams.isEmpty() && classArgs.size() != classParams.size()) return this.addMethodDeclForConstructor(ctor);
        if (!ctorParams.isEmpty()  && ctorArgs.size()  != ctorParams.size())  return this.addMethodDeclForConstructor(ctor);

        final List<Ref<? extends TypeDesc>> tdClassArgs = this.addTypeArguments(classArgs);
        final List<Ref<? extends TypeDesc>> tdCtorArgs  = this.addTypeArguments(ctorArgs);

        final boolean classUseful = this.isUsefulInstantiation(classParams, tdClassArgs);
        final boolean ctorUseful  = this.isUsefulInstantiation(ctorParams,  tdCtorArgs);
        if (!classUseful && !ctorUseful) return this.addMethodDeclForConstructor(ctor);

        // Look up (or create) the receiver BEFORE pushing the frame,
        // otherwise addObjectInst's own frame would inherit the ctor bindings
        // and produce an ObjectInst with too many instanceTypes.
        // When classUseful is false, still prefer the generic's ObjectDecl
        // as the receiver so the MethodInst points back at its class.
        final Ref<? extends TypeDesc> receiver = this.getReceiverForCall(receiverRef, declClass, classUseful);

        try {
            this.instantiator.pushFrame();
            for (int i = 0; i < classParams.size(); i++) this.instantiator.add(classParams.get(i), tdClassArgs.get(i));
            for (int i = 0; i < ctorParams.size();  i++) this.instantiator.add(ctorParams.get(i),  tdCtorArgs.get(i));

            return this.proj.methodInsts.create(this.log, this.instantiator,
                new ElementKey(ctor, this.instantiator.typeArgs()),
                "constructor for call site " + SpoonUtils.describeElem(ctor),
                () -> {
                    final Ref<ObjectDecl>               recvDecl      = this.addObjectDecl(declClass);
                    final Ref<MethodDecl>               generic       = this.addMethodDeclForConstructor(recvDecl, ctor);
                    final List<Ref<? extends TypeDesc>> instanceTypes = this.instantiator.typeArgs();
                    final Ref<Signature>                resolved      = this.addSignatureForConstructor(ctor);
                    return new MethodInst(generic, receiver, instanceTypes, resolved);
                },
                (Ref<MethodInst> ref, MethodInst mi) -> {
                    if (receiver != null && receiver.getResolved() instanceof ObjectInst recvInst)
                        recvInst.methods.add(ref);
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    /**
     * Get the receiver type reference for a method invocation. Prefer the
     * target's declared type (which carries actual type args like `Foo<Bar>`);
     * fall back to a parameterized reference of the declaring class so the
     * class-level type-param count still lines up.
     */
    private CtTypeReference<?> getCallReceiverTypeRef(CtInvocation<?> in, CtType<?> declType) {
        final CtExpression<?> target = in.getTarget();
        if (target != null) {
            final CtTypeReference<?> tt = target.getType();
            if (tt != null) return tt;
        }
        return SpoonUtils.parameterizedRef(declType);
    }

    /**
     * Resolve the receiver ref for a call-site MethodInst. When the class is
     * meaningfully instantiated (preferInst==true), returns an ObjectInst ref;
     * otherwise falls back to the generic ObjectDecl so the MethodInst still
     * points at its declaring class. Returns null only when neither can be
     * resolved (e.g. addObjectDecl couldn't produce anything).
     */
    private Ref<? extends TypeDesc> getReceiverForCall(CtTypeReference<?> receiverRef, CtClass<?> declClass, boolean preferInst) throws Exception {
        if (preferInst) return this.addObjectInst(receiverRef, declClass);
        return this.addObjectDecl(declClass);
    }

    private boolean isUsefulInstantiation(List<Ref<TypeParam>> params, List<Ref<? extends TypeDesc>> args) throws Exception {
        if (params.size() != args.size()) return false;
        final CmpOptions options = new CmpOptions();
        options.useResolved = true;
        for (int i = 0; i < params.size(); i++) {
            if (Cmp.run(args.get(i).getCmp(params.get(i), options)) != 0) return true;
        }
        return false;
    }

    public Ref<Abstract> addAbstract(CtMethod<?> m) throws Exception {
        Require.notObjectMethod(m);
        final Ref<Abstract> ref = this.proj.abstracts.create(this.log, this.instantiator,
            new ElementKey(m, this.instantiator.typeArgs()),
            "abstract " + SpoonUtils.describeElem(m),
            () -> {
                final String         name      = m.getSimpleName();
                final Ref<Signature> signature = this.addSignature(m);
                return new Abstract(name, signature);
            });

        // If the declaring interface is a shadow, attach this abstract to
        // its abstracts so the stub grows only with used methods.
        final CtType<?> declType = m.getDeclaringType();
        if (declType instanceof CtInterface<?> declInter && declInter.isShadow()) {
            final Ref<InterfaceDesc> descRef = this.addInterfaceDesc(declInter);
            // TODO: Determine if it would be possible for the descRef to not
            // be resolved. If so, determine how to handle adding the abstracts.
            // Maybe Ref could have a set of OnResolved methods that are called
            // when the reference is resolved.
            if (descRef != null && descRef.isResolved())
                descRef.getResolved().abstracts.add(ref);
        }
        return ref;
    }

    public Ref<Signature> addSignature(CtMethod<?> m) throws Exception {
        Require.notObjectMethod(m);
        return this.proj.signatures.create(this.log, this.instantiator,
            new ElementKey(m, this.instantiator.typeArgs()),
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

    public Ref<MethodDecl> addMethodDeclForConstructor(CtConstructor<?> ctor) throws Exception {
        if (ctor.isImplicit()) return null;
        if (!ctor.isParentInitialized()) {
            this.log.warning("failed to constructor: parent not initialized for " + SpoonUtils.describeElem(ctor));
            return null;
        }
        if (ctor.getParent() instanceof CtClass c) {
            final Ref<ObjectDecl> receiver = this.addObjectDecl(c);
            return this.addMethodDeclForConstructor(receiver, ctor);
        }
        this.log.warning("failed to constructor: unknown parent " +
            SpoonUtils.describeElem(ctor.getParent()) + " for " + SpoonUtils.describeElem(ctor));
        return null;
    }

    public Ref<MethodDecl> addMethodDeclForConstructor(Ref<ObjectDecl> receiver, CtConstructor<?> ctor) throws Exception {
        try {
            // All declarations must be added without type arguments.
            this.instantiator.pushCleanFrame();
            return this.proj.methodDecls.create(log, this.instantiator,
                new ElementKey(ctor),
                "constructor " + SpoonUtils.describeElem(ctor),
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
                    //md.setNest(this.getParent(ctor)); // Not needed because of receiver
                    final ObjectDecl recv = receiver.mustGetResolved();
                    recv.methodDecls.add(ref);
                    this.pendingMetrics.add(ctor);
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    public Ref<Signature> addSignatureForConstructor(CtConstructor<?> m) throws Exception {
        return this.proj.signatures.create(this.log, this.instantiator,
            new ElementKey(m, this.instantiator.typeArgs()),
            "constructor signature " + SpoonUtils.describeElem(m),
            () -> {
                final List<CtParameter<?>> ps = m.getParameters();
                final boolean variadic = ps.size() > 0 && ps.get(ps.size()-1).isVarArgs();

                final ArrayList<Ref<Argument>> params = new ArrayList<>();
                for (CtParameter<?> p : ps) params.add(this.addArgument(p));
                
                final ArrayList<Ref<Argument>> results = new ArrayList<>();
                results.add(this.addArgument(SpoonUtils.parameterizedRef(m.getDeclaringType())));

                return new Signature(variadic, params, results);
            });
    }

    public Ref<Argument> addArgument(CtParameter<?> p) throws Exception {
        return this.proj.arguments.create(this.log, this.instantiator,
            new ElementKey(p, this.instantiator.typeArgs()),
            "parameter " + SpoonUtils.describeElem(p),
            () -> {
                final String            name = p.getSimpleName();
                Ref<? extends TypeDesc> type = this.addTypeDesc(p.getType());
                if (type == null) {
                    this.log.notice("argument " + SpoonUtils.describeElem(p) + " had a null type. The type likely "+
                        "was an attribute, external dependency, or some other type not handled by the abstractor so using anyDesc.");
                    type = this.proj.baker.anyDesc();
                }
                return new Argument(name, type);
            });
    }
    
    public Ref<Argument> addArgument(CtTypeReference<?> p) throws Exception {
        return this.proj.arguments.create(this.log, this.instantiator,
            new ElementKey(p, this.instantiator.typeArgs()),
            "parameter <unnamed> " + SpoonUtils.describeGeneric(p),
            () -> {
                Ref<? extends TypeDesc> type = this.addTypeDesc(p);
                if (type == null) {
                    this.log.notice("argument " + SpoonUtils.describeElem(p) + " had a null type. The type likely "+
                        "was an attribute, external dependency, or some other type not handled by the abstractor so using anyDesc.");
                    type = this.proj.baker.anyDesc();
                }
                return new Argument("", type);
            });
    }
    
    public Ref<StructDesc> addStructDesc(CtType<?> c) throws Exception {
        return this.proj.structDescs.create(this.log, this.instantiator,
            new ElementKey(c, this.instantiator.typeArgs()),
            "struct " + SpoonUtils.describeElem(c),
            (Ref<StructDesc> ref) -> {
                final ArrayList<Ref<Field>> fields = new ArrayList<>();

                // For shadow types, fields (including $super/$nest) are attached on demand
                // by addSelection so only track what's actually referenced.
                if (c.isShadow()) return new StructDesc(fields, ref);

                // Collect all fields.
                for (CtFieldReference<?> fr : c.getAllFields())
                    fields.add(this.addField(fr));

                // Add extended class as a "$super" field.
                final CtTypeReference<?> superFr = c.getSuperclass();
                if (superFr != null) fields.add(this.addField("$super", superFr));

                // Add access to nesting class as a "$nest" field.
                if (SpoonUtils.isNested(c)) {
                    if (c.getParent() instanceof CtTypeReference<?> nest) {
                        fields.add(this.addField("$nest", nest));
                    } else if (c.getParent() instanceof CtType<?> nest) {
                        fields.add(this.addField("$nest", nest.getReference()));
                    } else {
                        this.log.warning("Unhandled nested object decl " + SpoonUtils.describeElem(c) +
                            " in " + SpoonUtils.describeElem(c.getParent()));
                    }
                }

                return new StructDesc(fields);
            });
    }

    public Ref<Field> addField(CtFieldReference<?> f) throws Exception {
        return this.addField(f.getFieldDeclaration());
    }

    public Ref<Field> addField(CtField<?> f) throws Exception {
        return this.proj.fields.create(this.log, this.instantiator,
            new ElementKey(f, this.instantiator.typeArgs()),
            "field " + SpoonUtils.describeElem(f),
            () -> {
                final String            name = f.getSimpleName();
                Ref<? extends TypeDesc> type = this.addTypeDesc(f.getType());
                if (type == null) {
                    this.log.notice("field " + SpoonUtils.describeElem(f) + " had a null type. The type likely "+
                        "was an attribute, external dependency, or some other type not handled by the abstractor so using anyDesc.");
                    type = this.proj.baker.anyDesc();
                }
                return new Field(name, type);
            },
            (Ref<Field> ref, Field field) -> {
                field.setVisibility(f);
            });
    }

    public Ref<Field> addField(String name, CtTypeReference<?> f) throws Exception {
        return this.proj.fields.create(this.log, this.instantiator,
            new ElementKey(f, this.instantiator.typeArgs()),
            "field " + name,
            () -> {
                final Ref<? extends TypeDesc> type = this.addTypeDesc(f);
                return new Field(name, type);
            });
    }

    public Ref<Selection> addSelection(CtFieldReference<?> ref) throws Exception {
        if (ref == null) return null;
        final CtField<?> field = ref.getFieldDeclaration();
        if (field == null) {
            this.log.notice("Skipping selection with no field declaration: " + SpoonUtils.describeElem(ref));
            return null;
        }

        // Resolve the origin from the reference's declaring type first so
        // capture the actual instantiation at the call site (e.g. Foo<Integer>
        // yields an ObjectInst, plain Foo yields an ObjectDecl). Fall back to
        // the field's declaring type if the reference doesn't carry one. If
        // neither is trackable the Selection would have a null origin (which
        // is not useful and rejected by the validator) so skip the Selection;
        // Analyzer callers (addRead / addWrite) are null-safe.
        final CtTypeReference<?> receiverRef = ref.getDeclaringType();
        Ref<? extends Construct> originResolved = null;
        if (receiverRef != null) originResolved = this.addTypeDesc(receiverRef);
        if (originResolved == null) originResolved = this.addDeclaration(field.getDeclaringType());
        if (originResolved == null) {
            this.log.notice("Skipping selection with no resolvable origin: " + SpoonUtils.describeElem(field));
            return null;
        }
        final Ref<? extends Construct> origin = originResolved;

        // Key on the field plus the receiver's actual type args so different
        // instantiations of the same field (Foo<Integer>.x vs Foo<String>.x)
        // produce distinct Selections.
        final ArrayList<CtTypeReference<?>> receiverArgs = new ArrayList<>();
        this.collectActualTypeArgs(receiverRef, receiverArgs);
        final List<Ref<? extends TypeDesc>> keyArgs = this.addTypeArguments(receiverArgs);

        final Ref<Selection> selRef = this.proj.selections.create(this.log, this.instantiator,
            new ElementKey(field, keyArgs),
            "select field " + SpoonUtils.describeElem(field),
            () -> new Selection(field.getSimpleName(), origin));

        // On-demand: if the origin is a shadow object (decl or inst), attach the
        // field to its StructDesc so the stub grows only with used fields.
        // TODO: Two empty shadow StructDescs currently dedup to the same construct
        // in the factory (getExisting), so attaches can leak between unrelated
        // shadow objects. Needs a per-owner discriminator on StructDesc (e.g. a
        // pin field) so empty stubs stay distinct.
        final Construct originCon = origin.isResolved() ? origin.getResolved() : null;
        if (originCon instanceof ObjectDecl od && od.isShadow) {
            this.attachShadowField(od.struct, field);
        } else if (originCon instanceof ObjectInst oi &&
                   oi.generic != null && oi.generic.isResolved() &&
                   oi.generic.getResolved().isShadow) {
            this.attachShadowField(oi.resData, field);
        }
        return selRef;
    }

    private void attachShadowField(Ref<StructDesc> structRef, CtField<?> field) throws Exception {
        if (structRef == null || !structRef.isResolved()) return;
        final StructDesc sd = structRef.getResolved();
        if (sd == null) return;
        final Ref<Field> fRef = this.addField(field);
        if (fRef != null && !sd.fields.contains(fRef)) sd.fields.add(fRef);
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
        return this.proj.basics.create(this.log, this.instantiator,
            new ElementKey(tr),
            "basic " + SpoonUtils.describeElem(tr),
            () -> {
                if (SpoonUtils.isVoid(tr))
                    throw new AbstractorException("A void was added as a basic");
                return new Basic(tr.getSimpleName());
            });
    }

    public ArrayList<Ref<? extends TypeDesc>> addTypeArguments(List<CtTypeReference<?>> trs) throws Exception {
        final ArrayList<Ref<? extends TypeDesc>> result = new ArrayList<>(trs.size());
        for (CtTypeReference<?> tr : trs) {
            Ref<? extends TypeDesc> trRef = this.addTypeDesc(tr);
            if (trRef == null) trRef = this.proj.baker.anyDesc();
            result.add(trRef);
        }
        return result;
    }

    private List<Ref<TypeParam>> addTypeParams(CtElement elem) throws Exception {
        final List<Ref<TypeParam>> result =
            (elem.getParent() instanceof CtType<?> parent && parent != null)
            ? this.addTypeParams(parent)
            : new ArrayList<>();

        if (elem instanceof CtFormalTypeDeclarer td) {
            for (CtTypeParameter tp : td.getFormalCtTypeParameters()) {
                Ref<TypeParam> tr = this.addTypeParam(tp);
                Require.notNull(tr, "type parameter may not be null for " +
                    SpoonUtils.describeElem(tp) + " in " + SpoonUtils.describeElem(elem));
                result.remove(tr); // remove any prior one
                result.add(tr);
            }
        }
        return result;
    }

    public Ref<TypeParam> addTypeParam(CtTypeParameter tp) throws Exception {
        // Do not use type arguments in the ElementKey for typeParams.
        // The typeParams will be replaced by the instantiator later.
        return this.proj.typeParams.create(this.log, this.instantiator,
            new ElementKey(tp, null),
            "type params " + SpoonUtils.describeElem(tp),
            () -> {
                final String                  name = tp.getSimpleName();
                // TODO: This does not seem to handle several bounds like `T extends A & B` (returns just `A`).
                final CtTypeReference<?>      tr   = tp.getTypeErasure();
                final Ref<? extends TypeDesc> type = this.addTypeDesc(tr);
                return new TypeParam(name, type);
            });
    }

    public Ref<TypeParam> addTypeParam(CtTypeParameterReference tpr) throws Exception {
        // Do not use type arguments in the ElementKey for typeParams.
        // The typeParams will be replaced by the instantiator later.
        final CtTypeParameter tp = tpr.getDeclaration();
        return this.proj.typeParams.create(this.log, this.instantiator,
            new ElementKey(tp != null ? tp : tpr, null),
            "type params reference " + SpoonUtils.describeElem(tpr),
            () -> {
                final String                  name = tpr.getSimpleName();
                final CtTypeReference<?>      tr   = tpr.getBoundingType();
                final Ref<? extends TypeDesc> type = this.addTypeDesc(tr);
                return new TypeParam(name, type);
            });
    }
    
    public Ref<Metrics> addMetrics(CtExecutable<?> m) throws Exception {
        return this.proj.metrics.create(this.log, this.instantiator,
            new ElementKey(m, this.instantiator.typeArgs()),
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
        if (c.isAnonymous()) {
            this.log.notice("Ignoring anonymous object declaration: " + SpoonUtils.describeElem(c));
            return null;
        }
        if (c.isLocalType()) {
            this.log.notice("Ignoring local object declaration: " + SpoonUtils.describeElem(c));
            return null;
        }
        try {
            // All declarations must be added without type arguments.
            this.instantiator.pushCleanFrame();
            return this.proj.objectDecls.create(this.log, this.instantiator,
                new ElementKey(c),
                "object decl " + SpoonUtils.describeElem(c),
                () -> {
                    final Ref<PackageCon>      pkg        = this.addPackageFor(c);
                    final Location             loc        = this.proj.locations.create(c.getPosition());
                    final String               name       = c.getSimpleName();
                    final Ref<StructDesc>      struct     = this.addStructDesc(c);
                    final List<Ref<TypeParam>> typeParams = this.addTypeParams(c);
                    Require.isIdentifier(name, "object decl name (" + name + ") was not an identifier: " + SpoonUtils.describeElem(c));
                    return new ObjectDecl(pkg, loc, name, struct, typeParams);
                },
                (Ref<ObjectDecl> ref, ObjectDecl obj) -> {
                    obj.setVisibility(c);
                    obj.setNest(this.getParent(c));

                    if (c.isShadow()) {
                        // Shadow object, so don't add fields, nested types, or methods
                        // proactively, they will be added as needed.
                        obj.isShadow = true;
                    } else {
                        for (CtType<?> nt : c.getNestedTypes())
                            obj.nestedTypes.add(this.addTypeDesc(nt.getReference()));
                        
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
                    }

                    // Synthesize an interface for this enum object.
                    // If a shadow, this will add a stub of the synthesized interface that can be added to.
                    obj.inter = this.synthesizeObjectInterface(c, ref);
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    private Ref<InterfaceDesc> synthesizeObjectInterface(CtClass<?> c, Ref<? extends Construct> pin) throws Exception {
        Require.require(pin != null || !c.isShadow(), "must have a pin for shadows: " + SpoonUtils.describeElem(c));

        // Synthesize the interface abstractions for the class.
        final TreeSet<Ref<Abstract>> abstracts = new TreeSet<Ref<Abstract>>();
        if (!c.isShadow()) {
            for (CtMethod<?> m : c.getAllMethods()) {
                if (!m.isStatic() && !SpoonUtils.isObjectMethod(m))
                    abstracts.add(this.addAbstract(m));
            }
        }

        if (abstracts.size() <= 0 && c.getSuperInterfaces().size() <= 0 && pin == null)
            return this.proj.baker.anyDesc();

        // Synthesize the interface description for the class.
        final InterfaceDesc it = new InterfaceDesc(abstracts, pin);
        final List<Ref<? extends TypeDesc>> typeArgs = this.instantiator.typeArgs();
        final Ref<InterfaceDesc> inter = this.proj.interfaceDescs.addOrGetRef(it, typeArgs, "interface for object");

        // Add direct super-interfaces this object extends.
        for (CtTypeReference<?> supRef : c.getSuperInterfaces()) {
            final CtType<?> supDecl = supRef.getTypeDeclaration(); // may be null for shadow/unresolved
            if (supDecl == null) {
                this.log.warning("Unhandled null super-interface for " + pin);
                continue;
            }
            if (!(supDecl instanceof CtInterface<?> supId)) {
                this.log.error("Unhandled super-interface " + SpoonUtils.describeElem(supDecl) + " for " + pin);
                continue;
            }
            it.inherits.add(this.addInterfaceDesc(supId));
        }
        return inter;
    }

    public Ref<? extends TypeDesc> addObjectInst(CtTypeReference<?> tr, CtClass<?> c) throws Exception {
        final Ref<ObjectDecl> decl = this.addObjectDecl(c);
        if (!SpoonUtils.isGenerics(c)) return decl;

        final List<Ref<TypeParam>> typeParams = this.addTypeParams(c);
        final List<Ref<? extends TypeDesc>> typeArgs = this.addTypeArguments(tr, typeParams);
        if (typeArgs == null) return decl;

        try {
            final boolean definedInNest = SpoonUtils.inSameNested(tr, c);
            if (definedInNest) this.instantiator.pushFrame();
            else this.instantiator.pushCleanFrame();
            for (int i = 0; i < typeParams.size(); i++)
                this.instantiator.add(typeParams.get(i), typeArgs.get(i));

            return this.proj.objectInsts.create(this.log, this.instantiator,
                new ElementKey(tr, this.instantiator.typeArgs()),
                "object instantiation "+SpoonUtils.describeGeneric(tr),
                (Ref<ObjectInst> ref) -> {
                    final Ref<StructDesc>    resData      = this.addStructDesc(c);
                    final Ref<InterfaceDesc> resInterface = this.synthesizeObjectInterface(c, ref);
                    return new ObjectInst(decl, this.instantiator.typeArgs(), resData, resInterface);
                },
                (Ref<ObjectInst> ref, ObjectInst obj) -> {
                    if (c.isShadow()) {
                        // Shadow object instantiation: methods and nested types are
                        // pulled in on demand via addMethodInstForCall / addTypeDesc.
                        return;
                    }

                    // Add constructors as (static) methods for the class instantiation.
                    for (CtConstructor<?> ctor : c.getConstructors()) {
                        if (ctor.getParent().equals(c)) {
                            if (ctor.isImplicit()) {
                                this.log.notice("skipping default constructor: " + ctor.getSignature());
                                continue;
                            }
                            this.addMethodInstForObjectInst(ref, ctor);
                        }
                    }

                    // Add methods for the class instantiation.
                    for (CtMethod<?> m : c.getAllMethods()) {
                        if (m.getParent().equals(c) && !SpoonUtils.isObjectMethod(m))
                            this.addMethodInstForObjectInst(ref, m);
                    }

                    // Create instances for all nested types too.
                    for (CtType<?> nt : c.getNestedTypes())
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
    private List<Ref<? extends TypeDesc>> addTypeArguments(CtTypeReference<?> tr, List<Ref<TypeParam>> typeParams) throws Exception {
        final ArrayList<CtTypeReference<?>> ctTypeArgs = new ArrayList<>();
        this.collectActualTypeArgs(tr, ctTypeArgs);

        final int count = ctTypeArgs.size();
        if (count <= 0) return null;
        if (count != typeParams.size()) return null;

        final ArrayList<Ref<? extends TypeDesc>> typeArgs = new ArrayList<>();
        for (CtTypeReference<?> ctTypeArg : ctTypeArgs) {
            Ref<? extends TypeDesc> tpRef = this.addTypeDesc(ctTypeArg);
            if (tpRef == null) tpRef = this.proj.baker.anyDesc();
            typeArgs.add(tpRef);
        }
        
        final CmpOptions options = new CmpOptions();
        options.useResolved = true;
        for (int i = 0; i < count; i++) {
            final Ref<? extends TypeDesc> ta = typeArgs.get(i);
            final Ref<? extends TypeDesc> tp = typeParams.get(i);
            final boolean isNotEq = Cmp.run(ta.getCmp(tp, options)) != 0;
            if (isNotEq) {
                // There was a difference so there is an instantiation
                return typeArgs;
            }
        }
        // There was no difference so the instantiation is not useful.
        return null;
    }

    /**
     * Collects the actual type arguments from a type reference's declaring chain,
     * outer-most first then the immediate arguments. This mirrors how
     * {@link #addTypeParams(CtElement)} walks the enclosing-type chain so that
     * nested generics (e.g. {@code Foo<T1>.Bar<Integer>}) line up by index.
     */
    private void collectActualTypeArgs(CtTypeReference<?> tr, List<CtTypeReference<?>> out) {
        if (tr == null) return;
        this.collectActualTypeArgs(tr.getDeclaringType(), out);
        final List<CtTypeReference<?>> args = tr.getActualTypeArguments();
        if (args != null) out.addAll(args);
    }

    public Ref<? extends TypeDesc> addTypeDesc(CtTypeReference<?> tr) throws Exception {
        if (tr == null) return null;

        // By default skip anonymous and local types since they can not escape the enclosing method,
        // e.g. `testData.java.test1004.Foo$1` is anonymous with `1` as the name.
        // They still will contribute to metrics via super-interfaces and extends.
        if (tr.isAnonymous()) {
            this.log.notice("Ignoring anonymous type: " + SpoonUtils.describeElem(tr));
            return null;
        }
        if (tr.isLocalType()) {
            this.log.notice("Ignoring local type: " + SpoonUtils.describeElem(tr));
            return null;
        }

        // Handle primitive types (e.g. `int` but not `String` nor `Integer`).
        if (tr.isPrimitive()) return this.addBasic(tr);
        
        // Handle an array (e.g. `T[]` not `List<T>`) type.
        if (tr.isArray()) return this.addArray((CtArrayTypeReference<?>)tr);

        // Handle wildcard types (e.g., `?`, `? extends Foo`, `? super Bar`).
        if (tr instanceof CtWildcardReference wr) return this.addWildcard(wr);

        // Type of the `null` literal in Spoon and not a real external type.
        if (SpoonUtils.isNull(tr)) return this.proj.baker.anyDesc();

        // A boxed type (e.g. Integer, String) that can alias as a basic.
        final Ref<Basic> boxed = this.proj.baker.basicForBoxedOrString(tr);
        if (boxed != null) return boxed;
        
        // If the type is an Object, return an any for the Object.
        if (SpoonUtils.isObject(tr)) return this.proj.baker.anyDesc();

        // If the type is a type parameter (reference), return a new type parameter.
        if (tr instanceof CtTypeParameterReference tpr) return this.instantiator.replace(this.addTypeParam(tpr));
        if (tr instanceof CtTypeParameter           tp) return this.instantiator.replace(this.addTypeParam(tp));

        // Get the actual type declaration that is being referenced.
        final CtType<?> ty = tr.getTypeDeclaration();
        if (ty == null) {
            this.log.warning("Type description did not have a declaration but "+
                "was not labelled as anonymous: " + SpoonUtils.describeElem(tr)+"\n"+
                "The type likely was an external dependency or some other type not handled by the abstractor so using anyDesc.");
            return this.proj.baker.anyDesc();
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
        if (e.isAnonymous()) {
            this.log.notice("Ignoring anonymous enumerator: " + SpoonUtils.describeElem(e));
            return null;
        }
        if (e.isLocalType()) {
            this.log.notice("Ignoring local enumerator: " + SpoonUtils.describeElem(e));
            return null;
        }
        try {
            // All declarations must be added without type arguments.
            this.instantiator.pushCleanFrame();
            return this.proj.objectDecls.create(this.log, this.instantiator,
                new ElementKey(e),
                "enum " + SpoonUtils.describeElem(e),
                () -> {
                    final String             name   = e.getSimpleName();
                    final Ref<PackageCon>    pkg    = this.addPackageFor(e);
                    final Location           loc    = this.proj.locations.create(e.getPosition());
                    final CtTypeReference<?> tr     = e.getSuperclass();
                    final Ref<StructDesc>    struct = this.proj.structDescs.create(this.log, this.instantiator,
                        new ElementKey(tr),
                        "enum struct " + SpoonUtils.describeElem(tr),
                        () -> {
                            final ArrayList<Ref<Field>> fields = new ArrayList<>();
                            fields.add(this.addField("$value", tr));
                            return new StructDesc(fields);
                        });

                    return new ObjectDecl(pkg, loc, name, struct, null);
                },
                (Ref<ObjectDecl> ref, ObjectDecl od) -> {
                    od.setVisibility(e);
                    od.setNest(this.getParent(e));

                    if (e.isShadow()) {
                        // Shadow enum, so don't add fields, nested types, and methods
                        // proactively, they will be added as needed.
                        od.isShadow = true;
                    } else {
                        for (CtType<?> nt : e.getNestedTypes())
                            od.nestedTypes.add(this.addTypeDesc(nt.getReference()));

                        // Finish by adding the "const values" to the package for each enumerator value.
                        for (CtEnumValue<?> ev: e.getEnumValues()) {
                            this.proj.values.create(this.log, this.instantiator,
                                new ElementKey(ev),
                                "enum value "+ SpoonUtils.describeElem(ev),
                                () -> {
                                    final String   name = ev.getSimpleName();
                                    final Location loc  = this.proj.locations.create(ev.getPosition());
                                    return new Value(od.pkg, loc, name, true, null, ref);
                                });
                        }

                        // Add constructors as (static) methods.
                        for (CtConstructor<?> ctor : e.getConstructors()) {
                            if (ctor.getParent().equals(e)) {
                                // Skip default constructors
                                if (ctor.isImplicit()) {
                                    this.log.notice("skipping default constructor: " + ctor.getSignature());
                                    continue;
                                }
                                this.addMethodDeclForConstructor(ref, ctor);
                            }
                        }

                        // Add methods for the enum.
                        for (CtMethod<?> m : e.getAllMethods()) {
                            if (!m.isImplicit() && m.getParent().equals(e) && !SpoonUtils.isObjectMethod(m))
                                this.addMethodDecl(ref, m);
                        }
                    }

                    // Synthesize an interface for this enum object.
                    // If a shadow, this will add a stub of the synthesized interface that can be added to.
                    od.inter = this.synthesizeObjectInterface(e, ref);
                });
        } finally {
            this.instantiator.popFrame();
        }
    }

    //===[ Processors ]=========================================================

    /**
     * performAbstraction will process all the packages and metrics, and
     * resolve anything else that needs to be done to finish the abstraction.
     */
    public void performAbstraction() throws Exception {
        this.log.measure("process pending packages",     () -> this.processPendingPackages());
        this.log.measure("short validation",             () -> this.shortValidate());
        this.log.measure("consolidate constructs",       () -> this.consolidateCons());
        this.log.measure("connect nests",                () -> this.connectNests());
        this.log.measure("collect package declarations", () -> this.collectPackageDeclarations());
        this.log.measure("addImports from usage",        () -> this.addImportsFromUsage());
        this.log.measure("remove empty packages",        () -> this.removeEmptyPackages());
    }

    private void processPendingPackages() throws Exception {
        while (!this.pendingPackages.isEmpty()) {
            final CtPackage pkg = this.pendingPackages.iterator().next();
            this.pendingPackages.remove(pkg);
            this.processPackage(pkg);
            this.processPendingMetrics();
            this.processDeferredFinishes();
        }
    }

    public Ref<PackageCon> processPackage(CtPackage pkg) throws Exception {
        return this.proj.packages.create(this.log, this.instantiator,
            new ElementKey(pkg),
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
                    this.proj.metrics.remove(this.log, met);
                }
            }
        }
    }

    private void processDeferredFinishes() throws Exception {
        boolean hasMore = true;
        while (hasMore) {
            hasMore = false;
            for (Factory<?> f : this.proj.factories)
                hasMore = f.runDeferredFinishes(this.log, this.instantiator) | hasMore;
        }
    }

    private void consolidateCons() throws Exception {
        new Consolidator(this.log, this.proj).consolidate();
    }

    private void connectNests() throws Exception {
        connectNests(this.proj.objectDecls);
        connectNests(this.proj.interfaceDecls);
        connectNests(this.proj.methodDecls);
        connectNests(this.proj.values);
    }

    private <T extends Declaration> void connectNests(Factory<T> factory) throws Exception {
        for (T decl : factory.getConSet())
            connectNests(factory, decl);
    }

    private <T extends Declaration> void connectNests(Factory<T> factory, T decl) throws Exception {
        final Ref<? extends Construct> nest = decl.getNest();
        if (nest != null && nest.isResolved() && nest.getResolved() instanceof Declaration p && decl instanceof TypeDesc) {
            final Ref<T> declRef = factory.addOrGetRef(decl, null, "connecting nests");
            @SuppressWarnings("unchecked")
            final Ref<? extends TypeDesc> tdRef = (Ref<? extends TypeDesc>)declRef;
            p.getNestedTypes().add(tdRef);
        }

        for (Ref<? extends TypeDesc> ref : decl.getNestedTypes()) {
            final TypeDesc res = ref.getResolved();
            if (res != null && res instanceof Declaration other && other.getNest() == null)
                other.setNest(ref);
        }
    }

    /**
     * This adds all the declarations into the lists in the packages for the 
     * ype of declaration in the package.
     */
    private void collectPackageDeclarations() throws Exception {
        for (MethodDecl m : this.proj.methodDecls.getConSet()) {
            final PackageCon pkg = m.pkg.mustGetResolved();
            if (pkg == null) this.log.error("package for method is null: " + m);
            pkg.methodDecls.add(this.proj.methodDecls.addOrGetRef(m, null, "method in package " + pkg));
        }

        for (ObjectDecl obj : this.proj.objectDecls.getConSet()) {
            final PackageCon pkg = obj.pkg.mustGetResolved();
            if (pkg == null) this.log.error("package for object is null: " + obj);
            pkg.objectDecls.add(this.proj.objectDecls.addOrGetRef(obj, null, "object in package " + pkg));
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

    /**
     * Populate each PackageCon.imports with references to every OTHER package
     * whose declarations are transitively referenced by this package's own
     * decls. Instantiations move up to their generic decl to find the package;
     * intermediate descriptions (StructDesc, InterfaceDesc's pin, Signature, ...)
     * are walked through so their embedded refs contribute too. Runs after
     * consolidation, so every ref is expected to be resolved (mustGetResolved).
     */
    private void addImportsFromUsage() throws Exception {
        for (Ref<PackageCon> pkgRef : this.proj.packages.getRefSet())
            this.addImportsFromUsageFor(pkgRef);
    }

    private void addImportsFromUsageFor(Ref<PackageCon> pkgRef) throws Exception {
        final HashSet<Construct> visited = new HashSet<>();
        final Bundle<Ref<? extends Construct>> bundle = new Bundle<>();
        final PackageCon pkg = pkgRef.mustGetResolved();
        bundle.add(pkg.subConstructs());
        while (bundle.hasNext()) {
            final Ref<? extends Construct> ref = bundle.next();
            if (ref == null) return;
            
            final Construct c = ref.mustGetResolved();
            if (!visited.add(c)) return;

            if (c instanceof Declaration decl) {
                final Ref<PackageCon> otherPkgRef = decl.pkgRef();
                if (otherPkgRef != null && otherPkgRef != pkgRef) {
                    pkg.imports.add(otherPkgRef);
                }
            }

            bundle.add(c.subConstructs());
        }
    }

    private void removeEmptyPackages() {
        this.proj.packages.removeIf(log, (PackageCon pc) -> pc.isEmpty());
        this.proj.packages.setIndices();
    }

    public void shortValidate() throws Exception {
        this.log.measure("performShortValidation", () -> this.performShortValidation());
    }

    private void performShortValidation() throws Exception {
        final boolean hadErrors = this.log.errorCount() > 0;
        new Validator(this.log, this.proj).shortValidate();
        if (this.log.errorCount() > 0) {
            if (hadErrors)
                throw new AbstractorException("Errors logged before short validation.");
            throw new AbstractorException("Errors logged during short validation.");
        }
    }

    public void validate() throws Exception {
        this.log.measure("performValidation", () -> this.performValidation());
    }

    private void performValidation() throws Exception {
        final boolean hadErrors = this.log.errorCount() > 0;
        new Validator(this.log, this.proj).validate();
        if (this.log.errorCount() > 0) {
            if (hadErrors)
                throw new AbstractorException("Errors logged before validation.");

            final boolean showAbstract = false;
            if (showAbstract) {
                JsonHelper h = new JsonHelper();
                this.log.notice("\n" + JsonFormat.Relaxed().format(this.proj.toJson(h)));
            }
            throw new AbstractorException("Errors logged during validation.");
        }
    }
}
