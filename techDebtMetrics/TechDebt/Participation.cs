namespace TechDebt;

public readonly record struct Participation(Method Method, double Value, Class Class) : System.IComparable<Participation> {

    public int CompareTo(Participation other) {
        int cmp = this.Method.CompareTo(other.Method);
        if (cmp != 0) return cmp;
        cmp = this.Class.CompareTo(other.Class);
        if (cmp != 0) return cmp;
        return this.Value.CompareTo(other.Value);
    }
}
