using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

/// <summary>A package level value or constant declaration.</summary>
/// <see cref="../../docs/genFeatureDef.md#value"/>
public class Value : IDeclaration, IInitializable {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The name of the object declaration.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The location the object was defined.</summary>
    public Location Location { get; private set; }

    /// <summary>Intidcates if the value is constant, readonly, or assign once.</summary>
    public bool Constant { get; private set; }

    /// <summary>
    /// Metrics gathered from the initialization expression for the varlue
    /// or constant, if it had an initialization expression.
    /// </summary>
    public Metrics? Metrics { get; private set; }

    /// <summary>The type of the value or constant.</summary>
    public ITypeDesc Type => this.inType ??
        throw new UninitializedException(nameof(this.Type));
    private ITypeDesc? inType;

    /// <summary>The package the value or constant was declared in.</summary>
    public Package Package => this.inPackage ??
        throw new UninitializedException(nameof(this.Package));
    private Package? inPackage;

    void IInitializable.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.Location = obj.TryReadLocation("loc", project);
        this.Constant = obj.TryReadBool("const");
        this.Metrics = obj.TryReadIndex("metrics", project.Metrics);
        this.inType = obj.ReadKey<ITypeDesc>("type", project);
        this.inPackage = obj.ReadIndex("package", project.Packages);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (j.Long)
            j.Write(this.Constant ? "const " : "var ").
                AsShort.Write(this.Type, suffix: " ");
        j.Write(this.Name);
    }
}
