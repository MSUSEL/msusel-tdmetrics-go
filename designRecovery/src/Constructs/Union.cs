using designRecovery.src.Extensions;

namespace designRecovery.src.Constructs;

internal class Union {
    public readonly List<ITypeDesc> Exact = new();
    public readonly List<ITypeDesc> Approx = new();
    public override string ToString() => this.Exact.ToStrings().Concat(this.Approx.ToStrings(prefix: "~")).Join("|");
}
