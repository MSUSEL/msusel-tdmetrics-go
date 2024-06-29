using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal class Basic : ITypeDesc {
    public string Name { get; private set; } = "";

    public void Initialize(ITypeGetter getter, JsonNode node) =>
        this.Name = node.AsValue().GetValue<string>();

    public override string ToString() => this.Name;
}
