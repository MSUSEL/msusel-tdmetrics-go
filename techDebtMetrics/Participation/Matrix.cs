using System;
using System.Collections;
using System.Collections.Generic;
using System.Text;

namespace Participation;

public class Matrix(int rowNo, int columnNo, double epsilon = 1.0e-9) : IEnumerable<Entry> {
    public readonly int Rows = rowNo;
    public readonly int Columns = columnNo;
    public readonly double Epsilon = epsilon;

    private readonly SortedDictionary<int, double>[] rows = new SortedDictionary<int, double>[rowNo];

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
        SortedDictionary<int, double> node = this.rows[row];
        if (node is null) return 0.0;
        node.TryGetValue(column, out double value);
        return value;
    }

    private void setValue(int row, int column, double value) {
        this.checkRange(row, column);
        SortedDictionary<int, double> node = this.rows[row];
        if (node is null) this.rows[row] = node = [];
        node[column] = value;
    }

    private bool removeValue(int row, int column) {
        this.checkRange(row, column);
        return this.rows[row]?.Remove(column) ?? false;
    }

    public IEnumerator<Entry> GetEnumerator(bool full) =>
        full ? this.fullEnumerator() : this.shortEnumerator();

    public IEnumerator<Entry> GetEnumerator() => this.GetEnumerator(true);
    IEnumerator IEnumerable.GetEnumerator() => this.GetEnumerator(true);

    private IEnumerator<Entry> shortEnumerator() {
        for (int row = 0; row < this.Rows; ++row)
            foreach (KeyValuePair<int, double> edge in this.rows[row])
                yield return new(row, edge.Key, edge.Value);
    }

    private IEnumerator<Entry> fullEnumerator() {
        for (int row = 0; row < this.Rows; ++row) {
            SortedDictionary<int, double> node = this.rows[row];
            int next = 0;
            if (node is not null) {
                foreach (KeyValuePair<int, double> edge in node) {
                    for (int column = next; column < edge.Key; ++column)
                        yield return new(row, column, 0.0);
                    yield return new(row, edge.Key, edge.Value);
                    next = edge.Key + 1;
                }
            }
            for (int column = next; column < this.Columns; ++column)
                yield return new(row, column, 0.0);
        }
    }

    public string Serialize() {
        StringBuilder sb = new();
        string version = "0"; // currently only version zero.
        sb.AppendFormat("{0} {1}x{2}", version, this.Rows, this.Columns);
        foreach (SortedDictionary<int, double> node in this.rows) {
            sb.Append('\n');
            if (node is not null) {
                bool first = true;
                foreach (KeyValuePair<int, double> edge in node) {
                    if (first) first = false;
                    else sb.Append(' ');
                    sb.AppendFormat("{0}:{1:0.0#####}", edge.Key, edge.Value);
                }
            }
        }
        return sb.ToString();
    }

    static public Matrix Deserialize(string data, double epsilon = 1.0e-9) {
        string[] lines = data.Split('\n');
        Matrix m = deserializeHeader(lines[0], epsilon);
        for (int i = 1; i < lines.Length; ++i)
            deserializeRow(m, i - 1, lines[i]);
        return m;
    }

    static private Matrix deserializeHeader(string header, double epsilon) {
        string[] parts = header.Trim().Split(' ');
        if (parts.Length != 2)
            throw new Exception("Expected the header to be two parts separated by a space, but got \"" + header + "\".");

        string version = parts[0].Trim();
        if (version != "0") // currently only version zero.
            throw new Exception("Unexpected version in header, got \"" + version + "\".");

        string[] sizes = parts[1].Trim().Split('x');
        if (sizes.Length != 2)
            throw new Exception("Expected the header size to be two values separated by an 'x', but got \"" + parts[1] + "\".");

        if (!int.TryParse(sizes[0].Trim(), out int rows))
            throw new Exception("Expected the number of rows to be an integer, but got \"" + sizes[0] + "\".");
        if (rows <= 0)
            throw new Exception("Expected the number of rows to be at least one, but got " + rows);

        if (!int.TryParse(sizes[1].Trim(), out int columns))
            throw new Exception("Expected the number of columns to be an integer, but got \"" + sizes[1] + "\".");
        if (columns <= 0)
            throw new Exception("Expected the number of columns to be at least one, but got " + columns);

        return new(rows, columns, epsilon);
    }

    static private void deserializeRow(Matrix m, int row, string line) {
        if (row < 0 || row >= m.Rows)
            throw new Exception("Expected a row number to be in the range [0.." + m.Rows + "), but it was " + row + ".");
        if (string.IsNullOrWhiteSpace(line)) return;

        string[] columns = line.Trim().Split(' ');
        SortedDictionary<int, double> data = [];
        foreach (string entry in columns) {
            string[] pair = entry.Trim().Split(':');
            if (pair.Length != 2)
                throw new Exception("Expected entry pairs separated by a ':', but got \"" + entry + "\" on row " + row + ".");

            if (!int.TryParse(pair[0].Trim(), out int column))
                throw new Exception("Expected the column number to be an integer, but got \"" + pair[0] + "\" on row " + row + ".");
            if (column < 0 || column >= m.Columns)
                throw new Exception("Expected a column number to be in the range [0.." + m.Columns + "), but it was " + column + ".");

            if (!double.TryParse(pair[1].Trim(), out double value))
                throw new Exception("Expected the edge weight to be a double, but got \"" + pair[1] + "\" on row " + row + ".");

            data[column] = value;
        }
        m.rows[row] = data;
    }

    public override string ToString() => this.ToString("{0:0.0###;-0.0###;-}");

    public string ToString(string format) {
        int[] lefts = new int[this.Columns];
        int[] rights = new int[this.Columns];
        StringBuilder sb = new();
        foreach (Entry entry in this) {
            (int left, int right) = measureNumber(string.Format(format, entry.Value));
            if (lefts[entry.Column] < left) lefts[entry.Column] = left;
            if (rights[entry.Column] < right) rights[entry.Column] = right;
        }

        sb.Append("[[ ");
        int maxRow = this.Rows - 1;
        int maxCol = this.Columns - 1;
        foreach (Entry entry in this) {
            string text = string.Format(format, entry.Value);
            (int left, int right) = measureNumber(text);
            sb.Append(' ', lefts[entry.Column] - left);
            sb.Append(text);
            if (entry.Column >= maxCol) {
                sb.Append(' ', rights[entry.Column] - right);
                if (entry.Row >= maxRow) sb.Append(" ]]");
                else sb.Append(" ],\n [ ");
            } else {
                sb.Append(", ");
                sb.Append(' ', rights[entry.Column] - right);
            }
        }
        return sb.ToString();
    }

    private static (int left, int right) measureNumber(string text) {
        if (text == "-") return (0, 1);
        char[] centers = ['.', 'e', 'E'];
        string[] parts = text.Split(centers, 2);
        return (parts[0].Length, parts.Length > 1 ? parts[1].Length + 1 : 0);
    }
}
