namespace Commons.Data.Yaml;

/// <summary>This is an interface for an object that can be turned into a node.</summary>
/// <typeparam name="D">The type of data to pass along when creating nodes.</typeparam>
public interface INodeable<D> {

    /// <summary>Creates a node for this object.</summary>
    /// <param name="data">The data to create the node with.</param>
    /// <param name="index">The index for this object in a list of these objects.</param>
    /// <returns>The created node for this object.</returns>
    Node ToNode(D data, int index);
}
