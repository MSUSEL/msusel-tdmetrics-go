using System;

namespace Constructs.Tooling;

public readonly record struct Location(int LineNo, string Path) : IComparable<Location> {

    public int CompareTo(Location other) {
        int cmp = this.LineNo.CompareTo(other.LineNo);
        if (cmp != 0) return cmp;
        return this.Path.CompareTo(other.Path);
    }
}
