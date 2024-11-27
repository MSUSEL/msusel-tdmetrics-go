namespace Constructs.Data;

/// <summary>A type that can be initialized while loading.</summary>
internal interface IInitializable {

    /// <summary>Initializes this object with the data from the given node.</summary>
    /// <param name="project">The project to get additional information from while initializing.</param>
    /// <param name="node">The node to initialize with.</param>
    void Initialize(Project project, Node node);
}
