using System.IO;
using YamlDotNet.RepresentationModel;

namespace Constructs.Data;

internal class Node(YamlNode source) {
    private readonly YamlNode source = source;

    static public Node Parse(string text) {
        YamlStream yaml = new();
        yaml.Load(new StringReader(text));
        return new Node(yaml.Documents[0].RootNode);
    }

    public Object AsObject() => new(this.source as YamlMappingNode ??
        throw new System.InvalidCastException("Not an object node at " + this.source.End));

    public Array AsArray() => new(this.source as YamlSequenceNode ??
        throw new System.InvalidCastException("Not an array node at " + this.source.End));

    private YamlScalarNode asValue() => this.source as YamlScalarNode ??
        throw new System.InvalidCastException("Not a value node at " + this.source.End);

    public string AsString() => this.asValue().Value ??
        throw new System.Exception("Null string from node value.");

    public bool AsBool() => bool.Parse(this.AsString());

    public int AsInt() => int.Parse(this.AsString());
    
    public uint AsUint() => uint.Parse(this.AsString());
}
