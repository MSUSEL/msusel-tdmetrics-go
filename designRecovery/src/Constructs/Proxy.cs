namespace designRecovery.src.Constructs;

internal class Proxy : ITypeDesc {
    public uint Index = 0;
    public override string ToString() => "@"+this.Index;
}
