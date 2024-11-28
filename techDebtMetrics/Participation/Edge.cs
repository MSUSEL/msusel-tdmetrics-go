using System;

namespace Participation;

internal readonly record struct Edge(int Column, double Value) : IComparable<Edge> {
    public int CompareTo(Edge other) => this.Column.CompareTo(other.Column);
}
