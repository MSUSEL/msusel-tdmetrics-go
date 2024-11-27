using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

/// <summary>A field is named types in a structure.</summary>
/// <see cref="../../docs/genFeatureDef.md#field"/>
public class Field : IConstruct, IInitializable {

    /// <summary>The name of the field.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The type of the field.</summary>
    public ITypeDesc Type => this.inType ??
        throw new UninitializedException(nameof(this.Type));
    private ITypeDesc? inType;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.inType = obj.ReadKey<ITypeDesc>("type", project);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (j.Long)
            j.AsShort.Write(this.Type, suffix: " ");
        j.Write(this.Name);
    }
}
