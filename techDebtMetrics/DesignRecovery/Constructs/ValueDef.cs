using DesignRecovery.Extensions;
using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public class ValueDef : IInitializable {
    public string Name { get; private set; } = "";
    public bool IsConst { get; private set; } = false;

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
