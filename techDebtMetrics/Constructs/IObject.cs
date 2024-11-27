namespace Constructs;

/// <summary>The interface for an object declaration or object instantiation.</summary>
public interface IObject : ITypeDesc {

    /// <summary>The name of the object.</summary>
    public string Name { get; }

    /// <summary>The interface that summarizes the methods of this object.</summary>
    public InterfaceDesc Interface { get; }

    /// <summary>The data contated by this object.</summary>
    public StructDesc Data { get; }
}
