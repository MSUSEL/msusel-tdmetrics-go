using Constructs.Tooling;
using System;
using System.Collections.Generic;
using System.IO;
using YamlDotNet.RepresentationModel;

namespace Constructs.Data;

/// <summary>This is a node of data to read with.</summary>
/// <param name="source">The underlying data source.</param>
internal class Node(YamlNode source) {
    private readonly YamlNode source = source;

    /// <summary>Reads a node and child nodes from the given input text.</summary>
    /// <param name="text">Text formatted as YAML or JSON to read from.</param>
    /// <returns>The read node.</returns>
    static public Node Parse(string text) {
        YamlStream yaml = [];
        yaml.Load(new StringReader(text));
        return new Node(yaml.Documents[0].RootNode);
    }

    /// <summary>Indicates this node is null.</summary>
    public bool IsNull => this.source is null;

    /// <summary>Indicates this node contains key/value pairs.</summary>
    public bool IsObject => this.source is YamlMappingNode;

    /// <summary>Indicates this node contains a list of nodes.</summary>
    public bool IsArray => this.source is YamlSequenceNode;

    /// <summary>Indicates this node is a single value, e.g. int, string.</summary>
    public bool IsScalar => this.source is YamlScalarNode;

    /// <summary>Gets this node as an object with key/value pairs.</summary>
    /// <returns>The object node.</returns>
    public Object AsObject() => new(this.source as YamlMappingNode ??
        throw new InvalidCastException("Not an object node at " + this.source.End));

    /// <summary>Gets this node as an array with a list of nodes.</summary>
    /// <returns>The array node.</returns>
    public Array AsArray() => new(this.source as YamlSequenceNode ??
        throw new InvalidCastException("Not an array node at " + this.source.End));

    /// <summary>Gets the underlying source as a scalar node.</summary>
    /// <returns>The scalar node source to read from.</returns>
    private YamlScalarNode getScalar() => this.source as YamlScalarNode ??
        throw new InvalidCastException("Not a value node at " + this.source.End);

    /// <summary>Gets this node as a string.</summary>
    /// <returns>The string value of this node.</returns>
    public string AsString() => this.getScalar().Value ??
        throw new Exception("Null string from value node.");

    /// <summary>Gets this node as a boolean.</summary>
    /// <returns>The boolean value of this node.</returns>
    public bool AsBool() => bool.Parse(this.AsString());

    /// <summary>Gets this node as a integer.</summary>
    /// <returns>The integer value of this node.</returns>
    public int AsInt() => int.Parse(this.AsString());

    /// <summary>Gets this node as a file location.</summary>
    /// <param name="project">The project to resolve the file location with.</param>
    /// <returns>The read file location.</returns>
    public Location AsLocation(Project project) =>
        project.Locations[this.AsInt()];

    /// <summary>Gets this node as an index lookup from the given source.</summary>
    /// <typeparam name="T">The type of item to lookup.</typeparam>
    /// <param name="source">The source to index into.</param>
    /// <returns>The value read from the given index.</returns>
    public T AsIndex<T>(IReadOnlyList<T> source) {
        int index = this.AsInt() - 1;
        if (index < 0 || index >= source.Count)
            throw new Exception("Index out of range [0.." + source.Count + "): " + index);
        return source[index];
    }

    /// <summary>Gets this node as a key lookup from the given project.</summary>
    /// <see cref="../../../docs/genFeatureDef.md#keys"/>
    /// <typeparam name="T">The type of item to lookup.</typeparam>
    /// <param name="project">The project to read from.</param>
    /// <returns>The value read with the given key.</returns>
    public T AsKey<T>(Project project) =>
        (T)readKey(this.AsString(), project);

    /// <summary>Find the place in the key to split it into a name and index.</summary>
    /// <param name="key">The key to find the split point in.</param>
    /// <returns>The index in the key string of the first digit of the index part of the key.</returns>
    static private int keySplitPoint(string key) {
        for (int i = key.Length - 1; i >= 0; --i) {
            if (!char.IsDigit(key[i])) return i + 1;
        }
        throw new Exception("bad key: " + key);
    }

    /// <summary>Reads the given index from the given source as part of reading the given key.</summary>
    /// <typeparam name="T">The type of the list to read from.</typeparam>
    /// <param name="key">The key that is being processed.</param>
    /// <param name="index">The index from the key used to read a value from the given list.</param>
    /// <param name="source">The list get an item at the given index from.</param>
    /// <returns>The item from the given list at the given index.</returns>
    static private T readKeyIndex<T>(string key, int index, IReadOnlyList<T> source) {
        if (index < 0 || index >= source.Count)
            throw new Exception("Key " + key + " out of range [0.." + source.Count + "): " + index);
        return source[index];
    }

    /// <summary>Reads a single key from the given project.</summary>
    /// <see cref="../../../docs/genFeatureDef.md#keys"/>
    /// <param name="key">The key of the value to read.</param>
    /// <param name="project">The project to read a key from.</param>
    /// <returns>The read key from the project.</returns>
    static private IConstruct readKey(string key, Project project) {
        int split = keySplitPoint(key);
        string name = key[..split];
        int index = int.Parse(key[split..]) - 1;
        return name switch {
            "abstract" => readKeyIndex(key, index, project.Abstracts),
            "argument" => readKeyIndex(key, index, project.Arguments),
            "basic" => readKeyIndex(key, index, project.Basics),
            "field" => readKeyIndex(key, index, project.Fields),
            "interfaceDecl" => readKeyIndex(key, index, project.InterfaceDecls),
            "interfaceDesc" => readKeyIndex(key, index, project.InterfaceDescs),
            "interfaceInst" => readKeyIndex(key, index, project.InterfaceInsts),
            "method" => readKeyIndex(key, index, project.MethodDecls),
            "methodInst" => readKeyIndex(key, index, project.MethodInsts),
            "metrics" => readKeyIndex(key, index, project.Metrics),
            "object" => readKeyIndex(key, index, project.ObjectDecls),
            "objectInst" => readKeyIndex(key, index, project.ObjectInsts),
            "package" => readKeyIndex(key, index, project.Packages),
            "selection" => readKeyIndex(key, index, project.Selections),
            "signature" => readKeyIndex(key, index, project.Signatures),
            "structDesc" => readKeyIndex(key, index, project.StructDescs),
            "typeParam" => readKeyIndex(key, index, project.TypeParams),
            "value" => readKeyIndex(key, index, project.Values),
            _ => throw new InvalidDataException(name)
        };
    }
}
