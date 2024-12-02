using System;
using System.Collections;
using System.Collections.Generic;
using System.Text;

namespace Participation;

public class Matrix : IEnumerable<Entry> {
    public readonly int Rows;
    public readonly int Columns;
    public readonly double Epsilon;

    private readonly SortedDictionary<int, double>[] data;

    public Matrix(int rows, int columns, double epsilon = 1.0e-9) {
        this.Rows = rows;
        this.Columns = columns;
        this.Epsilon = epsilon;
        this.data = new SortedDictionary<int, double>[rows];
    }

    public Matrix(double[,] data, double epsilon = 1.0e-9) {
        this.Rows = data.GetLength(0);
        this.Columns = data.GetLength(1);
        this.Epsilon = epsilon;
        this.data = new SortedDictionary<int, double>[this.Rows];
        for (int row = 0; row < this.Rows; ++row)
            for (int column = 0; column < this.Columns; ++column)
                this[row, column] = data[row, column];
    }

    public double this[int row, int column] {
        get {
            this.checkRange(row, column);
            return this.getValue(row, column);
        }
        set {
            this.checkRange(row, column);
            if (double.Abs(value) < this.Epsilon)
                this.removeValue(row, column);
            else this.setValue(row, column, value);
        }
    }

    private void checkRange(int row, int column) {
        if (row < 0 || row >= this.Rows)
            throw new IndexOutOfRangeException("Row must be in [0.." + this.Rows + "), the given row was " + row);
        if (column < 0 || column >= this.Columns)
            throw new IndexOutOfRangeException("Column must be in [0.." + this.Columns + "), the given column was " + column);
    }

    private double getValue(int row, int column) {
        SortedDictionary<int, double> node = this.data[row];
        if (node is null) return 0.0;
        node.TryGetValue(column, out double value);
        return value;
    }

    private void setValue(int row, int column, double value) {
        SortedDictionary<int, double> node = this.data[row];
        if (node is null) this.data[row] = node = [];
        node[column] = value;
    }

    private bool removeValue(int row, int column) =>
        this.data[row]?.Remove(column) ?? false;

    #region Enumerators

    public IEnumerator<Entry> GetEnumerator() => this.FullEnumerate().GetEnumerator();
    IEnumerator IEnumerable.GetEnumerator() => this.FullEnumerate().GetEnumerator();

    public IEnumerable<Entry> ShortEnumerate() {
        for (int row = 0; row < this.Rows; ++row)
            foreach (KeyValuePair<int, double> edge in this.data[row])
                yield return new(row, edge.Key, edge.Value);
    }

    public IEnumerable<Entry> FullEnumerate() {
        for (int row = 0; row < this.Rows; ++row) {
            SortedDictionary<int, double> node = this.data[row];
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

    // TODO: Finish
    //public IEnumerator<Entry> RegionEnumerate(int startRow, int startColumn, int rows, int columns)

    #endregion
    #region Serialization

    public string Serialize() {
        StringBuilder sb = new();
        string version = "0"; // currently the only version is "0".
        sb.AppendFormat("{0} {1}x{2}", version, this.Rows, this.Columns);
        foreach (SortedDictionary<int, double> node in this.data) {
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

            if (Math.Abs(value) >= m.Epsilon)
                data[column] = value;
        }
        m.data[row] = data;
    }

    #endregion
    #region Maths

    public static Matrix operator *(Matrix left, Matrix right) {
        if (left.Columns != right.Rows)
            throw new Exception("The left's columns (" + left.Columns + ") must be equal to the right's rows (" + right.Rows + ").");

        Matrix result = new(left.Rows, right.Columns, left.Epsilon);
        for (int row = 0; row < left.Rows; ++row) {
            SortedDictionary<int, double> node = left.data[row];
            if (node is not null) {
                for (int column = 0; column < right.Columns; ++column) {
                    double sum = 0.0;
                    foreach (KeyValuePair<int, double> edge in node)
                        sum += edge.Value * right[edge.Key, row];
                    result[row, column] = sum;
                }
            }
        }
        return result;
    }

    public static Matrix Zip(Matrix left, Matrix right, Func<double, double, double> joiner) {
        if (left.Rows != right.Rows || left.Columns != right.Columns)
            throw new Exception("The left (" + left.Rows + "x" + left.Columns + ") and right (" + right.Rows + "x" + right.Columns + ") matrices need to be the same size.");

        Matrix result = new(left.Rows, left.Columns, left.Epsilon);
        for (int row = 0; row < left.Rows; ++row) {
            SortedDictionary<int, double> leftNode = left.data[row];
            SortedDictionary<int, double> rightNode = right.data[row];

            if (leftNode is not null) {
                if (rightNode is not null) {
                    SortedDictionary<int, double>.Enumerator leftIt = leftNode.GetEnumerator();
                    SortedDictionary<int, double>.Enumerator rightIt = leftNode.GetEnumerator();
                    KeyValuePair<int, double>? leftPair = leftIt.MoveNext() ? leftIt.Current : null;
                    KeyValuePair<int, double>? rightPair = rightIt.MoveNext() ? rightIt.Current : null;

                    while (leftPair is not null && rightPair is not null) {
                        if (leftPair.Value.Key < rightPair.Value.Key) {
                            result[row, leftPair.Value.Key] = joiner(leftPair.Value.Value, 0.0);
                            leftPair = leftIt.MoveNext() ? leftIt.Current : null;

                        } else if (leftPair.Value.Key > rightPair.Value.Key) {
                            result[row, rightPair.Value.Key] = joiner(0.0, rightPair.Value.Value);
                            rightPair = rightIt.MoveNext() ? rightIt.Current : null;

                        } else {
                            result[row, leftPair.Value.Key] = joiner(leftPair.Value.Value, rightPair.Value.Value);
                            leftPair = leftIt.MoveNext() ? leftIt.Current : null;
                            rightPair = rightIt.MoveNext() ? rightIt.Current : null;
                        }
                    }

                    while (leftPair is not null) {
                        result[row, leftPair.Value.Key] = joiner(leftPair.Value.Value, 0.0);
                        leftPair = leftIt.MoveNext() ? leftIt.Current : null;
                    }

                    while (rightPair is not null) {
                        result[row, rightPair.Value.Key] = joiner(0.0, rightPair.Value.Value);
                        rightPair = rightIt.MoveNext() ? rightIt.Current : null;
                    }

                } else {
                    foreach (KeyValuePair<int, double> edge in leftNode)
                        result[row, edge.Key] = joiner(edge.Value, 0.0);
                }

            } else if (rightNode is not null) {
                foreach (KeyValuePair<int, double> edge in rightNode)
                    result[row, edge.Key] = joiner(0.0, edge.Value);
            }
        }
        return result;
    }

    public static Matrix operator +(Matrix left, Matrix right) =>
        Zip(left, right, (leftValue, rightValue) => leftValue + rightValue);

    public static Matrix operator -(Matrix left, Matrix right) =>
        Zip(left, right, (leftValue, rightValue) => leftValue - rightValue);

    public static Matrix operator -(Matrix matrix) {
        Matrix result = new(matrix.Rows, matrix.Columns, matrix.Epsilon);
        foreach (Entry entry in matrix.ShortEnumerate())
            result[entry.Row, entry.Column] = -entry.Value;
        return result;
    }

    public static Matrix operator *(Matrix left, double right) {
        Matrix result = new(left.Rows, left.Columns, left.Epsilon);
        foreach (Entry entry in left.ShortEnumerate())
            result[entry.Row, entry.Column] = entry.Value * right;
        return result;
    }

    public static Matrix operator *(double left, Matrix right) => right * left;

    #endregion

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
