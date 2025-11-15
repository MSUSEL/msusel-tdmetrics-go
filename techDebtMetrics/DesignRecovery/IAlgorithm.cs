namespace DesignRecovery;

/// <summary>Defines an algorithm for converting constructs into technical debt data.</summary>
public interface IAlgorithm {

    /// <summary>This runs the algorithm to perform the conversion.</summary>
    /// <param name="man">The data manager containing the input and output project.</param>
    void Run(Manager man);
}
