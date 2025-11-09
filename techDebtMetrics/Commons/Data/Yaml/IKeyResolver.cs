namespace Commons.Data.Yaml;

/// <summary>The interface used for resolving key values.</summary>
public interface IKeyResolver {

    /// <summary>Finds data inside this data by the given name and index from a key.</summary>
    /// <param name="name">The name of the data to look up.</param>
    /// <param name="index">The index into the named data to get.</param>
    /// <returns>The data found with that name and index.</returns>
    object FindData(string name, int index);
}
