namespace Constructs;

/// <summary>The interface for both interface declarations and interface instances.</summary>
public interface IInterface : ITypeDesc {

    /// <summary>The name of the interface.</summary>
    public string Name { get; }

    /// <summary>The interface type description.</summary>
    public InterfaceDesc Interface { get; }
}
