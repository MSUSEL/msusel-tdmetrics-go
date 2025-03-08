package abstractor.core.constructs;

public enum ConstructKind {
    ABSTRACT       ("abstract"),
    ARGUMENT       ("argument"),
    BASIC          ("basic"),
    FIELD          ("field"),
    INTERFACE_DECL ("interfaceDecl"),
    INTERFACE_DESC ("interfaceDesc"),
    INTERFACE_INST ("interfaceInst"),
    METHOD_DECL    ("method"),
    METHOD_INST    ("methodInst"),
    METRICS        ("metrics"),
    OBJECT_DECL    ("object"),
    OBJECT_INST    ("objectInst"),
    PACKAGE        ("package"),
    SELECTION      ("selection"),
    SIGNATURE      ("signature"),
    STRUCT_DESC    ("structDesc"),
    TYPE_PARAM     ("typeParam"),
    VALUE          ("value");


    static public ConstructKind fromName(String name) {
        for (ConstructKind k : ConstructKind.values()) {
            if (k.name.equals(name)) return k;
        }
        return null;
    }

    private final String name;
    ConstructKind(String name) { this.name = name; }

    public String toString() { return this.name; }

    public String plural() {
        return this.name.endsWith("s") ? this.name : (this.name + "s");
    }
}
