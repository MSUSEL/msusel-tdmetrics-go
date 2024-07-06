using YamlDotNet.RepresentationModel;

namespace Constructs.Data;

internal class Array(YamlSequenceNode source) : Node(source) {
    private readonly YamlSequenceNode source = source;

    public int Count => this.source.Children.Count;

    public Node this[int index] => new(this.source[index]);
}
