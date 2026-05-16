package abstractor.core.constructs;

import java.util.Set;
import java.util.TreeSet;

import abstractor.core.AbstractorException;
import abstractor.core.cmp.Cmp;
import abstractor.core.cmp.CmpOptions;
import abstractor.core.json.*;

public class PackageCon extends ConstructImp {
    public String name;
    public String path;

    public final TreeSet<Ref<PackageCon>>    imports        = new TreeSet<>();
    public final TreeSet<Ref<InterfaceDecl>> interfaceDecls = new TreeSet<>();
    public final TreeSet<Ref<MethodDecl>>    methodDecls    = new TreeSet<>();
    public final TreeSet<Ref<ObjectDecl>>    objectDecls    = new TreeSet<>();
    public final TreeSet<Ref<Value>>         values         = new TreeSet<>();

    public PackageCon() {}

    public PackageCon(String name, String path) {
        this.name = name;
        this.path = path;
    }
    
    public ConstructKind kind() { return ConstructKind.PACKAGE; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("name",       this.name);
        obj.putNotEmpty("path",       this.path);
        obj.putNotEmpty("imports",    indexSet(imports));
        obj.putNotEmpty("interfaces", indexSet(this.interfaceDecls));
        obj.putNotEmpty("methods",    indexSet(this.methodDecls));
        obj.putNotEmpty("objects",    indexSet(this.objectDecls));
        obj.putNotEmpty("values",     indexSet(this.values));
        return obj;
    }

    @Override
    public Cmp getCmp(Construct c, CmpOptions options) {
        return Cmp.or(super.getCmp(c, options),
            Cmp.defer(this.name, () -> ((PackageCon)c).name),
            Cmp.defer(this.path, () -> ((PackageCon)c).path)
        );
    }

    static private <T extends Construct> boolean tryToAdd(Set<Ref<T>> set, Ref<? extends Construct> ref, ConstructKind kind) {
        if (ref.kind() != kind) return false;
        @SuppressWarnings("unchecked")
        Ref<T> cast = (Ref<T>)ref;
        set.add(cast);
        return true;
    }

    public void add(Ref<? extends Construct> ref) throws Exception {
        if (ref == null) return;
        if (tryToAdd(this.imports,        ref, ConstructKind.PACKAGE))        return;
        if (tryToAdd(this.objectDecls,    ref, ConstructKind.OBJECT_DECL))    return;
        if (tryToAdd(this.interfaceDecls, ref, ConstructKind.INTERFACE_DECL)) return;
        if (tryToAdd(this.methodDecls,    ref, ConstructKind.METHOD_DECL))    return;
        if (tryToAdd(this.values,         ref, ConstructKind.VALUE))          return;
        throw new AbstractorException("Unexpected construct type being added to package, " + this.name + " (" + this.path + "): " + ref.kind());
    }
}
