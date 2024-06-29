using designRecovery.src.Extensions;
using System.Text.Json.Nodes;

namespace designRecovery.src.Constructs;

internal class Union: ITypeDesc {
    private readonly List<ITypeDesc> inExact = [];
    public IReadOnlyList<ITypeDesc> Exact => this.inExact.AsReadOnly();

    private readonly List<ITypeDesc> inApprox = [];
    public IReadOnlyList<ITypeDesc> Approx => this.inApprox.AsReadOnly();

    public void Initialize(ITypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        obj.ReadIndexTypeList("extern", getter, this.inExact);
        obj.ReadIndexTypeList("approx", getter, this.inApprox);
    }

    public override string ToString() =>
        this.Exact.ToStrings().Concat(this.Approx.ToStrings(prefix: "~")).Join("|");
}
