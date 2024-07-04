using System.Text.Json.Nodes;

namespace Constructs.Tooling;

internal interface IInitializable {

    void Initialize(TypeGetter getter, JsonNode node);
}
