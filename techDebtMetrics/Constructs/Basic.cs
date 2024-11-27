using Constructs.Data;
using Constructs.Tooling;

namespace Constructs;

/// <summary>
/// A basic type such as int, string, or bool. 
/// These may be unique to a language.
/// </summary>
/// <see cref="../../docs/genFeatureDef.md#basic"/>
public class Basic : ITypeDesc, IInitializable {

    /// <summary>This is the name for the basic type.</summary>
    public string Name { get; private set; } = "";

    void IInitializable.Initialize(Project project, Node node) =>
        this.Name = node.IsScalar ? node.AsString() : node.AsObject().ReadString("name");

    public override string ToString() => this.Name;

    public void ToStub(Journal j) => j.Write(this.Name);
}
