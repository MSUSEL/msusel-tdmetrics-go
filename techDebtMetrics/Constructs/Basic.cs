using Constructs.Data;
using Constructs.Tooling;

namespace Constructs;

/// <summary>
/// A basic type such as int, string, or bool. 
/// These may be unique to a language.
/// </summary>
/// <see cref="../../docs/genFeatureDef.md#basic"/>
public class Basic : ITypeDesc, IInitializable {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>This is the name for the basic type.</summary>
    public string Name { get; private set; } = "";

    void IInitializable.Initialize(Project project, int index, Node node) {
        this.Index = index;
        this.Name = node.IsScalar ? node.AsString() : node.AsObject().ReadString("name");
    }

    public override string ToString() => this.Name;

    public void ToStub(Journal j) => j.Write(this.Name);
}
