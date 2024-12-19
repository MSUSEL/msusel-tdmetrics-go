namespace Constructs;

/// <summary>The interface for both method declarations and method instances.</summary>
public interface IMethod : IConstruct {

    /// <summary>The name of the method.</summary>
    public string Name { get; }

    /// <summary>The signature type declaration.</summary>
    public Signature Signature { get; }
}
