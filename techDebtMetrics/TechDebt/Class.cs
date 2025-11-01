using System;
using System.Linq;
using SCG = System.Collections.Generic;

namespace TechDebt;

/// <summary>Represents a class in technical debt analysis.</summary>
public class Class(Source source) {

    /// <summary>The name and location from the source code that this class came from.</summary>
    public readonly Source Source = source;

    /// <summary>The set of method partitations and fuzzy estimates.</summary>
    public readonly SCG.SortedSet<Participation> Participation = [];

    /// <summary>Gets the Weighted Method Count (WMC) for this class.</summary>
    public double WeigthedMethodCount =>
        this.Participation.Sum(p => p.Value * p.Method.CyclomaticComplexity);

    // TODO: Add other god class metrics.
    // TODO: Create a Metric result record that detects the god classes and allows them to be sorted.

    public override string ToString() => "class " + this.Source.ToString();

    public override int GetHashCode() => HashCode.Combine("class", this.Source);

    public override bool Equals(object? obj) => obj is Class m && this.Source == m.Source;

    public int CompareTo(Class? other) => other is null ? 1 : this.Source.CompareTo(other.Source);
}
