using System.ComponentModel;

namespace GroundTruth;

/// <summary>Indicates the type of the declaration type.</summary>
/// <see cref="mauricioaniche/ck/src/main/java/com/github/mauricioaniche/ck/CKVisitor.java:385"/>
public enum DeclType: int {

    /// <summary>An unknown declaration type.</summary>
    [Description("unknown")]
    Unknown = 0,

    /// <summary>A top-level (compilation-unit level) Java class.</summary>
    [Description("class")]
    Class = 1,

    /// <summary>
    /// Any nested class TypeDeclaration — static nested, non-static inner,
    /// or local (method-local) class. Not just JLS "inner".
    /// </summary>
    [Description("innerclass")]
    InnerClass = 2,

    /// <summary>Any Java interface, top-level or nested.</summary>
    /// <remarks>
    /// CK's isInterface() is checked before the stack check,
    /// so there is no "innerinterface" like there is for "class" and "innerclass".
    /// </remarks>
    [Description("interface")]
    Interface = 3,

    /// <summary>Any enum type, top-level or nested.</summary>
    /// <remarks>There is no "innerenum".</remarks>
    [Description("enum")]
    Enum = 4,

    /// <summary>Java 16+ record, top-level or nested.</summary>
    /// <remarks>There is no "innerrecord".</remarks>
    [Description("record")]
    Record = 5,

    /// <summary>
    /// Always an anonymous class. Java's AST only has AnonymousClassDeclaration for
    /// the new Foo() { ... } form, which is only ever an anonymous class (Java has no anonymous interfaces / enums / records).
    /// Lambdas are LambdaExpressions and are not emitted as CK class rows (they're counted via lambdasQty on the enclosing class).
    /// </summary>
    [Description("anonymous")]
    Anonymous = 6
}
