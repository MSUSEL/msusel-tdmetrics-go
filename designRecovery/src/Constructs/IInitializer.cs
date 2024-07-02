using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal interface IInitializer {
    public void Initialize(TypeGetter getter, JsonNode node);
}
