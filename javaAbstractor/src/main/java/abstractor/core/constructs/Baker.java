package abstractor.core.constructs;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.SortedSet;
import java.util.TreeMap;
import java.util.TreeSet;

public class Baker {
    private final Project proj;
    private final Map<String, Construct> cache;

    public Baker(Project proj) {
        this.proj  = proj;
        this.cache = new TreeMap<String, Construct>();
    }

    public interface ConstructCreator<T extends Construct> { T create(); }

    @SuppressWarnings("unchecked")
    private <T extends Construct> T getConstruct(String name, ConstructCreator<T> creator) {
        Construct existing = cache.get(name);
        if (existing != null) return (T)existing;
        T value = creator.create();
        cache.put(name, value);
        return value;
    }

    public InterfaceDesc objectDesc() {
        return this.getConstruct("objectInterfaceDesc", () -> {
            final InterfaceDesc desc = new InterfaceDesc(Collections.emptySortedSet());
            return this.proj.interfaceDescs.addOrGet(desc);
        });
    }
    
    private Basic intBasic() {
        return this.getConstruct("intBasic", () -> {
            return this.proj.basics.addOrGet(new Basic("int"));
        });
    }

    private TypeParam genT() {
        return this.getConstruct("genT", () -> {
            final TypeParam tp = new TypeParam("T", this.objectDesc());
            return this.proj.typeParams.addOrGet(tp);
        });
    }

    private Argument intReturn() {
        return this.getConstruct("intReturn", () -> {
            return this.proj.arguments.addOrGet(new Argument(this.intBasic()));
        });
    }

    private Argument intIndexParam() {
        return this.getConstruct("intIndexParam", () -> {
            return this.proj.arguments.addOrGet(new Argument("index", this.intBasic()));
        });
    }

    private Argument genReturn(String tdName, TypeDesc td) {
        return this.getConstruct("genReturn<" + tdName + ">", () -> {
            return this.proj.arguments.addOrGet(new Argument(td));
        });
    }

    private Argument valueGenParam(String tdName, TypeDesc td) {
        return this.getConstruct("valueGenParam<" + tdName + ">", () -> {
            return this.proj.arguments.addOrGet(new Argument("value", td));
        });
    }

    private Signature lenSignature() {
        return this.getConstruct("lenSignature", () -> {
            final List<Argument> results = new ArrayList<Argument>();
            results.add(this.intReturn());
            final Signature sig = new Signature(false, Collections.emptyList(), results);
            return this.proj.signatures.addOrGet(sig);
        });
    }

    private Signature getSignature(String tdName, TypeDesc td) {
        return this.getConstruct("getSignature<" + tdName + ">", () -> {
            final List<Argument> params = new ArrayList<Argument>();
            params.add(this.intIndexParam());
            final List<Argument> results = new ArrayList<Argument>();
            results.add(this.genReturn(tdName, td));
            final Signature sig = new Signature(false, params, results);
            return this.proj.signatures.addOrGet(sig);
        });
    }

    private Signature setSignature(String tdName, TypeDesc td) {
        return this.getConstruct("setSignature<" + tdName + ">", () -> {
            final List<Argument> params = new ArrayList<Argument>();
            params.add(this.intIndexParam());
            params.add(this.valueGenParam(tdName, td));
            final Signature sig = new Signature(false, params, Collections.emptyList());
            return this.proj.signatures.addOrGet(sig);
        });
    }

    private Abstract lenAbstract() {
        return this.getConstruct("lenAbstract", () -> {
            final Abstract abs = new Abstract("$len", this.lenSignature());
            return this.proj.abstracts.addOrGet(abs);
        });
    }

    private Abstract getAbstract(String tdName, TypeDesc td) {
        return this.getConstruct("getAbstract<" + tdName + ">", () -> {
            final Abstract abs = new Abstract("$get", this.getSignature(tdName, td));
            return this.proj.abstracts.addOrGet(abs);
        });
    }

    private Abstract setAbstract(String tdName, TypeDesc td) {
        return this.getConstruct("setAbstract<" + tdName + ">", () -> {
            final Abstract abs = new Abstract("$set", this.setSignature(tdName, td));
            return this.proj.abstracts.addOrGet(abs);
        });
    }

    private InterfaceDesc arrayDesc(String tdName, TypeDesc td) {
        return this.getConstruct("arrayDesc<" + tdName + ">", () -> {
            final SortedSet<Abstract> abs = new TreeSet<Abstract>();
            abs.add(this.lenAbstract());
            abs.add(this.getAbstract(tdName, td));
            abs.add(this.setAbstract(tdName, td));
            final InterfaceDesc desc = new InterfaceDesc(abs);
            return this.proj.interfaceDescs.addOrGet(desc);
        });
    }

    private InterfaceDecl arrayDecl() {
        return this.getConstruct("arrayDecl", () -> {
            final TypeParam tdT = this.genT();
            final List<TypeParam> tp = new ArrayList<TypeParam>();
            tp.add(tdT);
            final InterfaceDesc desc = this.arrayDesc("$T", tdT);
            final InterfaceDecl decl = new InterfaceDecl(null, null, "$Array", desc, tp);
            return this.proj.interfaceDecls.addOrGet(decl);
        });
    }
    
    public InterfaceInst arrayInst(String tdName, TypeDesc td) {
        final InterfaceDecl generic = this.arrayDecl();
        final List<TypeDesc> its = new ArrayList<TypeDesc>();
        its.add(td);
        final InterfaceDesc resolved = this.arrayDesc(tdName, td);
        final InterfaceInst inst = new InterfaceInst(generic, its, resolved);
        return this.proj.interfaceInsts.addOrGet(inst);
    }
}
