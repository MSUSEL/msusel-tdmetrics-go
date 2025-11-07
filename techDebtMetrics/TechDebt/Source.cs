namespace TechDebt;

/// <summary>A source is the name and location of a construct in the source code.</summary>
/// <param name="Name">The name given to the construct.</param>
/// <param name="Path">The file path that the construct came from.</param>
/// <param name="LineNo">The line number in the file that the construct was at.</param>
public readonly record struct Source(string Name, string Path, int LineNo) : System.IComparable<Source> {

    public override string ToString() => this.Name + " @ " + this.Path + ":" + this.LineNo;

    public int CompareTo(Source other) {
        int cmp = this.Path.CompareTo(other.Path);
        if (cmp != 0) return cmp;
        cmp = this.LineNo.CompareTo(other.LineNo);
        if (cmp != 0) return cmp;
        return this.Name.CompareTo(other.Name);
    }
}
