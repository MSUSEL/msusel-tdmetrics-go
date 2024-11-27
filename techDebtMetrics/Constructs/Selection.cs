using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

/// <summary>A selction of information from another, e.g. selction of a field from an object.</summary>
/// <see cref="../../docs/genFeatureDef.md#selection"/>
public class Selection : IConstruct, IInitializable {

    /// <summary>The name of the field, abstract, etc that is being selected from the origin.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The origin object, interface, etc that is being seleted out of.</summary>
    public IConstruct Origin => this.inOrigin ??
        throw new UninitializedException(nameof(this.Origin));
    private IConstruct? inOrigin;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.inOrigin = obj.ReadKey<IConstruct>("origin", project);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Origin).Write(".").Write(this.Name);
}
