using Commons.Data.Locations;
using System;
using System.Collections.Generic;
using System.Linq;
using YamlDotNet.RepresentationModel;

namespace Commons.Data.Yaml;

/// <summary>An object with key/value pairs to read data from.</summary>
/// <param name="source">The underlying data source.</param>
public class Object(YamlMappingNode source) : Node(source) {
    private readonly YamlMappingNode source = source;

    /// <summary>Creates a new empty object.</summary>
    public Object() : this([]) { }

    /// <summary>The number of key/value pairs in this node.</summary>
    public int Count => this.source.Children.Count;

    /// <summary>Determines if this node contains the given name.</summary>
    /// <param name="name">The name of the key in a key/value pair to check for.</param>
    /// <returns>True if the given key name in a key/value pair exists.</returns>
    public bool Contains(string name) => this.source.Children.ContainsKey(new YamlScalarNode(name));

    /// <summary>Enumerates all the key/value pairs in this node.</summary>
    public IEnumerable<KeyValuePair<Node, Node>> Children => this.source.Children.
        Select(p => new KeyValuePair<Node, Node>(new Node(p.Key), new(p.Value)));

    /// <summary>Adds a new item for the node at the given key into the given list.</summary>
    /// <typeparam name="T">The type to preallocate.</typeparam
    /// <param name="name">The name for the node to preallocate for.</param>
    /// <param name="list">The list to add to.</param>
    public void PreallocateList<T>(string name, List<T> list)
        where T : new() =>
        this.TryReadNode(name)?.AsArray().PreallocateList(list);

    /// <summary>Adds a new item for the node at the given key into the given list.</summary>
    /// <typeparam name="T">The type to preallocate.</typeparam
    /// <param name="name">The name for the node to preallocate for.</param>
    /// <param name="list">The list to add to.</param>
    /// <param name="constructor">The construtor used for preallocating the list.</param>
    public void PreallocateList<T>(string name, List<T> list, Func<Node, T> constructor) =>
        this.TryReadNode(name)?.AsArray().PreallocateList(list, constructor);

    /// <summary>Initializes the given preallocated list with the node at the given name.</summary>
    /// <typeparam name="T">The type to call Initialize on.</typeparam>
    /// <typeparam name="D">The type of data to pass along while initializing.</typeparam>
    /// <param name="data">The data to to pass along while initializing.</param>
    /// <param name="name">The name for the node to initalize with.</param>
    /// <param name="list">The list of items to initialize.</param>
    public void InitializeList<T, D>(D data, string name, List<T> list)
        where T : IInitializable<D> =>
        this.TryReadNode(name)?.AsArray().InitializeList(data, list);

    #region Writes

    /// <summary>Adds a new key value pair to this object.</summary>
    /// <param name="key">The key to set the value to.</param>
    /// <param name="value">The value to set at the key.</param>
    public void Add(Node key, Node value) => this.source.Add(key.Source, value.Source);
    
    /// <summary>Adds a new key value pair to this object.</summary>
    /// <param name="key">The key to set the value to.</param>
    /// <param name="value">The value to set at the key.</param>
    public void Add(string key, Node value) => this.Add(new Node(key), value);
    
    /// <summary>Adds a new key value pair to this object.</summary>
    /// <param name="key">The key to set the value to.</param>
    /// <param name="value">The string value to set at the key.</param>
    public void Add(string key, string value) => this.Add(key, new Node(value));
    
    /// <summary>Adds a new key value pair to this object.</summary>
    /// <param name="key">The key to set the value to.</param>
    /// <param name="value">The bool value to set at the key.</param>
    public void Add(string key, bool value) => this.Add(key, new Node(value));
    
    /// <summary>Adds a new key value pair to this object.</summary>
    /// <param name="key">The key to set the value to.</param>
    /// <param name="value">The int value to set at the key.</param>
    public void Add(string key, int value) => this.Add(key, new Node(value));
    
    /// <summary>Adds a new key value pair to this object.</summary>
    /// <param name="key">The key to set the value to.</param>
    /// <param name="value">The double value to set at the key.</param>
    public void Add(string key, double value) => this.Add(key, new Node(value));
    
    /// <summary>Adds a new key location pair to this object.</summary>
    /// <param name="key">The key to set the location to.</param>
    /// <param name="locs">The location writer to use.</param>
    /// <param name="loc">The location to set at the key.</param>
    public void Add(string key, Writer locs, Location loc) => this.Add(key, new Node(locs, loc));

    #endregion
    #region Must Reads

    /// <summary>Gets the node with the given name that must exist.</summary>
    /// <param name="name">The name of the node to get.</param>
    /// <returns>The node with the given name.</returns>
    public Node ReadNode(string name) =>
        this.source.Children.TryGetValue(new YamlScalarNode(name), out YamlNode? value) ?
            new Node(value) :
            throw new Exception("Expected \"" + name + "\" in object at " + this.source.Start.ToString());

    /// <summary>Gets the node with the given name as a string.</summary>
    /// <param name="name">The name of the node to get the string from.</param>
    /// <returns>The string from the node with the given name.</returns>
    public string ReadString(string name) => this.ReadNode(name).AsString();

    /// <summary>Gets the node with the given name as a boolean.</summary>
    /// <param name="name">The name of the node to get the boolean from.</param>
    /// <returns>The boolean from the node with the given name.</returns>
    public bool ReadBool(string name) => this.ReadNode(name).AsBool();

    /// <summary>Gets the node with the given name as an integer.</summary>
    /// <param name="name">The name of the node to get the integer from.</param>
    /// <returns>The integer from the node with the given name.</returns>
    public int ReadInt(string name) => this.ReadNode(name).AsInt();

    /// <summary>Gets the node with the given name as a double.</summary>
    /// <param name="name">The name of the node to get the double from.</param>
    /// <returns>The double from the node with the given name.</returns>
    public double ReadDouble(string name) => this.ReadNode(name).AsDouble();

    /// <summary>Gets the item at the index from the node with the given name.</summary>
    /// <typeparam name="T">The type of item to get from the list.</typeparam>
    /// <param name="name">The name of the node to get the index from.</param>
    /// <param name="source">The list to read the value from.</param>
    /// <returns>The value at the index in the given list.</returns>
    public T ReadIndex<T>(string name, IReadOnlyList<T> source)
        where T : class =>
        this.ReadNode(name).AsIndex(source);

    /// <summary>Gets a list of values via an array of indices in the node with the given name.</summary>
    /// <typeparam name="T">The type of values to get.</typeparam>
    /// <param name="name">The name of the node containing a list of indices.</param>
    /// <param name="dest">The list to write all the read values into.</param>
    /// <param name="source">The list to read all the indices from.</param>
    public void ReadIndexList<T>(string name, List<T> dest, IReadOnlyList<T> source) =>
        this.ReadNode(name).AsArray().AsIndexList(dest, source);

    /// <summary>Gets the item with the given key from the node with the given name.</summary>
    /// <typeparam name="T">The type of item to get from the lookup.</typeparam>
    /// <param name="res">The key resolver to lookup the value with the key.</param>
    /// <param name="name">The name of the node to get the key from.</param>
    /// <returns>The value read from the key in the node with the given name.</returns>
    public T ReadKey<T>(IKeyResolver res, string name)
        where T : class =>
        this.ReadNode(name).AsKey<T>(res);

    /// <summary>Gets a list of values via an array of keys in the node with the given name.</summary>
    /// <typeparam name="T">The type of values to get.</typeparam>
    /// <param name="res">The key resolver to lookup the value with the key.</param>
    /// <param name="name">The name of the node containing a list of keys.</param>
    /// <param name="dest">The list to write all the read values into.</param>
    public void ReadKeyList<T>(IKeyResolver res, string name, List<T> dest) =>
        this.ReadNode(name).AsArray().AsKeyList(res, dest);

    /// <summary>Gets path and line number location from the node with the given name.</summary>
    /// <param name="locs">The locations to read the location from.</param>
    /// <param name="name">The name of the node to get the location from.</param>
    /// <returns>The location read from the project with the node at the given name.</returns>
    public Location ReadLocation(Reader locs, string name) =>
        this.ReadNode(name).AsLocation(locs);

    #endregion
    #region Optional Reads

    /// <summary>Tries to get the node with the given name that may not exist.</summary>
    /// <param name="name">The name of the node to get.</param>
    /// <returns>The node with the given name or null if it didn't exist.</returns>
    public Node? TryReadNode(string name) =>
        this.source.Children.TryGetValue(new YamlScalarNode(name), out YamlNode? value) ?
            new Node(value) : null;

    /// <summary>Try to get the node with the given name as a string.</summary>
    /// <param name="name">The name of the node to get the string from.</param>
    /// <param name="defaultValue">The default value to return if no node exists by the given name.</param>
    /// <returns>The string from the node with the given name or the default value if it didn't exist.</returns>
    public string TryReadString(string name, string defaultValue = "") =>
        this.TryReadNode(name)?.AsString() ?? defaultValue;

    /// <summary>Tries to get the node with the given name as a boolean.</summary>
    /// <param name="name">The name of the node to get the boolean from.</param>
    /// <param name="defaultValue">The default value to return if no node exists by the given name.</param>
    /// <returns>The boolean from the node with the given name or the default value if it didn't exist.</returns>
    public bool TryReadBool(string name, bool defaultValue = false) =>
        this.TryReadNode(name)?.AsBool() ?? defaultValue;

    /// <summary>Tries to get the node with the given name as an integer.</summary>
    /// <param name="name">The name of the node to get the integer from.</param>
    /// <param name="defaultValue">The default value to return if no node exists by the given name.</param>
    /// <returns>The integer from the node with the given name or the default value if it didn't exist.</returns>
    public int TryReadInt(string name, int defaultValue = 0) =>
        this.TryReadNode(name)?.AsInt() ?? defaultValue;

    /// <summary>Tries to get the node with the given name as a double.</summary>
    /// <param name="name">The name of the node to get the double from.</param>
    /// <param name="defaultValue">The default value to return if no node exists by the given name.</param>
    /// <returns>The double from the node with the given name or the default value if it didn't exist.</returns>
    public double TryReadDouble(string name, double defaultValue = 0.0) =>
        this.TryReadNode(name)?.AsDouble() ?? defaultValue;

    /// <summary>Tries to get the item at the index from the node with the given name.</summary>
    /// <typeparam name="T">The type of item to get from the list.</typeparam>
    /// <param name="name">The name of the node to get the index from.</param>
    /// <param name="source">The list to read the value from.</param>
    /// <returns>The value at the index in the given list or null if it didn't exist.</returns>
    public T? TryReadIndex<T>(string name, IReadOnlyList<T> source)
        where T : class =>
        this.TryReadNode(name)?.AsIndex(source);

    /// <summary>Tries to get a list of values via an array of indices in the node with the given name.</summary>
    /// <typeparam name="T">The type of values to get.</typeparam>
    /// <param name="name">The name of the node containing a list of indices.</param>
    /// <param name="dest">The list to write all the read values into.</param>
    /// <param name="source">The list to read all the indices from.</param>
    public void TryReadIndexList<T>(string name, List<T> dest, IReadOnlyList<T> source) =>
        this.TryReadNode(name)?.AsArray()?.AsIndexList(dest, source);

    /// <summary>Tries to get the item with the given key from the node with the given name.</summary>
    /// <typeparam name="T">The type of item to get from the lookup.</typeparam>
    /// <param name="res">The key resolver to lookup the value with the key.</param>
    /// <param name="name">The name of the node to get the key from.</param>
    /// <returns>The value read from the key in the node with the given name or null if it didn't exist.</returns>
    public T? TryReadKey<T>(IKeyResolver res, string name)
        where T : class =>
        this.TryReadNode(name)?.AsKey<T>(res);

    /// <summary>Tries to get a list of values via an array of keys in the node with the given name.</summary>
    /// <typeparam name="T">The type of values to get.</typeparam>
    /// <param name="res">The key resolver to lookup the value with the key.</param>
    /// <param name="name">The name of the node containing a list of keys.</param>
    /// <param name="dest">The list to write all the read values into.</param>
    public void TryReadKeyList<T>(IKeyResolver res, string name, List<T> dest) =>
        this.TryReadNode(name)?.AsArray()?.AsKeyList(res, dest);

    /// <summary>Tries to get path and line number location from the node with the given name.</summary>
    /// <param name="name">The name of the node to get the location from.</param>
    /// <param name="locs">The locations to read the location from.</param>
    /// <returns>The location read from the project with the node at the given name or the unknown location if it didn't exist..</returns>
    public Location TryReadLocation(Reader locs, string name) =>
        this.TryReadNode(name)?.AsLocation(locs) ?? Location.Unknown;

    #endregion
}
