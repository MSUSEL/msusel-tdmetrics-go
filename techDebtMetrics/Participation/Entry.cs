using System;
using System.Collections.Generic;
using System.Diagnostics.CodeAnalysis;

namespace Participation;

/// <summary>An entry in the matrix that is returned while enumerating.</summary>
/// <param name="Row">The zero based row of the entry.</param>
/// <param name="Column">The zero based column of the entry.</param>
/// <param name="Value">The value at this entry.</param>
public readonly record struct Entry(int Row, int Column, double Value) {

    /// <summary>This is a comparer for performing epsilon comparisons for entries.</summary>
    /// <param name="epsilon">The epsilon for comparing the values of the entries.</param>
    public class Comparer(double epsilon) : IEqualityComparer<Entry> {
        public readonly double epsilon = epsilon;

        public bool Equals(Entry x, Entry y) =>
            x.Row == y.Row && x.Column == y.Column && Math.Abs(x.Value - y.Value) < this.epsilon;

        public int GetHashCode([DisallowNull] Entry obj) =>
            HashCode.Combine(obj.Row, obj.Column, obj.Value);
    }
}
