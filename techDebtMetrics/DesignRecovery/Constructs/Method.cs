using DesignRecovery.Extensions;
using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public class Method : IInitializable {
    public string Name { get; private set; } = "";

    private ITypeDesc? inSignature;
    public ITypeDesc Signature => this.inSignature ??
        throw new UninitializedException("signature");

    public void Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        this.Name = obj.ReadValue<string>("name");
        this.inSignature = obj.ReadIndexType<ITypeDesc>("signature", getter);

        // TODO: Add "metrics"
    }
}
