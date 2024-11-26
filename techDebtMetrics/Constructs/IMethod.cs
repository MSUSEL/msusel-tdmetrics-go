namespace Constructs;

public interface IMethod : IConstruct {
    public string Name { get; }
    public Signature Signature { get; }
}
