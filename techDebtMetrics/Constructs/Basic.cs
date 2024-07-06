using Constructs.Tooling;

namespace Constructs;

public class Basic : ITypeDesc, IInitializable {
    public string Name { get; private set; } = "";

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) =>
        this.Name = node.AsString();

    public override string ToString() => this.Name;

    public string ToStub() => this.Name;
}
