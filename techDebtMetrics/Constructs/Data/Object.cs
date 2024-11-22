using Constructs.Exceptions;
using System.Collections.Generic;
using System.Linq;
using YamlDotNet.RepresentationModel;

namespace Constructs.Data;

internal class Object(YamlMappingNode source) : Node(source) {
    private readonly YamlMappingNode source = source;

    public bool Contains(string key) => this.source.Children.ContainsKey(new YamlScalarNode(key));

    public Node this[string key] => new(this.source[new YamlScalarNode(key)]);

    public IEnumerable<KeyValuePair<string, Node>> Children => this.source.Children.
        Select(p => new KeyValuePair<string, Node>(new Node(p.Key).AsString(), new(p.Value)));

    public void PreallocateList<T>(string name, List<T> list)
        where T : new() {
        if (this.Contains(name)) {
            int count = this[name].AsArray().Count;
            for (int i = 0; i < count; ++i)
                list.Add(new T());
        }
    }

    public void InitializeList<T>(Project project, string name, List<T> list)
        where T : IInitializable {
        if (this.Contains(name)) {
            Array listArr = this[name].AsArray();
            for (int i = 0; i < listArr.Count; ++i)
                list[i].Initialize(project, listArr[i]);
        }
    }

    private Node readValue(string name) => this[name] ?? throw new MissingDataException(name);
    
    public string ReadString(string name) => this.readValue(name).AsString();

    public bool ReadBool(string name) => this.readValue(name).AsBool();

    public int ReadInt(string name) => this.readValue(name).AsInt();
    
    public uint ReadUint(string name) => this.readValue(name).AsUint();

    public T ReadIndexType<T>(string name, Project project)
        where T : ITypeDesc =>
        project.GetTypeAtIndex<T>(this.ReadUint(name));

    public void ReadIndexTypeList<T>(string name, Project project, List<T> list)
        where T : ITypeDesc {
        if (this.Contains(name)) {
            Array exactArr = this[name].AsArray();
            for (int i = 0; i < exactArr.Count; ++i) {
                uint typeIndex = exactArr[i].AsUint();
                list.Add(project.GetTypeAtIndex<T>(typeIndex));
            }
        }
    }

    public T ReadKey<T>(string name, Project project)
        where T : ITypeDesc =>
        // TODO: FIX
        project.GetTypeAtIndex<T>(this.ReadUint(name));
}
