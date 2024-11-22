using Constructs.Data;
using Constructs.Tooling;

namespace Constructs;

public class Basic : ITypeDesc, IInitializable {
    public string Name { get; private set; } = "";

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        if (node.IsScalar) this.Name = node.AsString();
        else {
            Data.Object obj = node.AsObject();
            this.Name = obj.ReadString("name");
        }
    }

    public override string ToString() => this.Name;

    public void ToStub(Journal j) => j.Write(this.Name);
}
