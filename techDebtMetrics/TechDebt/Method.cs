using Commons.Data.Yaml;
using Commons.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;
using TechDebt.Exceptions;

namespace TechDebt;

/// <summary>Represents a method in technical debt analysis.</summary>
public class Method(Source source) :
    IComparable<Method>,
    IInitializable<Project.LoaderHelper>,
    INodeable<Project.ToNodeHelper> {

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

    /// <summary>Adds a participation between this method and the class.</summary>
    /// <param name="value">The value between zero exclusively and one inclusively.</param>
    /// <param name="c">the class that is participating.</param>
    /// <returns>The created participation.</returns>
    public Participation Add(double value, Class c) => Project.Add(this, value, c);

    /// <summary>The McCabe Cyclomatic Complexity for this method.</summary>
    public int CyclomaticComplexity;

    /// <summary>The number of lines in the method that are not comments or empty.</summary>
    public int CodeCount;

    /// <summary>The set of classes that this method reads from.</summary>
    public readonly SortedSet<Class> Reads = [];

    /// <summary>The set of classes that this method writes to.</summary>
    public readonly SortedSet<Class> Writes = [];

    /// <summary>The set of methods that this method invokes.</summary>
    public readonly SortedSet<Method> Invokes = [];

    /// <summary>
    /// Normalizes this method so that the participation from this method sum up to one.
    /// This adjusts the participation values and removes any that are too smal.
    /// </summary>
    public void Normalize() {
        if (this.Participation.Count <= 0) throw new NoParticipationException(this);

        List<(Class Class, double Value)> entries = [.. this.Participation.Select(p => (p.Class, p.Value))];
        entries.Sort((p1, p2) => p1.Value.CompareTo(p2.Value));

        void apply(Func<double, double> h) {
            for (int i = 0; i < entries.Count; i++)
                entries[i] = (entries[i].Class, h(entries[i].Value));
        }

        double min = entries.Min(p => p.Value);
        if (min < 0.0) apply(v => v + min);

        double max = entries.Max(p => p.Value);
        if (Math.IsZero(max)) apply(v => v + 1.0);

        double sum = entries.Sum(p => p.Value);
        apply(v => v / sum);

        entries.RemoveAll(e => Math.IsZero(e.Value));

        // Do a second normalization pass after removing nearly zero entries.
        sum = entries.Sum(p => p.Value);
        apply(v => v / sum);

        this.Participation.ForAll(p => p.Class.Participation.Remove(p));
        this.Participation.Clear();
        entries.ForAll(e => Project.Add(this, e.Value, e.Class));
    }

    public override string ToString() => "method " + this.Source.ToString();
    public override int GetHashCode() => HashCode.Combine("method", this.Source);
    public override bool Equals(object? obj) => obj is Method m && this.Source == m.Source;
    public int CompareTo(Method? other) => other is null ? 1 : this.Source.CompareTo(other.Source);

    void IInitializable<Project.LoaderHelper>.Initialize(Project.LoaderHelper lh, int index, Node node) {
        Commons.Data.Yaml.Object obj = node.AsObject();

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
        this.CodeCount = obj.TryReadInt("codeCount");

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

    Node INodeable<Project.ToNodeHelper>.ToNode(Project.ToNodeHelper data, int index) {
        Commons.Data.Yaml.Object obj = new();
        obj.Add("name", this.Source.Name);
        obj.Add("loc", data.Locations, this.Source.Location);

        Commons.Data.Yaml.Object parts = new();
        foreach (Participation p in this.Participation) {
            int classIndex = data.Classes[p.Class];
            parts.Add(new Node(classIndex), new Node(p.Value));
        }
        obj.Add("participation", parts);

        obj.Add("complexity", this.CyclomaticComplexity);
        obj.Add("codeCount", this.CodeCount);
        obj.AddIndexList("reads", this.Reads, data.Classes);
        obj.AddIndexList("writes", this.Writes, data.Classes);
        obj.AddIndexList("invokes", this.Invokes, data.Methods);
        return obj;
    }
}
