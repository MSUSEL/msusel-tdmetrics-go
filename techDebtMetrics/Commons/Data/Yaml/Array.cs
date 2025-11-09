using Commons.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;
using YamlDotNet.RepresentationModel;

namespace Commons.Data.Yaml;

/// <summary>A node containing a fixed length linear list of nodes.</summary>
/// <param name="source">The underlying data source.</param>
public class Array(YamlSequenceNode source) : Node(source) {
    private readonly YamlSequenceNode source = source;

    /// <summary>Creates a new empty array.</summary>
    public Array() : this(new YamlSequenceNode()) { }

    /// <summary>Creates a new array with the given nodes.</summary>
    /// <param name="nodes">The nodes to add to the array.</param>
    public Array(params Node[] nodes) : this() {
        foreach (Node n in nodes) this.Add(n);
    }

    /// <summary>Creates a new array with the given nodes.</summary>
    /// <param name="nodes">The nodes to add to the array.</param>
    public Array(IEnumerable<Node> nodes) : this() => nodes.ForAll(this.Add);

    /// <summary>The number of nodes in this node.</summary>
    public int Count => this.source.Children.Count;

    /// <summary>Gets the node at the given index.</summary>
    /// <param name="index">The index to get the node from. Must be [0..Count).</param>
    /// <returns>The node at the given index.</returns>
    public Node this[int index] => new(this.source[index]);

    /// <summary>Enumerates all the nodes in this node.</summary>
    public IEnumerable<Node> Items => this.source.Children.Select(x => new Node(x));

    /// <summary>Adds a new node to this sequence.</summary>
    /// <param name="value">The node to add to this sequence.</param>
    /// <returns>The node that was passed in.</returns>
    public Node Add(Node value) {
        this.source.Add(value.Source);
        return value;
    }

    /// <summary>Adds a new node with the given value.</summary>
    /// <param name="value">The value for the new node.</param>
    /// <returns>The new node that was created and added.</returns>
    public Node Add(string value) => this.Add(new Node(value));

    /// <summary>Adds a new node with the given value.</summary>
    /// <param name="value">The value for the new node.</param>
    /// <returns>The new node that was created and added.</returns>
    public Node Add(bool value) => this.Add(new Node(value));

    /// <summary>Adds a new node with the given value.</summary>
    /// <param name="value">The value for the new node.</param>
    /// <returns>The new node that was created and added.</returns>
    public Node Add(int value) => this.Add(new Node(value));

    /// <summary>Adds a new node with the given value.</summary>
    /// <param name="value">The value for the new node.</param>
    /// <returns>The new node that was created and added.</returns>
    public Node Add(double value) => this.Add(new Node(value));

    /// <summary>Adds a new item for each node in this array into the given list.</summary>
    /// <typeparam name="T">The type to preallocate.</typeparam>
    /// <param name="list">The list to add to.</param>
    public void PreallocateList<T>(List<T> list) where T : new() =>
        this.PreallocateList(list, _ => new T());

    /// <summary>Adds a new item for each node in this array into the given list.</summary>
    /// <typeparam name="T">The type to preallocate.</typeparam>
    /// <param name="list">The list to add to.</param>
    /// <param name="constructor">The construtor used for preallocating the list.</param>
    public void PreallocateList<T>(List<T> list, Func<Node, T> constructor) {
        for (int i = 0; i < this.Count; ++i)
            list.Add(constructor(this[i]));
    }

    /// <summary>Initializes the given preallocated list with the nodes in this node.</summary>
    /// <typeparam name="T">The type to call Initialize on.</typeparam>
    /// <typeparam name="D">The type of data to pass along while initializing.</typeparam>
    /// <param name="data">The data to to pass along while initializing.</param>
    /// <param name="list">The list of items to initialize.</param>
    public void InitializeList<T, D>(D data, IReadOnlyList<T> list)
        where T : IInitializable<D> {
        for (int i = 0; i < this.Count; ++i) {
            try {
                list[i].Initialize(data, i, this[i]);
            } catch (Exception ex) {
                throw new Exception("Failed to initialize #" + i + " in " + typeof(T).Name + " list:", ex);
            }
        }
    }

    /// <summary>Reads this node as a list of indices into the given source.</summary>
    /// <typeparam name="T">The type of items to read.</typeparam>
    /// <param name="dest">The list to add items read via indices in this node to.</param>
    /// <param name="source">The source list to index.</param>
    public void AsIndexList<T>(List<T> dest, IReadOnlyList<T> source) =>
        dest.AddRange(this.Items.Select(item => item.AsIndex(source)));

    /// <summary>Reads this node as a list of keys into the lookup.</summary>
    /// <see cref="docs/genFeatureDef.md#keys"/>
    /// <typeparam name="T">The type of items to read.</typeparam>
    /// <param name="res">The key resolver to lookup the value with the key.</param>
    /// <param name="dest">The list to add items read via keys from the lookup to.</param>
    public void AsKeyList<T>(IKeyResolver res, List<T> dest) =>
        dest.AddRange(this.Items.Select(item => item.AsKey<T>(res)));
}
