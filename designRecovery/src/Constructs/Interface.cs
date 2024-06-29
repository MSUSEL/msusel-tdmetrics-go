using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal class Interface: ITypeDesc {
    private readonly List<Named> inTypeParamss = [];
    public IReadOnlyList<Named> TypeParams => this.inTypeParamss.AsReadOnly();

    private readonly List<Interface> inInterfaces = [];
    public IReadOnlyList<Interface> Interfaces => this.inInterfaces.AsReadOnly();

    private readonly Dictionary<string, ITypeDesc> inMethods = [];
    public IReadOnlyDictionary<string, ITypeDesc> Methods => this.inMethods.AsReadOnly();

    public Union? Union { get; private set; }

    public void Initialize(ITypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();


        // TODO: Finish
    }
}
