using Constructs.Tooling;

namespace Constructs;

/// <summary>The interface for all constructs.</summary>
/// <see cref="../../docs/genFeatureDef.md#constructs"/>
public interface IConstruct {

    /// <summary>Writes the construct to the given journal.</summary>
    /// <param name="j">The journal to write to.</param>
    public void ToStub(Journal j);

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; }
}
