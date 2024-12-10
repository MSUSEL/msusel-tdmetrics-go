using System;
using System.Collections.Generic;
using System.Linq;

namespace Participation;

public class Matrix : Data {

    private readonly int rows;
    private readonly int columns;
    private readonly double epsilon;

    private readonly SortedDictionary<int, double>[] data;

    static public Matrix Deserialize(string data, double epsilon = 1.0e-9) =>
        deserialize((rows, columns) => new Matrix(rows, columns, epsilon), data);

    public Matrix(int rows, int columns, double epsilon = 1.0e-9) {
        this.rows = rows;
        this.columns = columns;
        this.epsilon = epsilon;
        this.data = new SortedDictionary<int, double>[rows];
    }

    public Matrix(double[,] data, double epsilon = 1.0e-9) {
        this.rows = data.GetLength(0);
        this.columns = data.GetLength(1);
        this.epsilon = epsilon;
        this.data = new SortedDictionary<int, double>[this.Rows];
        for (int row = 0; row < this.Rows; ++row)
            for (int column = 0; column < this.Columns; ++column)
                this[row, column] = data[row, column];
    }

    public override int Rows => this.rows;
    public override int Columns => this.columns;
    public override double Epsilon => this.epsilon;

    protected override double GetValue(int row, int column) {
        SortedDictionary<int, double> node = this.data[row];
        if (node is null) return 0.0;
        return node.TryGetValue(column, out double value) ? value : 0.0;
    }

    protected override void SetValue(int row, int column, double value) {
        SortedDictionary<int, double> node = this.data[row];
        if (node is null) this.data[row] = node = [];
        node[column] = value;
    }

    protected override bool RemoveValue(int row, int column) =>
        this.data[row]?.Remove(column) ?? false;

    public Vector GetColumn(int column) =>
        new(this.GetColumnNode(column), this.Rows, this.Epsilon);

    protected override bool ColumnHasZero(int column) =>
        this.data.Any((node) => !(node?.ContainsKey(column) ?? false));

    public override IEnumerable<Entry> ShortEnumerate() {
        for (int row = 0; row < this.Rows; ++row) {
            SortedDictionary<int, double> node = this.data[row];
            if (node is not null) {
                foreach (KeyValuePair<int, double> edge in node)
                    yield return new(row, edge.Key, edge.Value);
            }
        }
    }

    public override IEnumerable<Entry> FullEnumerate() {
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

    protected override SortedDictionary<int, double> GetColumnNode(int column) {
        SortedDictionary<int, double> result = [];
        for (int row = 0; row < this.Rows; ++row) {
            SortedDictionary<int, double> node = this.data[row];
            if (node is not null && node.TryGetValue(column, out double value))
                result[row] = value;
        }
        return result;
    }

    protected override SortedDictionary<int, double> GetRowNode(int row) =>
        this.data[row];

    public static Matrix operator *(Matrix left, Matrix right) {
        if (left.Columns != right.Rows)
            throw new Exception("The left's columns (" + left.Columns + ") must be equal to the right's rows (" + right.Rows + ").");

        Matrix result = new(left.Rows, right.Columns, left.Epsilon);
        for (int column = 0; column < right.Columns; ++column) {
            SortedDictionary<int, double> rightNode = right.GetColumnNode(column);
            for (int row = 0; row < left.Rows; ++row) {
                double sum = 0.0;
                zipAnd(left.data[row], rightNode, (column, leftValue, rightValue) =>
                    sum += leftValue * rightValue);
                result[row, column] = sum;
            }
        }
        return result;
    }

    public static Vector operator *(Matrix left, Vector right) {
        if (left.Columns != right.Rows)
            throw new Exception("The left's columns (" + left.Columns + ") must be equal to the right's rows (" + right.Rows + ").");

        Vector result = new(left.Rows, left.Epsilon);
        SortedDictionary<int, double> rightNode = right.GetDictionary();
        for (int row = 0; row < left.Rows; ++row) {
            double sum = 0.0;
            zipAnd(left.data[row], rightNode, (column, leftValue, rightValue) =>
                sum += leftValue * rightValue);
            result[row] = sum;
        }
        return result;
    }

    public static Matrix operator +(Matrix left, Matrix right) =>
        overlay(left, right, (leftValue, rightValue) => leftValue + rightValue);

    public static Matrix operator -(Matrix left, Matrix right) =>
        overlay(left, right, (leftValue, rightValue) => leftValue - rightValue);

    public static Matrix operator -(Matrix matrix) =>
        matrix.perEntry((v) => -v);

    public static Matrix operator *(Matrix left, double right) =>
        left.perEntry((v) => v*right);

    public static Matrix operator *(double left, Matrix right) =>
        right.perEntry((v) => v * left);

    private Matrix perEntry(Func<double, double> handle) {
        Matrix result = new(this.Rows, this.Columns, this.Epsilon);
        for (int row = 0; row < this.Rows; ++row) {
            SortedDictionary<int, double> node = this.data[row];
            if (node is not null && node.Count > 0) {
                SortedDictionary<int, double> resultNode = [];
                foreach (KeyValuePair<int, double> pair in node)
                    resultNode[pair.Key] = handle(pair.Value);
                result.data[row] = resultNode;
            }
        }
        return result;
    }

    private static Matrix overlay(Matrix left, Matrix right, Func<double, double, double> joiner) {
        if (left.Rows != right.Rows || left.Columns != right.Columns)
            throw new Exception("The left (" + left.Rows + "x" + left.Columns + ") and right (" + right.Rows + "x" + right.Columns + ") matrices need to be the same size.");

        Matrix result = new(left.Rows, left.Columns, left.Epsilon);
        for (int row = 0; row < left.Rows; ++row) {
            zipOr(left.data[row], right.data[row], (column, leftVal, rightVal) =>
                result[row, column] = joiner(leftVal, rightVal));
        }
        return result;
    }
}
