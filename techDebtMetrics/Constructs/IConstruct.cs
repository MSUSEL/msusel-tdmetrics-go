using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>The interface for all constructs.</summary>
/// <see cref="../../docs/genFeatureDef.md#constructs"/>
public interface IConstruct {

    /// <summary>Writes the construct to the given journal.</summary>
    /// <param name="j">The journal to write to.</param>
    public void ToStub(Journal j);

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; }

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    /// <remarks>
    /// This will not return packages, generics, receivers, or inheritance since
    /// those constructs are more superior/parents constructs.
    /// </remarks>
    public IEnumerable<IConstruct> SubConstructs { get; }
}
