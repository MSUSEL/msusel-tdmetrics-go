using Constructs.Data;
using Constructs.Tooling;

namespace Constructs;

public class Basic : ITypeDesc, IInitializable {
    public string Name { get; private set; } = "";

    void IInitializable.Initialize(Project project, Node node) =>
        this.Name = node.IsScalar ? node.AsString() :
            this.Name = node.AsObject().ReadString("name");

    public override string ToString() => this.Name;

    public void ToStub(Journal j) => j.Write(this.Name);
}
