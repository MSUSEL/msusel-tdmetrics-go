using System;

namespace Commons.Data.Locations;

/// <summary>The location in the source code that some construct came from.</summary>
/// <param name="LineNo">The line number in the file the construct came from.</param>
/// <param name="Path">The path of the source file the construct came from.</param>
public readonly record struct Location(int LineNo, string Path) : IComparable<Location> {

    /// <summary>The location returned for an unknown location.</summary>
    static public readonly Location Unknown = new(0, "<unknown>");

    /// <summary>Compares this location with another location.</summary>
    /// <param name="other">The other location to compare against.</param>
    /// <returns>The comparison results.</returns>
    public int CompareTo(Location other) {
        int cmp = this.Path.CompareTo(other.Path);
        return cmp != 0 ? cmp :
            this.LineNo.CompareTo(other.LineNo);
    }
}
