using System;
using System.Collections;
using System.Collections.Generic;
using System.Text;
using System.Text.RegularExpressions;

namespace Participation;

public class Matrix(int rows, int columns, double epsilon = 1.0e-9) : IEnumerable<Entry> {
    public readonly int Rows = rows;
    public readonly int Columns = columns;
    public readonly double Epsilon = epsilon;

    private Node[] nodes = [];

    public double this[int row, int column] {
        get => this.getValue(row, column);
        set {
            if (double.Abs(value) < this.Epsilon)
                this.removeValue(row, column);
            else this.setValue(row, column, value);
        }
    }

    private void checkRange(int row, int column) {
        if (row < 0 || row >= this.Rows)
            throw new IndexOutOfRangeException("Row must be in [0.." + this.Rows + "), the row was " + row);
        if (column < 0 || column >= this.Columns)
            throw new IndexOutOfRangeException("Column must be in [0.." + this.Columns + "), the column was " + column);
    }

    private double getValue(int row, int column) {
        this.checkRange(row, column);
        if (this.nodes.Length <= 0) return 0.0;
        Node node = this.nodes[row];
        Edge edge = new(column, 0.0);
        (int index, bool found) = node.FindEdge(edge);
        if (found) return node.Edges[index].Value;
        return 0.0;
    }

    private void setValue(int row, int column, double value) {
        this.checkRange(row, column);
        if (this.nodes.Length <= 0)
            this.nodes = new Node[this.Rows];

        Node node = this.nodes[row];
        Edge edge = new(column, value);
        (int index, bool found) = node.FindEdge(edge);
        if (found) node.Edges[index] = edge;
        else node.Insert(index, edge);
        this.nodes[row] = node;
    }

    private bool removeValue(int row, int column) {
        this.checkRange(row, column);
        if (this.nodes.Length <= 0) return false;

        Node node = this.nodes[row];
        Edge edge = new(column, 0.0);
        (int index, bool found) = node.FindEdge(edge);
        if (found) {
            node.Remove(index);
            this.nodes[row] = node;
            return true;
        }
        return false;
    }

    public IEnumerator<Entry> GetEnumerator(bool full) =>
        full ? this.fullEnumerator() : this.shortEnumerator();

    public IEnumerator<Entry> GetEnumerator() => this.GetEnumerator(true);
    IEnumerator IEnumerable.GetEnumerator() => this.GetEnumerator(true);

    private IEnumerator<Entry> shortEnumerator() {
        if (this.nodes is null || this.nodes.Length <= 0) yield break;

        for (int row = 0; row < this.Rows; ++row) {
            Node node = this.nodes[row];
            foreach (Edge edge in node.Edges)
                yield return new(row, edge.Column, edge.Value);
        }
    }
    
    private IEnumerator<Entry> fullEnumerator() {
        if (this.nodes is null || this.nodes.Length <= 0) {
            for (int row = 0; row < this.Rows; ++row) {
                for (int column = 0; column < this.Columns; ++column) {
                    yield return new(row, column, 0.0);
                }
            }
            yield break;
        }

        for (int row = 0; row < this.Rows; ++row) {
            Node node = this.nodes[row];
            if (node.Edges is null) {
                for (int column = 0; column < this.Columns; ++column)
                    yield return new(row, column, 0.0);
            } else {
                for (int column = 0, index = 0; column < this.Columns; ++column) {
                    if (index < node.Edges.Length) {
                        Edge edge = node.Edges[index];
                        if (edge.Column == column) {
                            yield return new(row, column, edge.Value);
                            ++index;
                            continue;
                        }
                    }
                    yield return new(row, column, 0.0);
                }
            }
        }
    }

    public string Serialize() {
        StringBuilder sb = new();
        sb.AppendFormat("0 {0}x{1}", this.Rows, this.Columns);
        foreach (Node node in this.nodes) {
            sb.Append('\n');
            bool first = true;
            if (node.Edges is not null) {
                foreach (Edge edge in node.Edges) {
                    if (first) first = false;
                    else sb.Append(' ');
                    sb.AppendFormat("{0}:{1:0.0000}", edge.Column, edge.Value);
                }
            }
        }
        return sb.ToString();
    }

    static public Matrix Deserialize(string data, double epsilon = 1.0e-9) {
        string[] lines = data.Split('\n');
        Regex header = new(@"^(\S+)\s+(\d+)x(\d)\s*$", RegexOptions.Compiled);

        Match mh = header.Match(lines[0]);
        if (!mh.Success || mh.Captures.Count != 3)
            throw new Exception("Invalid header in serialization.");
        string version = mh.Captures[0].Value;
        if (version != "0")
            throw new Exception("Unknown version: " + version);

        int rows    = int.Parse(mh.Captures[1].Value);
        int columns = int.Parse(mh.Captures[2].Value);
        Matrix m = new(rows, columns);

        int count = int.Max(rows, lines.Length-1);
        if (count <= 0) return m;

        // TODO: Add more checks to ensure a good deserialize.
        for (int i = 1; i < count; ++i) {
            string[] parts = lines[i].Split(' ');
            Edge[] edges = new Edge[parts.Length];
            for (int j = 0; j < parts.Length; ++j) {
                string[] p = parts[i].Split(':');
                int column = int.Parse(p[0].Trim());
                double value = double.Parse(p[1].Trim());
                edges[i] = new Edge(column, value);
            }
            m.nodes[i - 1].Edges = edges;
        }
        return m;
    }

    public override string ToString() => this.ToString("{0:0.0000}");

    public string ToString(string format) {
        int[] widths = new int[this.Columns];
        StringBuilder sb = new();
        foreach (Entry entry in this) {
            string text = string.Format(format, entry.Value);
            if (widths[entry.Column] < text.Length)
                widths[entry.Column] = text.Length;
        }

        sb.Append("[[");
        foreach (Entry entry in this) {
            if (entry.Column == 0) {
                if (entry.Row > 0) sb.Append("],\n [");
            } else if (entry.Column > 0) sb.Append(", ");
            string text = string.Format(format, entry.Value).
                PadLeft(widths[entry.Column], ' ');
            sb.Append(text);
        }
        sb.Append("]]");
        return sb.ToString();
    }
}
