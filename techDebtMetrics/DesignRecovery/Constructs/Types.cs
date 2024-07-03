// Ignore Spelling: Structs

using System.Collections.Generic;
using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public class Types {
    private readonly List<Basic> inBasics = [];
    public IReadOnlyList<Basic> Basics => this.inBasics.AsReadOnly<Basic>();

    private readonly List<Interface> inInterfaces = [];
    public IReadOnlyList<Interface> Interfaces => this.inInterfaces.AsReadOnly<Interface>();

    private readonly List<Named> inNamed = [];
    public IReadOnlyList<Named> Named => this.inNamed.AsReadOnly<Named>();

    private readonly List<Signature> inSignatures = [];
    public IReadOnlyList<Signature> Signatures => this.inSignatures.AsReadOnly<Signature>();

    private readonly List<Solid> inSolids = [];
    public IReadOnlyList<Solid> Solids => this.inSolids.AsReadOnly<Solid>();

    private readonly List<Struct> inStructs = [];
    public IReadOnlyList<Struct> Structs => this.inStructs.AsReadOnly<Struct>();

    private readonly List<Union> inUnions = [];
    public IReadOnlyList<Union> Unions => this.inUnions.AsReadOnly<Union>();

    static private void preallocate<T>(JsonObject obj, string name, List<T> list) where T : new() {
        int count = obj[name]?.AsArray()?.Count ?? 0;
        for (int i = 0; i < count; i++)
            list[i] = new T();
    }

    public Types(JsonNode node) {
        JsonObject obj = node.AsObject();
        preallocate(obj, "basics", this.inBasics);
        preallocate(obj, "interfaces", this.inInterfaces);
        preallocate(obj, "named", this.inNamed);
        preallocate(obj, "signatures", this.inSignatures);
        preallocate(obj, "solids", this.inSolids);
        preallocate(obj, "structs", this.inStructs);
        preallocate(obj, "unions", this.inUnions);
    }

    static private void initType<T>(JsonObject obj, TypeGetter getter, string name, List<T> list) where T : ITypeDesc {
        JsonArray? array = obj[name]?.AsArray();
        if (array is null) return;
        int count = array.Count;
        for (int i = 0; i < count; i++) {
            JsonNode? node = array[i];
            if (node is null) continue;
            list[i].Initialize(getter, node);
        }
    }

    public void Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        initType(obj, getter, "basic", this.inBasics);
        initType(obj, getter, "interfaces", this.inInterfaces);
        initType(obj, getter, "named", this.inNamed);
        initType(obj, getter, "signatures", this.inSignatures);
        initType(obj, getter, "solids", this.inSolids);
        initType(obj, getter, "structs", this.inStructs);
        initType(obj, getter, "unions", this.inUnions);
    }
}
