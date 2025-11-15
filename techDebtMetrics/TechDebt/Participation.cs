namespace TechDebt;

/// <summary>The participation (estimation of membership) that the class and method has.</summary>
/// <param name="Method">The method that the class is participating with.</param>
/// <param name="Value">The value between zero exclusively and one inclusively.</param>
/// <param name="Class">The class that participates with the method.</param>
public readonly record struct Participation(Method Method, double Value, Class Class) : System.IComparable<Participation> {

    public int CompareTo(Participation other) {
        int cmp = this.Method.CompareTo(other.Method);
        if (cmp != 0) return cmp;
        cmp = this.Class.CompareTo(other.Class);
        if (cmp != 0) return cmp;
        return this.Value.CompareTo(other.Value);
    }

    public override string ToString() => this.Method + " -> " + this.Value + " -> " + this.Class;
}
