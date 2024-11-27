using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

/// <summary>An argument is a parameter or result value in a signature.</summary>
/// <see cref="../../docs/genFeatureDef.md#argument"/>
public class Argument : IConstruct, IInitializable {

    /// <summary>Optional name of argument or an empty string.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The type of the argument.</summary>
    public ITypeDesc Type => this.inType ??
        throw new UninitializedException(nameof(this.Type));
    private ITypeDesc? inType;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name = obj.TryReadString("name");
        this.inType = obj.ReadKey<ITypeDesc>("type", project);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Type).Write(this.Name, prefix: " ");
}
