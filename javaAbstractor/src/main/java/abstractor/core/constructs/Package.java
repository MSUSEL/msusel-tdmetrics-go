package abstractor.core.constructs;

import java.util.TreeSet;

import abstractor.core.cmp.Cmp;
import abstractor.core.json.*;
import spoon.reflect.cu.SourcePosition;
import spoon.reflect.declaration.CtPackage;

public class Package extends Construct {
    public final CtPackage pkg;
    public final String name;
    public final String path;

    public final TreeSet<Package> imports = new TreeSet<Package>();
    public final TreeSet<InterfaceDecl> interfaceDecls = new TreeSet<InterfaceDecl>();
    public final TreeSet<MethodDecl> methodDecls = new TreeSet<MethodDecl>();
    public final TreeSet<ObjectDecl> objectDecls = new TreeSet<ObjectDecl>();
    // TODO: | `values`     | ⬤ | ◯ | List of [indices](#indices) of [values](#value) declared in this package. |
    
    static public Package Create(Project proj, CtPackage pkg) {
        Package existing = proj.packages.findWithSource(pkg);
        if (existing != null) return existing;
        return proj.packages.tryAdd(new Package(pkg));
    }

    static private String packagePath(CtPackage p) {
        SourcePosition pos = p.getPosition();
        if (!pos.isValidPosition()) return "";
        
        String path = pos.getFile().getPath();
        final String tail = "package-info.java";
        if (path.endsWith(tail))
            path = path.substring(0, path.length()-tail.length());
        return path;
    }

    private Package(CtPackage pkg) {
        this.pkg  = pkg;
        this.name = pkg.getQualifiedName();
        this.path = packagePath(pkg);
    }
    
    public Object source() { return this.pkg; }
    public String kind() { return "package"; }

    @Override
    public int compareTo(Construct c) {
        return Cmp.or(
            () -> super.compareTo(c),
            Cmp.defer(this.name, () -> ((Package)c).name),
            Cmp.defer(this.path, () -> ((Package)c).path)
        );
    }

    @Override
    public JsonNode toJson(JsonHelper h) {
        JsonObject obj = (JsonObject)super.toJson(h);
        obj.put("name", this.name);
        obj.putNotEmpty("path", this.path);
        obj.putNotEmpty("imports", indexSet(imports));
        obj.putNotEmpty("interfaces", Construct.indexSet(this.interfaceDecls));
        obj.putNotEmpty("methods", Construct.indexSet(this.methodDecls));
        obj.putNotEmpty("objects", Construct.indexSet(this.objectDecls));
        // TODO: values
        return obj;
    }
}
