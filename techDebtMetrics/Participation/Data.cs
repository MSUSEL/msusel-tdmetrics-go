using System;
using System.Collections;
using System.Collections.Generic;
using System.Text;

namespace Participation;

public abstract class Data : IEnumerable<Entry> {

    public abstract int Rows { get; }
    public abstract int Columns { get; }
    public abstract double Epsilon { get; }

    protected abstract bool ColumnHasZero(int column);

    public IEnumerator<Entry> GetEnumerator() => this.FullEnumerate().GetEnumerator();
    IEnumerator IEnumerable.GetEnumerator() => this.FullEnumerate().GetEnumerator();

    public abstract IEnumerable<Entry> ShortEnumerate();

    public abstract IEnumerable<Entry> FullEnumerate();

    public double this[int row, int column] {
        get {
            this.CheckRange(row, column);
            return this.GetValue(row, column);
        }
        set {
            this.CheckRange(row, column);
            if (double.Abs(value) < this.Epsilon)
                this.RemoveValue(row, column);
            else this.SetValue(row, column, value);
        }
    }

    protected void CheckRange(int row, int column) {
        if (row < 0 || row >= this.Rows)
            throw new IndexOutOfRangeException("Row must be in [0.." + this.Rows + "), the given row was " + row);
        if (column < 0 || column >= this.Columns)
            throw new IndexOutOfRangeException("Column must be in [0.." + this.Columns + "), the given column was " + column);
    }

    protected abstract double GetValue(int row, int column);

    protected abstract void SetValue(int row, int column, double value);

    protected abstract bool RemoveValue(int row, int column);

    protected abstract SortedDictionary<int, double> GetColumnNode(int column);

    protected abstract SortedDictionary<int, double> GetRowNode(int row);

    public override string ToString() => this.ToString("{0:0.0###;-0.0###;-}");

    public string ToString(string format) {
        int[] lefts = new int[this.Columns];
        int[] rights = new int[this.Columns];
        StringBuilder sb = new();
        foreach (Entry entry in this.ShortEnumerate()) {
            (int left, int right) = measureNumber(string.Format(format, entry.Value));
            if (lefts[entry.Column] < left) lefts[entry.Column] = left;
            if (rights[entry.Column] < right) rights[entry.Column] = right;
        }
        for (int column = 0; column < this.Columns; ++column) {
            if (this.ColumnHasZero(column)) {
                // TODO: Measure once
                (int left, int right) = measureNumber(string.Format(format, 0.0));
                if (lefts[column] < left) lefts[column] = left;
                if (rights[column] < right) rights[column] = right;
            }
        }

        sb.Append("[[ ");
        int maxRow = this.Rows - 1;
        int maxCol = this.Columns - 1;
        foreach (Entry entry in this.FullEnumerate()) {
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

        // TODO: Improve zero values

        if (text == "-") return (0, 1);
        char[] centers = ['.', 'e', 'E'];
        string[] parts = text.Split(centers, 2);
        return (parts[0].Length, parts.Length > 1 ? parts[1].Length + 1 : 0);
    }

    static protected bool zipOr(SortedDictionary<int, double> left, SortedDictionary<int, double> right, Action<int, double, double> handle) {
        if (left is null) {
            if (right is null) return false;

            foreach (KeyValuePair<int, double> edge in right)
                handle(edge.Key, 0.0, edge.Value);
            return right.Count > 0;
        }

        if (right is null) {
            foreach (KeyValuePair<int, double> edge in left)
                handle(edge.Key, edge.Value, 0.0);
            return left.Count > 0;
        }

        SortedDictionary<int, double>.Enumerator leftIt = left.GetEnumerator();
        SortedDictionary<int, double>.Enumerator rightIt = right.GetEnumerator();
        KeyValuePair<int, double>? leftEdge = leftIt.MoveNext() ? leftIt.Current : null;
        KeyValuePair<int, double>? rightEdge = rightIt.MoveNext() ? rightIt.Current : null;
        bool called = leftEdge is not null || rightEdge is not null;

        while (leftEdge is not null && rightEdge is not null) {
            if (leftEdge.Value.Key < rightEdge.Value.Key) {
                handle(leftEdge.Value.Key, leftEdge.Value.Value, 0.0);
                leftEdge = leftIt.MoveNext() ? leftIt.Current : null;

            } else if (leftEdge.Value.Key > rightEdge.Value.Key) {
                handle(rightEdge.Value.Key, 0.0, rightEdge.Value.Value);
                rightEdge = rightIt.MoveNext() ? rightIt.Current : null;

            } else {
                handle(leftEdge.Value.Key, leftEdge.Value.Value, rightEdge.Value.Value);
                leftEdge = leftIt.MoveNext() ? leftIt.Current : null;
                rightEdge = rightIt.MoveNext() ? rightIt.Current : null;
            }
        }

        while (leftEdge is not null) {
            handle(leftEdge.Value.Key, leftEdge.Value.Value, 0.0);
            leftEdge = leftIt.MoveNext() ? leftIt.Current : null;
        }

        while (rightEdge is not null) {
            handle(rightEdge.Value.Key, 0.0, rightEdge.Value.Value);
            rightEdge = rightIt.MoveNext() ? rightIt.Current : null;
        }
        return called;
    }

    static protected bool zipAnd(SortedDictionary<int, double> left, SortedDictionary<int, double> right, Action<int, double, double> handle) {
        if (left is null || right is null) return false;

        SortedDictionary<int, double>.Enumerator leftIt = left.GetEnumerator();
        SortedDictionary<int, double>.Enumerator rightIt = right.GetEnumerator();
        KeyValuePair<int, double>? leftEdge = leftIt.MoveNext() ? leftIt.Current : null;
        KeyValuePair<int, double>? rightEdge = rightIt.MoveNext() ? rightIt.Current : null;
        bool called = false;

        while (leftEdge is not null && rightEdge is not null) {
            if (leftEdge.Value.Key < rightEdge.Value.Key)
                leftEdge = leftIt.MoveNext() ? leftIt.Current : null;

            else if (leftEdge.Value.Key > rightEdge.Value.Key)
                rightEdge = rightIt.MoveNext() ? rightIt.Current : null;

            else {
                handle(leftEdge.Value.Key, leftEdge.Value.Value, rightEdge.Value.Value);
                leftEdge = leftIt.MoveNext() ? leftIt.Current : null;
                rightEdge = rightIt.MoveNext() ? rightIt.Current : null;
                called = true;
            }
        }
        return called;
    }

    public string Serialize() {
        StringBuilder sb = new();
        string version = "0"; // currently the only version is "0".
        sb.AppendFormat("{0} {1}x{2}", version, this.Rows, this.Columns);
        for (int rows = 0; rows < this.Rows; ++rows) {
            SortedDictionary<int, double> node = this.GetRowNode(rows);
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

    static protected T deserialize<T>(Func<int, int, T> factory, string data)
        where T: Data {
        string[] lines = data.Split('\n');
        (int rows, int columns) = deserializeHeader(lines[0]);
        T d = factory(rows, columns);
        for (int i = 1; i < lines.Length; ++i)
            deserializeRow(d, i - 1, lines[i]);
        return d;
    }

    static private (int rows, int columns) deserializeHeader(string header) {
        string[] parts = header.Trim().Split(' ');
        if (parts.Length != 2)
            throw new Exception("Expected the header to be two parts separated by a space, but got \"" + header + "\".");

        string version = parts[0].Trim();
        if (version != "0") // currently the only version is "0".
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

        return (rows, columns);
    }

    static private void deserializeRow(Data d, int row, string line) {
        if (row < 0 || row >= d.Rows)
            throw new Exception("Expected a row number to be in the range [0.." + d.Rows + "), but it was " + row + ".");
        if (string.IsNullOrWhiteSpace(line)) return;

        string[] columns = line.Trim().Split(' ');
        foreach (string entry in columns) {
            string[] pair = entry.Trim().Split(':');
            if (pair.Length != 2)
                throw new Exception("Expected entry pairs separated by a ':', but got \"" + entry + "\" on row " + row + ".");

            if (!int.TryParse(pair[0].Trim(), out int column))
                throw new Exception("Expected the column number to be an integer, but got \"" + pair[0] + "\" on row " + row + ".");
            if (column < 0 || column >= d.Columns)
                throw new Exception("Expected a column number to be in the range [0.." + d.Columns + "), but it was " + column + ".");

            if (!double.TryParse(pair[1].Trim(), out double value))
                throw new Exception("Expected the edge weight to be a double, but got \"" + pair[1] + "\" on row " + row + ".");

            if (Math.Abs(value) >= d.Epsilon)
                d.SetValue(row, column, value);
        }
    }
}
