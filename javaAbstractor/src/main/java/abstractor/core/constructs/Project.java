package abstractor.core.constructs;

import abstractor.core.json.*;

public class Project implements Jsonable {
    public final Baker                   baker          = new Baker(this);
    public final Locations               locations      = new Locations();
    public final Factory<Abstract>       abstracts      = new Factory<Abstract>      (ConstructKind.ABSTRACT,       () -> new Abstract());
    public final Factory<Argument>       arguments      = new Factory<Argument>      (ConstructKind.ARGUMENT,       () -> new Argument());
    public final Factory<Basic>          basics         = new Factory<Basic>         (ConstructKind.BASIC,          () -> new Basic());
    public final Factory<Field>          fields         = new Factory<Field>         (ConstructKind.FIELD,          () -> new Field());
    public final Factory<InterfaceDecl>  interfaceDecls = new Factory<InterfaceDecl> (ConstructKind.INTERFACE_DECL, () -> new InterfaceDecl());
    public final Factory<InterfaceDesc>  interfaceDescs = new Factory<InterfaceDesc> (ConstructKind.INTERFACE_DESC, () -> new InterfaceDesc());
    public final Factory<InterfaceInst>  interfaceInsts = new Factory<InterfaceInst> (ConstructKind.INTERFACE_INST, () -> new InterfaceInst());
    public final Factory<MethodDecl>     methodDecls    = new Factory<MethodDecl>    (ConstructKind.METHOD_DECL,    () -> new MethodDecl());
    public final Factory<MethodInst>     methodInsts    = new Factory<MethodInst>    (ConstructKind.METHOD_INST,    () -> new MethodInst());
    public final Factory<Metrics>        metrics        = new Factory<Metrics>       (ConstructKind.METRICS,        () -> new Metrics());
    public final Factory<ObjectDecl>     objectDecls    = new Factory<ObjectDecl>    (ConstructKind.OBJECT_DECL,    () -> new ObjectDecl());
    public final Factory<ObjectInst>     objectInsts    = new Factory<ObjectInst>    (ConstructKind.OBJECT_INST,    () -> new ObjectInst());
    public final Factory<PackageCon>     packages       = new Factory<PackageCon>    (ConstructKind.PACKAGE,        () -> new PackageCon());
    public final Factory<Selection>      selections     = new Factory<Selection>     (ConstructKind.SELECTION,      () -> new Selection());
    public final Factory<Signature>      signatures     = new Factory<Signature>     (ConstructKind.SIGNATURE,      () -> new Signature());
    public final Factory<StructDesc>     structDescs    = new Factory<StructDesc>    (ConstructKind.STRUCT_DESC,    () -> new StructDesc());
    public final Factory<TypeParam>      typeParams     = new Factory<TypeParam>     (ConstructKind.TYPE_PARAM,     () -> new TypeParam());
    public final Factory<Value>          values         = new Factory<Value>         (ConstructKind.VALUE,          () -> new Value());

    public final Factory<?>[] factories = new Factory<?>[] {
        this.abstracts,
        this.arguments,
        this.basics,
        this.fields,
        this.interfaceDecls,
        this.interfaceDescs,
        this.interfaceInsts,
        this.methodDecls,
        this.methodInsts,
        this.metrics,
        this.objectDecls,
        this.objectInsts,
        this.packages,
        this.selections,
        this.signatures,
        this.structDescs,
        this.typeParams,
        this.values
    };

    public final Factory<?>[] declarations = new Factory<?>[] {
        this.interfaceDecls,
        this.methodDecls,
        this.objectDecls,
        this.values,
    };

    private void setAllIndices() {
        for (Factory<?> factory : this.factories)
            factory.setIndices();
    }

    public JsonNode toJson(JsonHelper h) {
        this.locations.prepareForOutput();
        this.setAllIndices();

        JsonObject obj = new JsonObject();
        obj.put("language", "java");
        obj.putNotEmpty("locs", this.locations.toJson(h));
        for (Factory<?> factory : this.factories)
            obj.putNotEmpty(factory.kind().plural(), factory.toJson(h));
        return obj;
    }

    public Factory<?> getFactory(ConstructKind kind) {
        for (Factory<?> factory : factories) {
            if (factory.kind().equals(kind)) return factory;
        }
        return null;
    }

    public Construct getConstructWithKey(String key) {
        int split = -1;
        for (int i = 0; i < key.length(); i++) {
            final char c = key.charAt(i);
            if (c >= '0' && c <= '9') {
                split = i;
                break;
            }
        }

        if (split < 0) {
            System.out.println("Failed to find split point for key " + key + ".");
            return null;
        }

        final String kindStr = key.substring(0, split);
        final String indStr  = key.substring(split);

        final ConstructKind kind = ConstructKind.fromName(kindStr);
        if (kind == null) {
            System.out.println("Failed to find kind with " + kindStr + " for " + key + ".");
            return null;
        }

        final Factory<?> factory = getFactory(kind);
        if (factory == null) {
            System.out.println("Failed to find factory of kind " + kind + " for " + key + ".");
            return null;
        }

        int index = 0;
        try {
            index = Integer.parseInt(indStr) - 1;
        } catch(Exception ex) {
            System.out.println("Failed to get index with " + indStr + " for " + key + ".");
            return null;
        }

        final Construct con = factory.get(index);
        if (con == null) {
            System.out.println("Failed to get construct from factory " + factory.kind() +
                " (size: " + factory.size() + ") with index " + index + " for " + key + ".");
            return null;
        }

        return con;
    }
}
