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
    public void addMavenProject(String mavenProject) {
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
    public void addClassFromSource(String ...sourceLines) {
        String source = String.join("\n", sourceLines);
        this.addObjectDecl(Launcher.parseClass(source));
    }

    public interface ConstructCreator<T extends Construct> { T create(); }
    public interface FinishConstruct<T extends Construct> { void finish(T con); }
    public interface InProgressHandle<T extends Construct> { T handle(); }

    public <T extends Construct, U extends T> T create(Factory<U> factory, CtElement elem,
        String title, ConstructCreator<U> c, FinishConstruct<U> f, InProgressHandle<T> h) {
        final T existing = factory.get(elem);
        if (existing != null) return existing;

        try {
            if (log != null) {
                log.log("Adding " + title);
                log.push();
            }

            if (factory.inProgress(elem)) {
                if (h != null) return h.handle();
                log.error("Error: Already in progress: " + title);
                return null;
            }

            factory.startProgress(elem);
            final U newCon = c.create();
            factory.stopProgress(elem);

            final U other = factory.get(newCon);
            if (other != null) return other;
            factory.add(elem, newCon);
            
            this.resolveReferencesFor(elem, newCon);
            if (f != null) f.finish(newCon);
            return newCon;
        } finally {
            if (log != null) log.pop();
        }
    }
    
    public <T extends Construct, U extends T> T create(Factory<U> factory, CtElement elem,
        String title, ConstructCreator<U> c, FinishConstruct<U> f) {
        return this.create(factory, elem, title, c, f, null);
    }
    
    public <T extends Construct, U extends T> T create(Factory<U> factory, CtElement elem,
        String title, ConstructCreator<U> c) {
        return this.create(factory, elem, title, c, null, null);
    }

    public <T extends Construct> void resolveReferencesFor(CtElement elem, T con) {
        if (con instanceof Reference<?>) return;
        
        if (con instanceof Declaration decl) {
            for (DeclarationRef ref : this.proj.declRefs) {
                if (ref.elem.equals(elem)) ref.setResolved(decl);
            }
        }

        if (con instanceof TypeDesc td) {
            for (TypeDescRef ref : this.proj.typeDescRefs) {
                if (ref.elem.equals(elem)) ref.setResolved(td);
            }
        }

        if (con instanceof TypeParam tp) {
            for (TypeParamRef ref : this.proj.typeParamRefs) {
                if (ref.elem.equals(elem)) ref.setResolved(tp);
            }
        }
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
        return this.create(this.proj.packages, pkg,
            "package " + pkg.getQualifiedName(),
            () -> {
                final String name = pkg.getQualifiedName();
                final String path = packagePath(pkg);
                return new PackageCon(name, path);
            },
            (PackageCon pkgCon) -> {
                for (CtType<?> t : pkg.getTypes()) {
                    if (t instanceof CtClass<?> c) {
                        if (!this.proj.objectDecls.inProgress(c)) this.addObjectDecl(c);
                    } else if (t instanceof CtInterface<?> i) {
                        if (!this.proj.interfaceDecls.inProgress(i)) this.addInterfaceDecl(i);
                    } else this.log.error("Unhandled (" + t.getClass().getName() + ") "+t.getQualifiedName());
                }
            });
    }

    /**
     * Handles adding and processing classes, enums, and records.
     * @param c The class to process.
     */
    private TypeDesc addObjectDecl(CtClass<?> c) {
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
        
                // TODO: Finish implementing
            },
            () -> {
                return this.addTypeDescRef(c.getReference());
            });
    }

    private MethodDecl addMethod(ObjectDecl receiver, CtMethod<?> m) {
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

    private Signature addSignature(CtMethod<?> m) {
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

    private Argument addArgument(CtParameter<?> p) {
        return this.create(this.proj.arguments, p,
            "parameter " + p.getSimpleName(),
            () -> {
                final String   name = p.getSimpleName();
                final TypeDesc type = this.addTypeDesc(p.getType());
                return new Argument(name, type);
            });
    }
    
    private Argument addArgument(CtTypeReference<?> t) {
        return this.create(this.proj.arguments, t,
            "parameter <unnamed> " + t.getSimpleName(),
            () -> {
                final TypeDesc type = this.addTypeDesc(t);
                return new Argument("", type);
            });
    }

    private StructDesc addStruct(CtClass<?> c) {
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

    private Field addField(CtField<?> f) {
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
    
    private TypeDesc addInterfaceDecl(CtInterface<?> i) {
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

    private TypeDesc addTypeDesc(CtTypeReference<?> tr) {
        if (tr.isPrimitive()) return this.addBasic(tr);
        if (tr.isArray())     return this.addArray((CtArrayTypeReference<?>)tr);
        if (tr.isClass())     return this.addNamedTypeDesc(tr);
        if (tr.isInterface()) return this.addNamedTypeDesc(tr);
        if (tr.isGenerics())  return this.addTypeParam((CtTypeParameterReference)tr);

        // TODO: Finish implementing.
        return this.unknownTypeDesc(tr);
    }

    private TypeDesc unknownTypeDesc(CtTypeReference<?> tr) {
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

    private TypeDesc addNamedTypeDesc(CtTypeReference<?> tr) {



        // TODO: Switch based on availbility and some lock

        return this.addTypeDescRef(tr);
    }

    private TypeDescRef addTypeDescRef(CtTypeReference<?> tr) {
        return this.create(this.proj.typeDescRefs, tr,
            "type desc ref "+ tr.getSimpleName(),
            () -> {
                final String name = tr.getSimpleName();
                final String pkgPath = tr.getPackage().toString();
                final List<TypeDesc> tps = this.addTypeArguments(tr.getActualTypeArguments());
                return new TypeDescRef(tr, pkgPath, name, tps);
            });
    }

    private InterfaceInst addArray(CtArrayTypeReference<?> tr) {
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
    
    private Basic addBasic(CtTypeReference<?> tr) {
        return this.create(this.proj.basics, tr,
            "basic " + tr.getSimpleName(),
            () -> {
                final String name = tr.getSimpleName();
                if (name == "void")
                    this.log.error("A void was added as a basic.");
                return new Basic(name);
            });
    }

    private InterfaceDesc addInterfaceDesc(CtInterface<?> i) {
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

    private Abstract addAbstract(CtMethod<?> m) {
        return this.create(this.proj.abstracts, m,
            "abstract " + m.getSimpleName(),
            () -> {
                final String name = m.getSimpleName();
                final Signature signature = this.addSignature(m);
                return new Abstract(name, signature);
            });
    }

    private List<TypeDesc> addTypeArguments(List<CtTypeReference<?>> trs) {
        List<TypeDesc> result = new ArrayList<TypeDesc>(trs.size());
        for (CtTypeReference<?> tr : trs) result.add(this.addTypeDesc(tr));
        return result;
    }

    private List<TypeParam> addTypeParams(List<CtTypeParameter> tps) {
        List<TypeParam> result = new ArrayList<TypeParam>(tps.size());
        for (CtTypeParameter tp : tps) result.add(this.addTypeParam(tp));
        return result;
    }

    private TypeParam addTypeParam(CtTypeParameterReference tp) {
        return this.addTypeParam(tp.getDeclaration());
    }

    private TypeParam addTypeParam(CtTypeParameter tp) {
        return this.create(this.proj.typeParams, tp,
            "type params " + tp.getQualifiedName(),
            () -> {
                final String name = tp.getQualifiedName();
                
                // TODO: Remove
                //System.out.println(">> " + name + " >> " + tp.prettyprint());
                //System.out.println(">> >> " + tp.getTypeErasure());
                //for (CtTypeReference<?> tpr : tp.getSuperInterfaces())
                //    System.out.println(">>  >> " + tpr.getSimpleName() + " >> " + tpr.prettyprint());

                final TypeDesc type = this.addTypeDesc(tp.getTypeErasure());

                // TODO: Finish
                return new TypeParam(name, type);
            });
    }

    private Metrics addMetrics(CtMethod<?> m) {
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
