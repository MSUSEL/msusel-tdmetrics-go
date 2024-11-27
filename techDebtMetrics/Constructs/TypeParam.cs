using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

/// <summary>A type parameter name and type description.</summary>
/// <see cref="../../docs/genFeatureDef.md#type-parameter"/>
public class TypeParam : ITypeDesc, IInitializable {

    /// <summary>The name of the type parameter.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The type description of the type parameter.</summary>
    public ITypeDesc Type => this.inType ??
        throw new UninitializedException(nameof(this.Type));
    private ITypeDesc? inType;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name = obj.TryReadString("name");
        this.inType = obj.ReadKey<ITypeDesc>("type", project);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (j.Long)
            j.AsShort.Write(this.Type, suffix: " ");
        j.Write(this.Name);
    }
}
