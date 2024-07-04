using Constructs.Exceptions;
using Constructs.Extensions;
using Constructs.Tooling;
using System.Text.Json.Nodes;

namespace Constructs;

public class ValueDef : IConstruct, IInitializable {
    public string Name { get; private set; } = "";
    public bool IsConst { get; private set; } = false;

    public ITypeDesc Type => this.inType ??
        throw new UninitializedException("type");
    private ITypeDesc? inType;

    void IInitializable.Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        this.Name = obj.ReadValue<string>("name");
        this.IsConst = obj.ReadValue<bool>("isConst");
        this.inType = obj.ReadIndexType<ITypeDesc>("type", getter);
    }
    
    public string ToStub() => (this.IsConst?"const ":"") + this.Name + " "+ this.Type.ToStub();
}
