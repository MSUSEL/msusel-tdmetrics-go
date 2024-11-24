using Constructs.Tooling;
using System.Collections.Generic;
using System.IO;
using YamlDotNet.RepresentationModel;

namespace Constructs.Data;

internal class Node(YamlNode source) {
    private readonly YamlNode source = source;

    static public Node Parse(string text) {
        YamlStream yaml = [];
        yaml.Load(new StringReader(text));
        return new Node(yaml.Documents[0].RootNode);
    }

    public bool IsNull => this.source is null;

    public bool IsObject => this.source is YamlMappingNode;

    public bool IsArray => this.source is YamlSequenceNode;

    public bool IsScalar => this.source is YamlScalarNode;

    public Object AsObject() => new(this.source as YamlMappingNode ??
        throw new System.InvalidCastException("Not an object node at " + this.source.End));

    public Array AsArray() => new(this.source as YamlSequenceNode ??
        throw new System.InvalidCastException("Not an array node at " + this.source.End));

    private YamlScalarNode asValue() => this.source as YamlScalarNode ??
        throw new System.InvalidCastException("Not a value node at " + this.source.End);

    public string AsString() => this.asValue().Value ??
        throw new System.Exception("Null string from value node.");

    public bool AsBool() => bool.Parse(this.AsString());

    public int AsInt() => int.Parse(this.AsString());

    public Location AsLocation(Project project) =>
        project.Locations[this.AsInt()];

    public T AsIndex<T>(IReadOnlyList<T> source)
        where T : IConstruct =>
        source[this.AsInt()];

    public T AsKey<T>(Project project)
        where T : IConstruct =>
        (T)readKey(this.AsString(), project);

    static private int keySplitPoint(string key) {
        for (int i = key.Length - 1; i >= 0; --i {
            if (!char.IsDigit(key[i])) return i + 1;
        }
        throw new System.Exception("bad key: " + key);
    }

    static private IConstruct readKey(string key, Project project) {
        int split = keySplitPoint(key);
        string name = key[..split];
        int index = int.Parse(key[split..]);
        return name switch {
            "abstract" => project.Abstracts[index],
            "argument" => project.Arguments[index],
            "basic" => project.Basics[index],
            "field" => project.Fields[index],
            "interfaceDecl" => project.InterfaceDecls[index],
            "interfaceDesc" => project.InterfaceDescs[index],
            "interfaceInst" => project.InterfaceInsts[index],
            "method" => project.Methods[index],
            "methodInst" => project.MethodInsts[index],
            "metrics" => project.Metrics[index],
            "object" => project.ObjectDecls[index],
            "objectInst" => project.ObjectInsts[index],
            "package" => project.Packages[index],
            "selection" => project.Selections[index],
            "signature" => project.Signatures[index],
            "structDesc" => project.StructDescs[index],
            "typeParam" => project.TypeParams[index],
            "value" => project.Values[index],
            _ => throw new InvalidDataException(name)
        };
    }
}
 