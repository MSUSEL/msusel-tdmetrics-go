using Commons.Data.Locations;

namespace TechDebt;

/// <summary>A source is the name and location of a construct in the source code.</summary>
/// <param name="Name">The name given to the construct.</param>
/// <param name="Location">The file path and line number that the construct came from.</param>
public readonly record struct Source(string Name, Location Location) : System.IComparable<Source> {

    /// <summary>The file path that the construct came from.</summary>
    public string Path => this.Location.Path;

    /// <summary>The line number in the file that the construct was at.</summary>
    public int LineNo => this.Location.LineNo;

    public override string ToString() => this.Name + " @ " + this.Location;

    public int CompareTo(Source other) {
        int cmp = this.Path.CompareTo(other.Path);
        return cmp != 0 ? cmp : this.Location.CompareTo(other.Location);
    }
}
