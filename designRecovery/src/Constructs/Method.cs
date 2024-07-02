using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal class Method : IInitializer{
    
    /*
        name       string
        signature  TypeDesc
        metrics    metrics.Metrics

        Add(ctx2, `name`, m.name).
        Add(ctx2, `signature`, m.signature).
        AddNonZero(ctx2, `metrics`, m.metrics)
     */

    public void Initialize(TypeGetter getter, JsonNode node) {
    }
}
