using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Participation.LinearMath;

/// <summary>This is the shared parts of a sparse matrix and vector.</summary>
public abstract class Data : IEnumerable<Entry> {

    /// <summary>The default epsilon value to use for comparing data.-</summary>
    public const double DefaultEpsilon = 1.0e-9;

    /// <summary>Creates a new data.</summary>
    /// <param name="rows">The number of rows for the data.</param>
    /// <param name="columns">The number of columns for the data.</param>
    /// <param name="epsilon">The epsilon comparator for the values in the data.</param>
    protected Data(int rows, int columns, double epsilon) {
        if (rows < 1)
            throw new ArgumentException("Must have a positive non-zero number of rows.", nameof(rows));
        if (columns < 1)
            throw new ArgumentException("Must have a positive non-zero number of columns.", nameof(columns));
        if (double.IsNegative(epsilon))
            throw new ArgumentException("Must have a positive value epsilon.", nameof(epsilon));

        this.Rows = rows;
        this.Columns = columns;
        this.Epsilon = epsilon;
    }

    /// <summary>The fixed number of rows.</summary>
    public readonly int Rows;

    /// <summary>The fixed number of columns.</summary>
    public readonly int Columns;

    /// <summary>The epsilon for an epsilon comparison.</summary>
    public readonly double Epsilon;

    /// <summary>Enumerates all of the values in the data.</summary>
    /// <returns>The enumerator of all the values in the data.</returns>
    public IEnumerator<Entry> GetEnumerator() => this.FullEnumerate().GetEnumerator();
    IEnumerator IEnumerable.GetEnumerator() => this.FullEnumerate().GetEnumerator();

    /// <summary>Enumerates only the non-zero values in the data.</summary>
    /// <returns>The enumerator of all non-zero values in the data.</returns>
    public abstract IEnumerable<Entry> ShortEnumerate();

    /// <summary>Enumerates all of the values in the data.</summary>
    /// <returns>The enumerator of all the values in the data.</returns>
    public abstract IEnumerable<Entry> FullEnumerate();

    /// <summary>Gets or sets the value at the given row and column.</summary>
    /// <param name="row">The row to get or set, [0..Rows).</param>
    /// <param name="column">The column to get or set, [0..Columns).</param>
    /// <returns>The value at the given row and column.</returns>
    public double this[int row, int column] {
        get {
            this.CheckRange(row, column);
            return this.GetValue(row, column);
        }
        set {
            this.CheckRange(row, column);
            if (this.IsZero(value))
                this.RemoveValue(row, column);
            else this.SetValue(row, column, value);
        }
    }

    /// <summary>Gets or sets the value at the given row and column.</summary>
    /// <param name="row">The row to get or set, [0..Rows).</param>
    /// <param name="column">The column to get or set, [0..Columns).</param>
    /// <returns>The value at the given row and column.</returns>
    public double this[Index row, Index column] {
        get => this[row.GetOffset(this.Rows), column.GetOffset(this.Columns)];
        set => this[row.GetOffset(this.Rows), column.GetOffset(this.Columns)] = value;
    }

    /// <summary>Determines if the given value is zero with the epsilon.</summary>
    /// <param name="value">The value to check if zero.</param>
    /// <returns>True if the value is zero, false otherwise.</returns>
    protected bool IsZero(double value) => double.Abs(value) < this.Epsilon;

    /// <summary>Determines if the two values are equal with the given epsilon.</summary>
    /// <param name="x">The first value in the check.</param>
    /// <param name="y">The second value in the check.</param>
    /// <returns>True if the two values are equal, false otherwise.</returns>
    protected bool Equal(double x, double y) => this.IsZero(x - y);

    /// <summary>Determines if the given column has at least one zero.</summary>
    /// <param name="column">The valid column to check.</param>
    /// <returns>True if the given column contains a zero, otherwise false.</returns>
    protected abstract bool ColumnHasZero(int column);

    /// <summary>Checks if the given row and column are valid, otherwise an exception is thrown.</summary>
    /// <param name="row">The row to check.</param>
    /// <param name="column">The column to check.</param>
    protected void CheckRange(int row, int column) {
        this.CheckRow(row);
        this.CheckColumn(column);
    }

    /// <summary>Checks if the given row are valid, otherwise an exception is thrown.</summary>
    /// <param name="row">The row to check.</param>
    protected void CheckRow(int row) {
        if (row < 0 || row >= this.Rows)
            throw new IndexOutOfRangeException("Row must be in [0.." + this.Rows + "), the given row was " + row);
    }

    /// <summary>Checks if the given column are valid, otherwise an exception is thrown.</summary>
    /// <param name="column">The column to check.</param>
    protected void CheckColumn(int column) {
        if (column < 0 || column >= this.Columns)
            throw new IndexOutOfRangeException("Column must be in [0.." + this.Columns + "), the given column was " + column);
    }

    /// <summary>Sets the given value if the value is non-zero, otherwise has no effect.</summary>
    /// <param name="row">The valid row to set to.</param>
    /// <param name="column">The valid column to set to.</param>
    /// <param name="value">The value to set if non-zero.</param>
    internal void SetIfNonZero(int row, int column, double value) {
        if (!this.IsZero(value)) this.SetValue(row, column, value);
    }

    /// <summary>Gets the value at the given row and column.</summary>
    /// <param name="row">The valid row to get from.</param>
    /// <param name="column">The value column to get from.</param>
    /// <returns>The value at the given row and column.</returns>
    protected abstract double GetValue(int row, int column);

    /// <summary>This will set the non-zero value to the given row and column.</summary>
    /// <param name="row">The valid row to set to.</param>
    /// <param name="column">The valid column to set to.</param>
    /// <param name="value">The non-zero value to set.</param>
    protected abstract void SetValue(int row, int column, double value);

    /// <summary>Removes the value at the given row and column.</summary>
    /// <remarks>This is called when setting a zero value.</remarks>
    /// <param name="row">The valid row to remove.</param>
    /// <param name="column">The valid column to remove.</param>
    /// <returns>True if the value was removed, false if the value didn't exist or was already zero.</returns>
    protected abstract bool RemoveValue(int row, int column);

    /// <summary>
    /// Gets the whole column as a dictionary containing the non-zero values in the column,
    /// keyed with the column number of each non-zero value.
    /// </summary>
    /// <param name="column">The column to get.</param>
    /// <returns>The dictionary containing the column numbers paired with the non-zero values.</returns>
    internal abstract SortedDictionary<int, double> GetColumnNode(int column);

    /// <summary>
    /// Gets the whole row as a dictionary containing the non-zero values in the row,
    /// keyed with the row number of each non-zero value.
    /// </summary>
    /// <param name="column">The row to get.</param>
    /// <returns>The dictionary containing the row numbers paired with the non-zero values.</returns>
    internal abstract SortedDictionary<int, double> GetRowNode(int row);

    /// <summary>Gets the string for the data.</summary>
    /// <returns>The string for the data.</returns>
    public override string ToString() => this.ToString("{0:0.0###;-0.0###;-}");

    /// <summary>Gets the string for the data.</summary>
    /// <param name="format">The format to use for each number in the data.</param>
    /// <returns>The string for the data.</returns>
    public string ToString(string format) {
        int[] lefts = new int[this.Columns];
        int[] rights = new int[this.Columns];
        StringBuilder sb = new();
        foreach (Entry entry in this.ShortEnumerate()) {
            (int left, int right) = measureNumber(string.Format(format, entry.Value));
            if (lefts[entry.Column] < left) lefts[entry.Column] = left;
            if (rights[entry.Column] < right) rights[entry.Column] = right;
        }
        bool zeroNotMeasured = true;
        int zeroLeft = 0, zeroRight = 0;
        for (int column = 0; column < this.Columns; ++column) {
            if (this.ColumnHasZero(column)) {
                if (zeroNotMeasured) {
                    (zeroLeft, zeroRight) = measureNumber(string.Format(format, 0.0));
                    zeroNotMeasured = false;
                }
                if (lefts[column] < zeroLeft) lefts[column] = zeroLeft;
                if (rights[column] < zeroRight) rights[column] = zeroRight;
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

    /// <summary>Determines if this data is equal to the given object.</summary>
    /// <param name="obj">The other data to check.</param>
    /// <returns>True if the other data is equal to this data.</returns>
    public override bool Equals(object? obj) {
        if (obj is null) return false;
        if (ReferenceEquals(this, obj)) return true;
        if (obj is not Data other ||
            this.Rows != other.Rows ||
            this.Columns != other.Columns) return false;
        return this.ShortEnumerate().SequenceEqual(other.ShortEnumerate(), new Entry.Comparer(this.Epsilon));
    }

    /// <summary>Gets the hash for this data.</summary>
    /// <returns>The hash of this data.</returns>
    public override int GetHashCode() =>
        HashCode.Combine(this.Rows, this.Columns, this.ShortEnumerate());

    /// <summary>
    /// This measures the left and right side of a decimal point in the stringified number.
    /// This will attempt to center non-numbers and numbers without a decimal point.
    /// </summary>
    /// <remarks>
    /// This could be made better to handle numbers without a decimal point, check for numbers
    /// next to 'e' or 'E', and handle numbers using decimal point other than '.', however,
    /// for now, that isn't required for this project.
    /// This works well for the default formatting passed into ToString().
    /// </remarks>
    /// <param name="text">The stringified number to measure the size of.</param>
    /// <returns>
    /// The size of the number to the left of the decimal point
    /// and to the right and including the decimal point.
    /// </returns>
    private static (int left, int right) measureNumber(string text) {
        int length = text.Length;
        if (length <= 0) return (0, 0);

        char[] centers = ['.', 'e', 'E'];
        int index = text.IndexOfAny(centers);
        if (index < 0) index = length >> 1;
        return (index, length - index);
    }

    /// <summary>
    /// This zips together two sorted dictionaries with the given actions.
    /// If the left OR the right has a value for a column, this will call the given action.
    /// Zero will be used for whichever dictionary doesn't have a value for a column.
    /// </summary>
    /// <param name="left">The left dictionary in the zip. This may be null.</param>
    /// <param name="right">The right dictionary in the zip. This may be null.</param>
    /// <param name="handle">The handle for processing each column found in either dictionary.</param>
    /// <returns>True if any values were passed to the action, false otherwise.</returns>
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

    /// <summary>
    /// This zips together two sorted dictionaries with the given actions.
    /// If the left AND the right has a value for a column, this will call the given action.
    /// </summary>
    /// <param name="left">The left dictionary in the zip. This may be null.</param>
    /// <param name="right">The right dictionary in the zip. This may be null.</param>
    /// <param name="handle">The handle for processing each column found in both dictionary.</param>
    /// <returns>True if any values were passed to the action, false otherwise.</returns>
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

    /// <summary>Serializes the data into a string.</summary>
    /// <returns>The serialized data.</returns>
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

    /// <summary>Deserializes the data from a string.</summary>
    /// <typeparam name="T">The implementation of data to deserialize into.</typeparam>
    /// <param name="factory">The factory for the implementation to deserialize into.</param>
    /// <param name="data">The data to deserialize from.</param>
    /// <returns>The created and populated with the deserialized data.</returns>
    static protected T deserialize<T>(Func<int, int, T> factory, string data)
        where T : Data {
        string[] lines = data.Split('\n');
        (int rows, int columns) = deserializeHeader(lines[0]);
        T d = factory(rows, columns);
        for (int i = 1; i < lines.Length; ++i)
            deserializeRow(d, i - 1, lines[i]);
        return d;
    }

    /// <summary>Deserializes the header information from the serialization.</summary>
    /// <param name="header">The first line of the serialized data.</param>
    /// <returns>The number of rows and columns specified in the header.</returns>
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

    /// <summary>Deserialized a row of data from the serialized string.</summary>
    /// <param name="d">The data to write the row to.</param>
    /// <param name="row">The row being written.</param>
    /// <param name="line">The line of serialized data for the row.</param>
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
