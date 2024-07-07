using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Linq;

namespace Constructs;

public class Union : ITypeDesc, IInitializable {
    public IReadOnlyList<ITypeDesc> Exact => this.inExact.AsReadOnly();
    private readonly List<ITypeDesc> inExact = [];

    public IReadOnlyList<ITypeDesc> Approx => this.inApprox.AsReadOnly();
    private readonly List<ITypeDesc> inApprox = [];

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();
        obj.ReadIndexTypeList("extern", getter, this.inExact);
        obj.ReadIndexTypeList("approx", getter, this.inApprox);
    }

    public override string ToString() =>
        this.Exact.ToStrings().Concat(this.Approx.ToStrings(prefix: "~")).Join("|");

    public void ToStub(Journal j) {
        j.Write(this.Exact, separator: "|");
        string prefix = this.Exact.Count > 0 ? "|~" : "~";
        j.Write(this.Approx, prefix, "", "|~");
    }
}
