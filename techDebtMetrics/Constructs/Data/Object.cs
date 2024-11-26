using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Linq;
using YamlDotNet.RepresentationModel;

namespace Constructs.Data;

internal class Object(YamlMappingNode source) : Node(source) {
    private readonly YamlMappingNode source = source;

    public int Count => this.source.Children.Count;

    public bool Contains(string name) => this.source.Children.ContainsKey(new YamlScalarNode(name));

    public IEnumerable<KeyValuePair<string, Node>> Children => this.source.Children.
        Select(p => new KeyValuePair<string, Node>(new Node(p.Key).AsString(), new(p.Value)));

    public void PreallocateList<T>(string name, List<T> list)
        where T : new() =>
        this.TryReadNode(name)?.AsArray().PreallocateList(list);

    public void InitializeList<T>(Project project, string name, List<T> list)
        where T : IInitializable =>
        this.TryReadNode(name)?.AsArray().InitializeList<T>(project, list);

    #region Must Reads

    public Node ReadNode(string name) =>
        this.source.Children.TryGetValue(new YamlScalarNode(name), out YamlNode? value) ?
            new Node(value) :
            throw new MissingDataException("expected \"" + name + "\" in object at " + this.source.Start.ToString());

    public string ReadString(string name) => this.ReadNode(name).AsString();

    public bool ReadBool(string name) => this.ReadNode(name).AsBool();

    public int ReadInt(string name) => this.ReadNode(name).AsInt();

    public T ReadIndex<T>(string name, IReadOnlyList<T> source)
        where T : IConstruct =>
        this.ReadNode(name).AsIndex(source);

    public List<T> ReadIndexList<T>(string name, IReadOnlyList<T> source)
        where T : IConstruct =>
        this.ReadNode(name).AsArray().AsIndexList(source);

    public T ReadKey<T>(string name, Project project)
        where T : IConstruct =>
        this.ReadNode(name).AsKey<T>(project);

    public List<T> ReadKeyList<T>(string name, Project project)
        where T : IConstruct =>
        this.ReadNode(name).AsArray().AsKeyList<T>(project);

    public Location ReadLocation(string name, Project project) =>
        this.ReadNode(name).AsLocation(project);

    #endregion
    #region Optional Reads

    public Node? TryReadNode(string name) =>
        this.source.Children.TryGetValue(new YamlScalarNode(name), out YamlNode? value) ?
            new Node(value) : null;

    public string TryReadString(string name) => this.TryReadNode(name)?.AsString() ?? "";

    public bool TryReadBool(string name) => this.TryReadNode(name)?.AsBool() ?? false;

    public int TryReadInt(string name) => this.TryReadNode(name)?.AsInt() ?? 0;

    public T? TryReadIndex<T>(string name, IReadOnlyList<T> source)
        where T : class, IConstruct =>
        this.TryReadNode(name)?.AsIndex(source);

    public List<T> TryReadIndexList<T>(string name, IReadOnlyList<T> source)
        where T : IConstruct =>
        this.TryReadNode(name)?.AsArray()?.AsIndexList(source) ?? [];

    public T? TryReadKey<T>(string name, Project project)
        where T : class, IConstruct =>
        this.TryReadNode(name)?.AsKey<T>(project);

    public List<T> TryReadKeyList<T>(string name, Project project)
        where T : IConstruct =>
        this.TryReadNode(name)?.AsArray()?.AsKeyList<T>(project) ?? [];

    public Location TryReadLocation(string name, Project project) =>
        this.TryReadNode(name)?.AsLocation(project) ?? Locations.Unknown;

    #endregion
}
