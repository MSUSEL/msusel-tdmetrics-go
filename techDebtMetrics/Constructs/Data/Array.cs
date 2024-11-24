using System.Collections.Generic;
using System.Linq;
using YamlDotNet.RepresentationModel;

namespace Constructs.Data;

internal class Array(YamlSequenceNode source) : Node(source) {
    private readonly YamlSequenceNode source = source;

    public int Count => this.source.Children.Count;

    public Node this[int index] => new(this.source[index]);

    public IEnumerable<Node> Items => this.source.Children.Select(x => new Node(x));

    public void PreallocateList<T>(List<T> list)
        where T : new() {
        for (int i = this.Count-1; i >= 0; --i)
            list.Add(new T());
    }

    public void InitializeList<T>(Project project, IReadOnlyList<T> list)
        where T : IInitializable {
        for (int i = 0; i < this.Count; ++i)
            list[i].Initialize(project, this[i]);
    }

    public List<T> AsIndexList<T>(IReadOnlyList<T> source)
        where T : IConstruct =>
        this.Items.Select(item => item.AsIndex(source)).ToList();

    public List<T> AsKeyList<T>(Project project)
        where T : IConstruct =>
        this.Items.Select(item => item.AsKey<T>(project)).ToList();
}
