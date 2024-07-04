using Constructs.Exceptions;
using Constructs.Extensions;
using Constructs.Tooling;
using System.Text.Json.Nodes;

namespace Constructs;

public class Named : ITypeDesc, IInitializable {
    public string Name { get; private set; } = "";

    public ITypeDesc Type => this.inType ??
        throw new UninitializedException("type");
    private ITypeDesc? inType;

    void IInitializable.Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        this.Name = obj.ReadValue<string>("name");
        this.inType = obj.ReadIndexType<ITypeDesc>("type", getter);
    }

    public override string ToString() => this.ToStub();

    public string ToStub() => this.Name + ": " + this.Type ?? "<null>";
}
