using Constructs.Tooling;
using System.Text.Json.Nodes;

namespace Constructs;

public class Basic : ITypeDesc, IInitializable {
    public string Name { get; private set; } = "";

    void IInitializable.Initialize(TypeGetter getter, JsonNode node) =>
        this.Name = node.GetValue<string>();

    public override string ToString() => this.Name;

    public string ToStub() => this.Name;
}
