using designRecovery.src.Extensions;
using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal class ValueDef : IInitializer {
    public string Name { get; private set; } = "";
    public bool IsConst  { get; private set; }

    private ITypeDesc? inType;
    public ITypeDesc Type => this.inType ??
        throw new UninitializedException("type");

    public void Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        this.Name = obj.ReadValue<string>("name");
        this.IsConst = obj.ReadValue<bool>("isConst");
        this.inType = obj.ReadIndexType<ITypeDesc>("type", getter);
    }
}
