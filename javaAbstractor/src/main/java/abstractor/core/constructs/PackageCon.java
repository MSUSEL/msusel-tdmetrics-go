package abstractor.core.constructs;

import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;

public class PackageCon extends ConstructImp {
    public String name;
    public String path;

    public final TreeSet<Ref<PackageCon>>    imports        = new TreeSet<Ref<PackageCon>>();
    public final TreeSet<Ref<InterfaceDecl>> interfaceDecls = new TreeSet<Ref<InterfaceDecl>>();
    public final TreeSet<Ref<MethodDecl>>    methodDecls    = new TreeSet<Ref<MethodDecl>>();
    public final TreeSet<Ref<ObjectDecl>>    objectDecls    = new TreeSet<Ref<ObjectDecl>>();
    public final TreeSet<Ref<Value>>         values         = new TreeSet<Ref<Value>>();

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
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.name, () -> ((PackageCon)c).name),
            Cmp.defer(this.path, () -> ((PackageCon)c).path)
        );
    }
}
