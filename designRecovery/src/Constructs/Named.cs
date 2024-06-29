using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal class Named : ITypeDesc{
    public string Name { get; private set; } = "";

    private ITypeDesc? inType;
    public ITypeDesc Type =>
        this.inType ?? throw new UninitializedException(nameof(Named), nameof(Type));

    public void Initialize(ITypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        
        this.Name = obj["name"]?.GetValue<string>() ??
            throw new MissingDataException(nameof(Named), "name");

        uint typeIndex = obj["type"]?.GetValue<uint>() ??
            throw new MissingDataException(nameof(Named), "type");
        this.inType = getter.Get<ITypeDesc>(typeIndex);
    }

    public override string ToString() => this.Name + ": " + this.Type ?? "<null>";
}
