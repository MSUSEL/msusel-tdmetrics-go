using Constructs.Tooling;
using System;
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
        throw new InvalidCastException("Not an object node at " + this.source.End));

    public Array AsArray() => new(this.source as YamlSequenceNode ??
        throw new InvalidCastException("Not an array node at " + this.source.End));

    private YamlScalarNode asValue() => this.source as YamlScalarNode ??
        throw new InvalidCastException("Not a value node at " + this.source.End);

    public string AsString() => this.asValue().Value ??
        throw new Exception("Null string from value node.");

    public bool AsBool() => bool.Parse(this.AsString());

    public int AsInt() => int.Parse(this.AsString());

    public Location AsLocation(Project project) =>
        project.Locations[this.AsInt()];

    public T AsIndex<T>(IReadOnlyList<T> source)
        where T : IConstruct {
        int index = this.AsInt()-1;
        if (index < 0 || index >= source.Count)
            throw new Exception("Index out of range [0.." + source.Count + "): " + index);
        return source[index];
    }

    public T AsKey<T>(Project project)
        where T : IConstruct =>
        (T)readKey(this.AsString(), project);

    static private int keySplitPoint(string key) {
        for (int i = key.Length - 1; i >= 0; --i) {
            if (!char.IsDigit(key[i])) return i + 1;
        }
        throw new Exception("bad key: " + key);
    }

    static private IConstruct readKeyIndex<T>(string key, int index, IReadOnlyList<T> source)
        where T : IConstruct {
        if (index < 0 || index >= source.Count)
            throw new Exception("Key " + key + " out of range [0.." + source.Count + "): " + index);
        return source[index];
    }

    static private IConstruct readKey(string key, Project project) {
        int split = keySplitPoint(key);
        string name = key[..split];
        int index = int.Parse(key[split..])-1;
        return name switch {
            "abstract"      => readKeyIndex(key, index, project.Abstracts),
            "argument"      => readKeyIndex(key, index, project.Arguments),
            "basic"         => readKeyIndex(key, index, project.Basics),
            "field"         => readKeyIndex(key, index, project.Fields),
            "interfaceDecl" => readKeyIndex(key, index, project.InterfaceDecls),
            "interfaceDesc" => readKeyIndex(key, index, project.InterfaceDescs),
            "interfaceInst" => readKeyIndex(key, index, project.InterfaceInsts),
            "method"        => readKeyIndex(key, index, project.MethodDecls),
            "methodInst"    => readKeyIndex(key, index, project.MethodInsts),
            "metrics"       => readKeyIndex(key, index, project.Metrics),
            "object"        => readKeyIndex(key, index, project.ObjectDecls),
            "objectInst"    => readKeyIndex(key, index, project.ObjectInsts),
            "package"       => readKeyIndex(key, index, project.Packages),
            "selection"     => readKeyIndex(key, index, project.Selections),
            "signature"     => readKeyIndex(key, index, project.Signatures),
            "structDesc"    => readKeyIndex(key, index, project.StructDescs),
            "typeParam"     => readKeyIndex(key, index, project.TypeParams),
            "value"         => readKeyIndex(key, index, project.Values),
            _ => throw new InvalidDataException(name)
        };
    }
}
 