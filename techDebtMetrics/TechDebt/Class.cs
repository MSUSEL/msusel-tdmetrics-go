using Commons.Data.Reader;
using System;
using System.Linq;
using System.Collections.Generic;

namespace TechDebt;

/// <summary>Represents a class in technical debt analysis.</summary>
public class Class(Source source): IComparable<Class>, IInitializable<Project.LoaderHelper>  {

    /// <summary>The name and location from the source code that this class came from.</summary>
    public readonly Source Source = source;

    /// <summary>The set of method partitations and fuzzy estimates.</summary>
    public readonly SortedSet<Participation> Participation = [];

    /// <summary>Gets the Weighted Method Count (WMC) for this class.</summary>
    public double WeigthedMethodCount =>
        this.Participation.Sum(p => p.Value * p.Method.CyclomaticComplexity);

    // TODO: Add other god class metrics.
    // TODO: Create a Metric result record that detects the god classes and allows them to be sorted.

    public override string ToString() => "class " + this.Source.ToString();

    public override int GetHashCode() => HashCode.Combine("class", this.Source);

    public override bool Equals(object? obj) => obj is Class m && this.Source == m.Source;

    public int CompareTo(Class? other) => other is null ? 1 : this.Source.CompareTo(other.Source);
    
    void IInitializable<Project.LoaderHelper>.Initialize(Project.LoaderHelper lh, int index, Node node) {
        // Currently this does nothing since the source is loaded before and the participation comes from the methods.
    }
}
