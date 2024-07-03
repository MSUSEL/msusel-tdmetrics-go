using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public class Basic : ITypeDesc, IInitializable {
    public string Name { get; private set; } = "";

    public void Initialize(TypeGetter getter, JsonNode node) =>
        this.Name = node.GetValue<string>();

    public override string ToString() => this.Name;
}
