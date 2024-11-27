using Constructs.Tooling;

namespace Constructs;

/// <summary>The interface to anything that is a named declaration.</summary>
/// <see cref="../../docs/genFeatureDef.md#declarations"/>
public interface IDeclaration : IConstruct {

    /// <summary>The name of the declaration.</summary>
    public string Name { get; }

    /// <summary>The location that the declaration was defined in the source code.</summary>
    public Location Location { get; }

    /// <summary>The package the declaration was defined in.</summary>
    public Package Package { get; }
}
