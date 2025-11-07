using System;
using SCG = System.Collections.Generic;
using System.Linq;

namespace TechDebt;

/// <summary>Represents a method in technical debt analysis.</summary>
public class Method(Source source) : IComparable<Method> {

    /// <summary>The name and location from the source code that this method came from.</summary>
    public readonly Source Source = source;

    /// <summary>The set of class partitations and fuzzy estimates.</summary>
    public readonly SCG.SortedSet<Participation> Participation = [];

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
    public double CyclomaticComplexity;

    public readonly SCG.SortedSet<Class> Reads = [];

    public readonly SCG.SortedSet<Class> Writes = [];

    public readonly SCG.SortedSet<Method> Invokes = [];

    public override string ToString() => "method " + this.Source.ToString();

    public override int GetHashCode() => HashCode.Combine("method", this.Source);

    public override bool Equals(object? obj) => obj is Method m && this.Source == m.Source;

    public int CompareTo(Method? other) => other is null ? 1 : this.Source.CompareTo(other.Source);
}
