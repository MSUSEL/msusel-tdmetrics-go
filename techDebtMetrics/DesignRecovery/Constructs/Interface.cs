using DesignRecovery.Extensions;
using System.Collections.Generic;
using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public class Interface : ITypeDesc, IInitializable {
    private readonly List<Named> inTypeParams = [];
    public IReadOnlyList<Named> TypeParams => this.inTypeParams.AsReadOnly();

    private readonly List<Interface> inInterfaces = [];
    public IReadOnlyList<Interface> Interfaces => this.inInterfaces.AsReadOnly();

    private readonly Dictionary<string, ITypeDesc> inMethods = [];
    public IReadOnlyDictionary<string, ITypeDesc> Methods => this.inMethods.AsReadOnly();

    public Union? Union { get; private set; }

    public void Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();

        obj.ReadIndexTypeList("typeParams", getter, this.inTypeParams);
        obj.ReadIndexTypeList("interfaces", getter, this.inInterfaces);

        JsonObject? methodObj = obj["methods"]?.AsObject();
        if (methodObj is not null) {
            foreach (KeyValuePair<string, JsonNode?> pair in methodObj) {
                string key = pair.Key;
                uint typeIndex = pair.Value?.GetValue<uint>() ??
                    throw new MissingDataException("methods." + key + ".value");
                ITypeDesc value = getter.GetTypeAtIndex<ITypeDesc>(typeIndex);
                this.inMethods[key] = value;
            }
        }

        if (obj.ContainsKey("union"))
            this.Union = obj.ReadIndexType<Union>("union", getter);
    }
}
