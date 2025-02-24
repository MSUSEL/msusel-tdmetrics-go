package abstractor.core.constructs;

import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.declaration.CtPackage;

public class PackageCon extends ConstructImp {
    public final String name;
    public final String path;

    public final TreeSet<PackageCon> imports = new TreeSet<PackageCon>();
    public final TreeSet<InterfaceDecl> interfaceDecls = new TreeSet<InterfaceDecl>();
    public final TreeSet<MethodDecl> methodDecls = new TreeSet<MethodDecl>();
    public final TreeSet<ObjectDecl> objectDecls = new TreeSet<ObjectDecl>();
    public final TreeSet<Value> values = new TreeSet<Value>();

    public PackageCon(CtPackage src, String name, String path) {
        super(src);
        this.name = name;
        this.path = path;
    }
    
    public String kind() { return "package"; }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.putNotEmpty("name", this.name);
        obj.putNotEmpty("path", this.path);
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
