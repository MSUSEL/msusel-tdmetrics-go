using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json.Nodes;

namespace Constructs;

public class Union : ITypeDesc, IInitializable {
    public IReadOnlyList<ITypeDesc> Exact => this.inExact.AsReadOnly();
    private readonly List<ITypeDesc> inExact = [];

    public IReadOnlyList<ITypeDesc> Approx => this.inApprox.AsReadOnly();
    private readonly List<ITypeDesc> inApprox = [];

    void IInitializable.Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        obj.ReadIndexTypeList("extern", getter, this.inExact);
        obj.ReadIndexTypeList("approx", getter, this.inApprox);
    }

    public override string ToString() => this.ToStub();
 
    public string ToStub() =>
        this.Exact.ToStrings().Concat(this.Approx.ToStrings(prefix: "~")).Join("|");
}
