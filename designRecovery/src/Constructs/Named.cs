using designRecovery.src.Extensions;
using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal class Named : ITypeDesc{
    public string Name { get; private set; } = "";

    private ITypeDesc? inType;
    public ITypeDesc Type => this.inType ??
        throw new UninitializedException("type");

    public void Initialize(ITypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        this.Name = obj.ReadValue<string>("name");
        this.inType = obj.ReadIndexType<ITypeDesc>("type", getter);
    }

    public override string ToString() => this.Name + ": " + this.Type ?? "<null>";
}
