package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.Collections;
import java.util.Map;
import java.util.TreeMap;
import java.util.TreeSet;

import abstractor.core.spoonUtils.SpoonUtils;
import spoon.reflect.reference.CtTypeReference;

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

    private <T extends Construct> Ref<T> getConstruct(String name, Factory<T> factory, ConstructCreator<T> creator) throws Exception {
        return this.getConstruct(name, () -> factory.addOrGetRef(creator.create(), "baker: " + name));
    }

    // anyDesc creates a new empty interface (any) that is the base type of all non-basic types.
    public Ref<InterfaceDesc> anyDesc() throws Exception {
        return this.getConstruct("objectInterfaceDesc", this.proj.interfaceDescs,
            () -> new InterfaceDesc(Collections.emptySortedSet()));
    }
    
    private Ref<Basic> intBasic() throws Exception {
        return this.getConstruct("intBasic", this.proj.basics,
            () -> new Basic("int"));
    }

    private Ref<TypeParam> genT() throws Exception {
        return this.getConstruct("genT", this.proj.typeParams,
            () -> new TypeParam("T", this.anyDesc()));
    }

    private Ref<Argument> intReturn() throws Exception {
        return this.getConstruct("intReturn", this.proj.arguments,
            () -> new Argument(this.intBasic()));
    }

    private Ref<Argument> intIndexParam() throws Exception {
        return this.getConstruct("intIndexParam", this.proj.arguments,
            () -> new Argument("index", this.intBasic()));
    }

    private Ref<Argument> genReturn(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("genReturn<" + tdName + ">", this.proj.arguments,
            () -> new Argument(td));
    }

    private Ref<Argument> valueGenParam(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("valueGenParam<" + tdName + ">", this.proj.arguments,
            () -> new Argument("value", td));
    }

    private Ref<Signature> lenSignature() throws Exception {
        return this.getConstruct("lenSignature", this.proj.signatures, () -> {
            final ArrayList<Ref<Argument>> results = new ArrayList<Ref<Argument>>();
            results.add(this.intReturn());
            return new Signature(false, Collections.emptyList(), results);
        });
    }

    private Ref<Signature> getSignature(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("getSignature<" + tdName + ">", this.proj.signatures, () -> {
            final ArrayList<Ref<Argument>> params = new ArrayList<Ref<Argument>>();
            params.add(this.intIndexParam());
            final ArrayList<Ref<Argument>> results = new ArrayList<Ref<Argument>>();
            results.add(this.genReturn(tdName, td));
            return new Signature(false, params, results);
        });
    }

    private Ref<Signature> setSignature(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("setSignature<" + tdName + ">", this.proj.signatures, () -> {
            final ArrayList<Ref<Argument>> params = new ArrayList<Ref<Argument>>();
            params.add(this.intIndexParam());
            params.add(this.valueGenParam(tdName, td));
            return new Signature(false, params, Collections.emptyList());
        });
    }

    private Ref<Abstract> lenAbstract() throws Exception {
        return this.getConstruct("lenAbstract", this.proj.abstracts,
            () -> new Abstract("$len", this.lenSignature()));
    }

    private Ref<Abstract> getAbstract(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("getAbstract<" + tdName + ">", this.proj.abstracts,
            () -> new Abstract("$get", this.getSignature(tdName, td)));
    }

    private Ref<Abstract> setAbstract(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("setAbstract<" + tdName + ">", this.proj.abstracts,
            () -> new Abstract("$set", this.setSignature(tdName, td)));
    }

    private Ref<InterfaceDesc> arrayDesc(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("arrayDesc<" + tdName + ">", this.proj.interfaceDescs, () -> {
            final TreeSet<Ref<Abstract>> abs = new TreeSet<Ref<Abstract>>();
            abs.add(this.lenAbstract());
            abs.add(this.getAbstract(tdName, td));
            abs.add(this.setAbstract(tdName, td));
            return new InterfaceDesc(abs);
        });
    }

    private Ref<InterfaceDecl> arrayDecl() throws Exception {
        return this.getConstruct("arrayDecl", this.proj.interfaceDecls, () -> {
            final Ref<TypeParam> tdT = this.genT();
            final ArrayList<Ref<TypeParam>> tp = new ArrayList<Ref<TypeParam>>();
            tp.add(tdT);
            final Ref<InterfaceDesc> desc = this.arrayDesc("$T", tdT);
            return new InterfaceDecl(null, null, "$Array", desc, tp);
        });
    }
    
    public Ref<InterfaceInst> arrayInst(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("arrayInst<" + tdName + ">", this.proj.interfaceInsts, () -> {
            final Ref<InterfaceDecl> generic = this.arrayDecl();
            final ArrayList<Ref<? extends TypeDesc>> its = new ArrayList<Ref<? extends TypeDesc>>();
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
        return this.getConstruct("boxedBasic: " + SpoonUtils.describeElem(tr), this.proj.basics,
            () -> new Basic(basicName));
    }
}
