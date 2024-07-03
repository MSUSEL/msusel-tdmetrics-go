// Ignore Spelling: Initializable

using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public interface IInitializable {
    void Initialize(TypeGetter getter, JsonNode node);
}
