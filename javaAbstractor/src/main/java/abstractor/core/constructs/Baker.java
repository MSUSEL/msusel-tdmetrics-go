package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.Collections;
import java.util.Map;
import java.util.TreeMap;
import java.util.TreeSet;

public class Baker {
    private final Project proj;
    private final Map<String, Construct> cache;

    public Baker(Project proj) {
        this.proj  = proj;
        this.cache = new TreeMap<String, Construct>();
    }

    public interface ConstructCreator<T extends Construct> { T create() throws Exception; }

    @SuppressWarnings("unchecked")
    private <T extends Construct> T getConstruct(String name, ConstructCreator<T> creator) throws Exception {
        Construct existing = cache.get(name);
        if (existing != null) return (T)existing;
        T value = creator.create();
        cache.put(name, value);
        return value;
    }

    public Ref<InterfaceDesc> objectDesc() throws Exception {
        return this.getConstruct("objectInterfaceDesc", () -> {
            final InterfaceDesc desc = new InterfaceDesc(Collections.emptySortedSet());
            return this.proj.interfaceDescs.addOrGetRef(desc);
        });
    }
    
    private Ref<Basic> intBasic() throws Exception {
        return this.getConstruct("intBasic", () -> {
            return this.proj.basics.addOrGetRef(new Basic("int"));
        });
    }

    private Ref<TypeParam> genT() throws Exception {
        return this.getConstruct("genT", () -> {
            final TypeParam tp = new TypeParam("T", this.objectDesc());
            return this.proj.typeParams.addOrGetRef(tp);
        });
    }

    private Ref<Argument> intReturn() throws Exception {
        return this.getConstruct("intReturn", () -> {
            return this.proj.arguments.addOrGetRef(new Argument(this.intBasic()));
        });
    }

    private Ref<Argument> intIndexParam() throws Exception {
        return this.getConstruct("intIndexParam", () -> {
            return this.proj.arguments.addOrGetRef(new Argument("index", this.intBasic()));
        });
    }

    private Ref<Argument> genReturn(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("genReturn<" + tdName + ">", () -> {
            return this.proj.arguments.addOrGetRef(new Argument(td));
        });
    }

    private Ref<Argument> valueGenParam(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("valueGenParam<" + tdName + ">", () -> {
            return this.proj.arguments.addOrGetRef(new Argument("value", td));
        });
    }

    private Ref<Signature> lenSignature() throws Exception {
        return this.getConstruct("lenSignature", () -> {
            final ArrayList<Ref<Argument>> results = new ArrayList<Ref<Argument>>();
            results.add(this.intReturn());
            final Signature sig = new Signature(false, Collections.emptyList(), results);
            return this.proj.signatures.addOrGetRef(sig);
        });
    }

    private Ref<Signature> getSignature(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("getSignature<" + tdName + ">", () -> {
            final ArrayList<Ref<Argument>> params = new ArrayList<Ref<Argument>>();
            params.add(this.intIndexParam());
            final ArrayList<Ref<Argument>> results = new ArrayList<Ref<Argument>>();
            results.add(this.genReturn(tdName, td));
            final Signature sig = new Signature(false, params, results);
            return this.proj.signatures.addOrGetRef(sig);
        });
    }

    private Ref<Signature> setSignature(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("setSignature<" + tdName + ">", () -> {
            final ArrayList<Ref<Argument>> params = new ArrayList<Ref<Argument>>();
            params.add(this.intIndexParam());
            params.add(this.valueGenParam(tdName, td));
            final Signature sig = new Signature(false, params, Collections.emptyList());
            return this.proj.signatures.addOrGetRef(sig);
        });
    }

    private Ref<Abstract> lenAbstract() throws Exception {
        return this.getConstruct("lenAbstract", () -> {
            final Abstract abs = new Abstract("$len", this.lenSignature());
            return this.proj.abstracts.addOrGetRef(abs);
        });
    }

    private Ref<Abstract> getAbstract(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("getAbstract<" + tdName + ">", () -> {
            final Abstract abs = new Abstract("$get", this.getSignature(tdName, td));
            return this.proj.abstracts.addOrGetRef(abs);
        });
    }

    private Ref<Abstract> setAbstract(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("setAbstract<" + tdName + ">", () -> {
            final Abstract abs = new Abstract("$set", this.setSignature(tdName, td));
            return this.proj.abstracts.addOrGetRef(abs);
        });
    }

    private Ref<InterfaceDesc> arrayDesc(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        return this.getConstruct("arrayDesc<" + tdName + ">", () -> {
            final TreeSet<Ref<Abstract>> abs = new TreeSet<Ref<Abstract>>();
            abs.add(this.lenAbstract());
            abs.add(this.getAbstract(tdName, td));
            abs.add(this.setAbstract(tdName, td));
            final InterfaceDesc desc = new InterfaceDesc(abs);
            return this.proj.interfaceDescs.addOrGetRef(desc);
        });
    }

    private Ref<InterfaceDecl> arrayDecl() throws Exception {
        return this.getConstruct("arrayDecl", () -> {
            final Ref<TypeParam> tdT = this.genT();
            final ArrayList<Ref<TypeParam>> tp = new ArrayList<Ref<TypeParam>>();
            tp.add(tdT);
            final Ref<InterfaceDesc> desc = this.arrayDesc("$T", tdT);
            final InterfaceDecl decl = new InterfaceDecl(null, null, "$Array", desc, tp);
            return this.proj.interfaceDecls.addOrGetRef(decl);
        });
    }
    
    public Ref<InterfaceInst> arrayInst(String tdName, Ref<? extends TypeDesc> td) throws Exception {
        final Ref<InterfaceDecl> generic = this.arrayDecl();
        final ArrayList<Ref<? extends TypeDesc>> its = new ArrayList<Ref<? extends TypeDesc>>();
        its.add(td);
        final Ref<InterfaceDesc> resolved = this.arrayDesc(tdName, td);
        final InterfaceInst inst = new InterfaceInst(generic, its, resolved);
        return this.proj.interfaceInsts.addOrGetRef(inst);
    }
}
