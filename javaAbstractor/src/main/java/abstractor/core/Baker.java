package abstractor.core;

import java.util.*;

import spoon.reflect.reference.CtTypeReference;
import abstractor.core.constructs.Abstract;
import abstractor.core.constructs.Argument;
import abstractor.core.constructs.Basic;
import abstractor.core.constructs.Construct;
import abstractor.core.constructs.Factory;
import abstractor.core.constructs.InterfaceDecl;
import abstractor.core.constructs.InterfaceDesc;
import abstractor.core.constructs.InterfaceInst;
import abstractor.core.constructs.PackageCon;
import abstractor.core.constructs.Project;
import abstractor.core.constructs.Ref;
import abstractor.core.constructs.Signature;
import abstractor.core.constructs.TypeDesc;
import abstractor.core.constructs.TypeParam;
import abstractor.core.require.Require;
import abstractor.core.spoonUtils.SpoonUtils;

public class Baker {
    private final Project proj;
    private final Map<String, Construct> cache;

    public Baker(Project proj) {
        this.proj  = proj;
        this.cache = new TreeMap<>();
    }

    public interface ConstructCreator<T extends Construct> { T create() throws Exception; }

    @SuppressWarnings("unchecked")
    private <T extends Construct> T getConstruct(String name, ConstructCreator<T> creator) throws Exception {
        final Construct existing = cache.get(name);
        if (existing != null) return (T)existing;

        final T value = creator.create();
        cache.put(name, value);
        return value;
    }

    private <T extends Construct> Ref<T> getConstruct(String name, Factory<T> factory, List<Construct> typeArgs, ConstructCreator<T> creator) throws Exception {
        return this.getConstruct(name, () -> {
            T con = creator.create();
            final T other = factory.getExisting(con);
            if (other != null) con = other;
            return factory.addOrGetRef(con, typeArgs, "baker: " + name);
        });
    }

    public Ref<PackageCon> builtinPackage() throws Exception {
        return this.getConstruct("builtinPackage", this.proj.packages, null,
            () -> new PackageCon("$builtin", ""));
    }

    // Creates a new empty interface (any) that is the base interface type of all non-basic types.
    public Ref<InterfaceDesc> anyDesc() throws Exception {
        return this.getConstruct("objectInterfaceDesc", this.proj.interfaceDescs, null,
            () -> new InterfaceDesc(Collections.emptySortedSet()));
    }

    private Ref<Basic> intBasic() throws Exception {
        return this.getConstruct("intBasic", this.proj.basics, null,
            () -> new Basic("int"));
    }

    public Ref<TypeParam> genT() throws Exception {
        return this.getConstruct("genT", this.proj.typeParams, null,
            () -> new TypeParam("T", this.anyDesc()));
    }

    private Ref<Argument> intReturn() throws Exception {
        return this.getConstruct("intReturn", this.proj.arguments, null,
            () -> new Argument(this.intBasic()));
    }

    private Ref<Argument> intIndexParam() throws Exception {
        return this.getConstruct("intIndexParam", this.proj.arguments, null,
            () -> new Argument("index", this.intBasic()));
    }

    private Ref<Argument> genReturn(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("genReturn<" + tdName + ">", this.proj.arguments, Arrays.asList(td),
            () -> new Argument(td));
    }

    private Ref<Argument> valueGenParam(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("valueGenParam<" + tdName + ">", this.proj.arguments, Arrays.asList(td),
            () -> new Argument("value", td));
    }

    private Ref<Signature> lenSignature() throws Exception {
        return this.getConstruct("lenSignature", this.proj.signatures, null, () -> {
            final ArrayList<Ref<Argument>> results = new ArrayList<>();
            results.add(this.intReturn());
            return new Signature(false, Collections.emptyList(), results);
        });
    }

    private Ref<Signature> getSignature(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("getSignature<" + tdName + ">", this.proj.signatures, Arrays.asList(td), () -> {
            final ArrayList<Ref<Argument>> params = new ArrayList<>();
            params.add(this.intIndexParam());
            final ArrayList<Ref<Argument>> results = new ArrayList<>();
            results.add(this.genReturn(tdName, td));
            return new Signature(false, params, results);
        });
    }

    private Ref<Signature> setSignature(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("setSignature<" + tdName + ">", this.proj.signatures, Arrays.asList(td), () -> {
            final ArrayList<Ref<Argument>> params = new ArrayList<>();
            params.add(this.intIndexParam());
            params.add(this.valueGenParam(tdName, td));
            return new Signature(false, params, Collections.emptyList());
        });
    }

    private Ref<Abstract> lenAbstract() throws Exception {
        return this.getConstruct("lenAbstract", this.proj.abstracts, null,
            () -> new Abstract("$len", this.lenSignature()));
    }

    private Ref<Abstract> getAbstract(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("getAbstract<" + tdName + ">", this.proj.abstracts, Arrays.asList(td),
            () -> new Abstract("$get", this.getSignature(tdName, td)));
    }

    private Ref<Abstract> setAbstract(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("setAbstract<" + tdName + ">", this.proj.abstracts, Arrays.asList(td),
            () -> new Abstract("$set", this.setSignature(tdName, td)));
    }

    private Ref<InterfaceDesc> arrayDesc(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("arrayDesc<" + tdName + ">", this.proj.interfaceDescs, Arrays.asList(td), () -> {
            final TreeSet<Ref<Abstract>> abs = new TreeSet<>();
            abs.add(this.lenAbstract());
            abs.add(this.getAbstract(tdName, td));
            abs.add(this.setAbstract(tdName, td));
            return new InterfaceDesc(abs);
        });
    }

    public Ref<InterfaceDecl> arrayDecl() throws Exception {
        return this.getConstruct("arrayDecl", this.proj.interfaceDecls, null, () -> {
            final Ref<PackageCon> pkg = this.builtinPackage();
            final Ref<TypeParam>  tdT = this.genT();
            final ArrayList<Ref<TypeParam>> tp = new ArrayList<>();
            tp.add(tdT);
            final Ref<InterfaceDesc> desc = this.arrayDesc("T", tdT);
            return new InterfaceDecl(pkg, null, "$Array", desc, tp);
        });
    }
    
    public Ref<InterfaceInst> arrayInst(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        // Check that `td` is not `T` to prevent $Array<T> being instantiated with T.
        if (td.isResolved()) {
            final Ref<TypeParam> tdT = this.genT();
            Require.notEqual(td.getResolved(), tdT.getResolved());
        }

        return this.getConstruct("arrayInst<" + tdName + ">", this.proj.interfaceInsts, Arrays.asList(td), () -> {
            final Ref<InterfaceDecl> generic = this.arrayDecl();
            final ArrayList<Ref<? extends TypeDesc>> its = new ArrayList<>();
            its.add(td);
            final Ref<InterfaceDesc> resolved = this.arrayDesc(tdName, td);
            return new InterfaceInst(generic, its, resolved);
        });
    }

    /**
     * Qualified erasure name (e.g. java.lang.Integer) to basic type name (i.e. int).
     */
    private static final Map<String, String> boxedQualifiedNameToBasic = Map.of(
        "java.lang.Byte",      "byte",
        "java.lang.Short",     "short",
        "java.lang.Integer",   "int",
        "java.lang.Long",      "long",
        "java.lang.Float",     "float",
        "java.lang.Double",    "double",
        "java.lang.Character", "char",
        "java.lang.Boolean",   "boolean",
        "java.lang.String",    "string"
    );

    /**
     * If qualifiedErasureName is a boxed primitive (e.g. java.lang.Integer)
     * or java.lang.String, returns the shared Basic, otherwise null.
     */
    public Ref<Basic> basicForBoxedOrString(CtTypeReference<?> tr) throws Exception {
        final String name      = tr.getQualifiedName();
        final String basicName = boxedQualifiedNameToBasic.get(name);
        if (basicName == null) return null;
        return this.getConstruct("boxedBasic: " + SpoonUtils.describeElem(tr), this.proj.basics,  null,
            () -> new Basic(basicName));
    }
}
