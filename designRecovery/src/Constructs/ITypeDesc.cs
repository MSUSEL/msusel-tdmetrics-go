using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal interface ITypeDesc {
    public void Initialize(ITypeGetter getter, JsonNode node);
}
