using DesignRecovery.Extensions;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json.Nodes;

namespace DesignRecovery.Constructs;

public class Union : ITypeDesc, IInitializable {
    private readonly List<ITypeDesc> inExact = [];
    public IReadOnlyList<ITypeDesc> Exact => this.inExact.AsReadOnly();

    private readonly List<ITypeDesc> inApprox = [];
    public IReadOnlyList<ITypeDesc> Approx => this.inApprox.AsReadOnly();

    public void Initialize(TypeGetter getter, JsonNode node) {
        JsonObject obj = node.AsObject();
        obj.ReadIndexTypeList("extern", getter, this.inExact);
        obj.ReadIndexTypeList("approx", getter, this.inApprox);
    }

    public override string ToString() =>
        this.Exact.ToStrings().Concat(this.Approx.ToStrings(prefix: "~")).Join("|");
}
