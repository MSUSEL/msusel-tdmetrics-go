using DesignRecovery.Extensions;
using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public class Named : ITypeDesc, IInitializable {
    public string Name { get; private set; } = "";

    private ITypeDesc? inType;
    public ITypeDesc Type => this.inType ??
        throw new UninitializedException("type");

    public void Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        this.Name = obj.ReadValue<string>("name");
        this.inType = obj.ReadIndexType<ITypeDesc>("type", getter);
    }

    public override string ToString() => this.Name + ": " + this.Type ?? "<null>";
}
