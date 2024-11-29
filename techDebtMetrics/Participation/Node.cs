using System;

namespace Participation;

internal record struct Node(Edge[] Edges) {

    public void Insert(int index, Edge edge) {
        if (this.Edges is null) {
            this.Edges = [edge];
            return;
        }

        Edge[] edges = new Edge[this.Edges.Length+1];
        Array.Copy(this.Edges, edges, index);
        Array.Copy(this.Edges, index, edges, index + 1, edges.Length - index);
        edges[index] = edge;
        this.Edges = edges;
    }

    public void Remove(int index) {
        Edge[] edges = new Edge[this.Edges.Length+1];
        Array.Copy(this.Edges, edges, index - 1);
        Array.Copy(this.Edges, index - 1, edges, index, edges.Length - index - 1);
        this.Edges = edges;
    }

    public readonly (int index, bool found) FindEdge(Edge edge) {
        if (this.Edges is null) return (0, false);
        int index = Array.BinarySearch(this.Edges, edge);
        return index switch {
            >= 0 => (index, true),            // exact match
            ~0 => (this.Edges.Length, false), // not found
            _ => (~index - 1, false),         // nearest match
        };
    }
}
