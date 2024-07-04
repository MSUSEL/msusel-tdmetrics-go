// Ignore Spelling: Structs

using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Text.Json.Nodes;

namespace Constructs;

public class Types : IInitializable {
    public IReadOnlyList<Basic> Basics => this.inBasics.AsReadOnly<Basic>();
    private readonly List<Basic> inBasics = [];

    public IReadOnlyList<Interface> Interfaces => this.inInterfaces.AsReadOnly<Interface>();
    private readonly List<Interface> inInterfaces = [];

    public IReadOnlyList<Named> Named => this.inNamed.AsReadOnly<Named>();
    private readonly List<Named> inNamed = [];

    public IReadOnlyList<Signature> Signatures => this.inSignatures.AsReadOnly<Signature>();
    private readonly List<Signature> inSignatures = [];

    public IReadOnlyList<Solid> Solids => this.inSolids.AsReadOnly<Solid>();
    private readonly List<Solid> inSolids = [];

    public IReadOnlyList<Struct> Structs => this.inStructs.AsReadOnly<Struct>();
    private readonly List<Struct> inStructs = [];

    public IReadOnlyList<Union> Unions => this.inUnions.AsReadOnly<Union>();
    private readonly List<Union> inUnions = [];

    public Types(JsonNode node) {
        JsonObject obj = node.AsObject();
        obj.PreallocateList("basics", this.inBasics);
        obj.PreallocateList("interfaces", this.inInterfaces);
        obj.PreallocateList("named", this.inNamed);
        obj.PreallocateList("signatures", this.inSignatures);
        obj.PreallocateList("solids", this.inSolids);
        obj.PreallocateList("structs", this.inStructs);
        obj.PreallocateList("unions", this.inUnions);
    }

    void IInitializable.Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        obj.InitializeList(getter, "basic", this.inBasics);
        obj.InitializeList(getter, "interfaces", this.inInterfaces);
        obj.InitializeList(getter, "named", this.inNamed);
        obj.InitializeList(getter, "signatures", this.inSignatures);
        obj.InitializeList(getter, "solids", this.inSolids);
        obj.InitializeList(getter, "structs", this.inStructs);
        obj.InitializeList(getter, "unions", this.inUnions);
    }
}
