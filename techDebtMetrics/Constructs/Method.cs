using Constructs.Exceptions;
using Constructs.Extensions;
using Constructs.Tooling;
using System.Text.Json.Nodes;

namespace Constructs;

public class Method : IConstruct, IInitializable {
    public string Name { get; private set; } = "";

    public ITypeDesc Signature => this.inSignature ??
        throw new UninitializedException("signature");
    private ITypeDesc? inSignature;

    void IInitializable.Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        this.Name = obj.ReadValue<string>("name");
        this.inSignature = obj.ReadIndexType<ITypeDesc>("signature", getter);

        // TODO: Add "metrics"
    }

    public string ToStub() => this.Name + " " + this.Signature.ToStub();
}
