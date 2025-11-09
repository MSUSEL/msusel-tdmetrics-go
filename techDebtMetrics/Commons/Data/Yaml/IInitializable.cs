namespace Commons.Data.Yaml;

/// <summary>A type that can be initialized while loading.</summary>
/// <typeparam name="D">The type of data to pass along when initializing.</typeparam>
public interface IInitializable<D> {

    /// <summary>Initializes this object with the data from the given node.</summary>
    /// <param name="data">The data to pass along when initializing.</param>
    /// <param name="index">The index of this object in the list being initalized.</param>
    /// <param name="node">The node to initialize with.</param>
    void Initialize(D data, int index, Node node);
}
