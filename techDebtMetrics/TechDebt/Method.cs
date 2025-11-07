using Commons.Data.Reader;
using System;
using System.Collections.Generic;

namespace TechDebt;

/// <summary>Represents a method in technical debt analysis.</summary>
public class Method(Source source) : IComparable<Method>, IInitializable<Project.LoaderHelper> {

    /// <summary>The name and location from the source code that this method came from.</summary>
    public readonly Source Source = source;

    /// <summary>The set of class partitations and fuzzy estimates.</summary>
    public readonly SortedSet<Participation> Participation = [];

    public bool TryFind(Class c, out Participation? p) {
        foreach (Participation par in this.Participation) {
            if (par.Class == c) {
                p = par;
                return true;
            }
        }
        p = null;
        return false;
    }

    /// <summary>The McCabe Cyclomatic Complexity for this method.</summary>
    public int CyclomaticComplexity;

    public readonly SortedSet<Class> Reads = [];

    public readonly SortedSet<Class> Writes = [];

    public readonly SortedSet<Method> Invokes = [];

    public override string ToString() => "method " + this.Source.ToString();

    public override int GetHashCode() => HashCode.Combine("method", this.Source);

    public override bool Equals(object? obj) => obj is Method m && this.Source == m.Source;

    public int CompareTo(Method? other) => other is null ? 1 : this.Source.CompareTo(other.Source);

    void IInitializable<Project.LoaderHelper>.Initialize(Project.LoaderHelper lh, int index, Node node) {
        Commons.Data.Reader.Object obj = node.AsObject();

        Node? parNode = obj.TryReadNode("participation");
        if (parNode is not null) {
            foreach (KeyValuePair<Node, Node> child in parNode.AsObject().Children) {
                Class c = child.Key.AsIndex(lh.Classes);
                double value = child.Value.AsDouble();
                Participation p = new(this, value, c);
                this.Participation.Add(p);
                c.Participation.Add(p);
            }
        }

        this.CyclomaticComplexity = obj.TryReadInt("complexity");

        List<Class> reads = [];
        obj.ReadIndexList("reads", reads, lh.Classes);
        reads.ForAll(this.Reads.Add);

        List<Class> writes = [];
        obj.ReadIndexList("writes", writes, lh.Classes);
        writes.ForAll(this.Writes.Add);

        List<Method> invokes = [];
        obj.ReadIndexList("invokes", invokes, lh.Methods);
        invokes.ForAll(this.Invokes.Add);
    }
}
