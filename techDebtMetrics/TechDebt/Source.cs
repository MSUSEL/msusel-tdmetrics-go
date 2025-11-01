namespace TechDebt;

/// <summary>A source is the name and location of a construct in the source code.</summary>
/// <param name="Name">The name given to the construct.</param>
/// <param name="File">The file path that the construct came from.</param>
/// <param name="Line">The line number in the file that the construct was at.</param>
public readonly record struct Source(string Name, string File, int Line) : System.IComparable<Source> {

    public override string ToString() => this.Name + " @ " + this.File + ":" + this.Line;

    public int CompareTo(Source other) {
        int cmp = this.File.CompareTo(other.File);
        if (cmp != 0) return cmp;
        cmp = this.Line.CompareTo(other.Line);
        if (cmp != 0) return cmp;
        return this.Name.CompareTo(other.Name);
    }
}
