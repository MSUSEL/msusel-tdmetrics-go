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
        this.tryReadNode(name)?.AsArray().PreallocateList(list);

    public void InitializeList<T>(Project project, string name, List<T> list)
        where T : IInitializable =>
        this.tryReadNode(name)?.AsArray().InitializeList<T>(project, list);

    #region Must Reads

    private Node readNode(string name) =>
        this.source.Children.TryGetValue(new YamlScalarNode(name), out YamlNode? value) ?
            new Node(value) :
            throw new MissingDataException("Expected "+name+" in object at "+this.source.Start.ToString());

    public string ReadString(string name) => this.readNode(name).AsString();

    public bool ReadBool(string name) => this.readNode(name).AsBool();

    public int ReadInt(string name) => this.readNode(name).AsInt();

    public T ReadIndex<T>(string name, IReadOnlyList<T> source)
        where T : IConstruct =>
        this.readNode(name).AsIndex(source);

    public List<T> ReadIndexList<T>(string name, IReadOnlyList<T> source)
        where T : IConstruct =>
        this.readNode(name).AsArray().AsIndexList(source);

    public T ReadKey<T>(string name, Project project)
        where T : IConstruct =>
        this.readNode(name).AsKey<T>(project);

    public List<T> ReadKeyList<T>(string name, Project project)
        where T : IConstruct =>
        this.readNode(name).AsArray().AsKeyList<T>(project);

    public Location ReadLocation(string name, Project project) =>
        this.readNode(name).AsLocation(project);

    #endregion
    #region Optional Reads

    private Node? tryReadNode(string name) =>
        this.source.Children.TryGetValue(new YamlScalarNode(name), out YamlNode? value) ?
            new Node(value) : null;

    public string TryReadString(string name) => this.tryReadNode(name)?.AsString() ?? "";

    public bool TryReadBool(string name) => this.tryReadNode(name)?.AsBool() ?? false;

    public int TryReadInt(string name) => this.tryReadNode(name)?.AsInt() ?? 0;

    public T? TryReadIndex<T>(string name, IReadOnlyList<T> source)
        where T : class, IConstruct =>
        this.tryReadNode(name)?.AsIndex(source);

    public List<T> TryReadIndexList<T>(string name, IReadOnlyList<T> source)
        where T : IConstruct =>
        this.tryReadNode(name)?.AsArray()?.AsIndexList(source) ?? [];

    public T? TryReadKey<T>(string name, Project project)
        where T : class, IConstruct =>
        this.tryReadNode(name)?.AsKey<T>(project);

    public List<T> TryReadKeyList<T>(string name, Project project)
        where T : IConstruct =>
        this.tryReadNode(name)?.AsArray()?.AsKeyList<T>(project) ?? [];

    public Location TryReadLocation(string name, Project project) =>
        this.tryReadNode(name)?.AsLocation(project) ?? Locations.Unknown;

    #endregion
}
