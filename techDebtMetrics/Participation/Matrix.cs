using System;
using System.Collections.Generic;
using System.Linq;

namespace Participation;

/// <summary>A 2D sparse matrix.</summary>
/// <remarks>
/// This matrix assumes that there will likely be a value in every row and every column
/// however the number of non-zero entries is expected to be very small.
/// </remarks>
public class Matrix : Data {
    private readonly SortedDictionary<int, double>[] data;

    /// <summary>Deserialized the given data into a matrix.</summary>
    /// <param name="data">The serialized data to populate this matrix with.</param>
    /// <param name="epsilon">The epsilon comparitor used for determining if a value is zero or not.</param>
    /// <returns>The deserialized matrix.</returns>
    static public Matrix Deserialize(string data, double epsilon = DefaultEpsilon) =>
        deserialize((rows, columns) => new Matrix(rows, columns, epsilon), data);

    /// <summary>Creates a new sparse matrix.</summary>
    /// <param name="rows">The number of rows for the matrix.</param>
    /// <param name="columns">The number of columns for the matrix.</param>
    /// <param name="epsilon">The epsilon comparitor used for determining if a value is zero or not.</param>
    public Matrix(int rows, int columns, double epsilon = DefaultEpsilon) :
        base(rows, columns, epsilon) =>
        this.data = new SortedDictionary<int, double>[rows];

    /// <summary>Creates a sparse matrix populated with the given data.</summary>
    /// <param name="data">The data to populate the matrix with.</param>
    /// <param name="epsilon">The epsilon comparitor used for determining if a value is zero or not.</param>
    public Matrix(double[,] data, double epsilon = DefaultEpsilon) :
        this(data.GetLength(0), data.GetLength(1), epsilon) {
        for (int row = 0; row < this.Rows; ++row)
            for (int column = 0; column < this.Columns; ++column)
                this.SetIfNonZero(row, column, data[row, column]);
    }

    /// <summary>Creates a new sparse matrix.</summary>
    /// <param name="rows">The number of rows for the matrix.</param>
    /// <param name="columns">The number of columns for the matrix.</param>
    /// <param name="entries">The data to populate thr matrix with.</param>
    /// <param name="epsilon">The epsilon comparitor used for determining if a value is zero or not.</param>
    public Matrix(int rows, int columns, IEnumerable<Entry> entries, double epsilon = DefaultEpsilon) :
        this(rows, columns, epsilon) {
        foreach (Entry entry in entries) {
            this.CheckRange(entry.Row, entry.Column);
            this.SetIfNonZero(entry.Row, entry.Column, entry.Value);
        }
    }

    #region Data overrides

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

    internal override SortedDictionary<int, double> GetColumnNode(int column) {
        SortedDictionary<int, double> result = [];
        for (int row = 0; row < this.Rows; ++row) {
            SortedDictionary<int, double> node = this.data[row];
            if (node is not null && node.TryGetValue(column, out double value))
                result[row] = value;
        }
        return result;
    }

    internal override SortedDictionary<int, double> GetRowNode(int row) =>
        this.data[row];

    #endregion

    /// <summary>Gets a column of the matrix as a vector.</summary>
    /// <param name="column">The column to get.</param>
    /// <returns>The vector for the column.</returns>
    public Vector GetColumn(int column) =>
        new(this.GetColumnNode(column), this.Rows, this.Epsilon);

    /// <summary>Creates a copy of this matrix.</summary>
    /// <returns>The clone of this matrix.</returns>
    public Matrix Clone() =>
        new(this.Rows, this.Columns, this, this.Epsilon);

    /// <summary>Creates a transposed matrix of this matrix.</summary>
    /// <returns>The transposed matrix.</returns>
    public Matrix Transpose() =>
        new(this.Columns, this.Rows, this.Select(e => new Entry(e.Column, e.Row, e.Value)), this.Epsilon);

    /// <summary>This multiplies two matrices together.</summary>
    /// <param name="left">The left matrix in the multiplication.</param>
    /// <param name="right">The right matrix in the multiplication.</param>
    /// <returns>The result matrix of the multiplication.</returns>
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
                result.SetIfNonZero(row, column, sum);
            }
        }
        return result;
    }

    /// <summary>This multiplies a matrix with a vector.</summary>
    /// <param name="left">The matrix in the multiplication.</param>
    /// <param name="right">The vector in the multiplication.</param>
    /// <returns>The result vector of the multiplication..</returns>
    public static Vector operator *(Matrix left, Vector right) {
        if (left.Columns != right.Rows)
            throw new Exception("The left's columns (" + left.Columns + ") must be equal to the right's rows (" + right.Rows + ").");

        Vector result = new(left.Rows, left.Epsilon);
        SortedDictionary<int, double> rightNode = right.GetColumnNode(0);
        for (int row = 0; row < left.Rows; ++row) {
            double sum = 0.0;
            zipAnd(left.data[row], rightNode, (column, leftValue, rightValue) =>
                sum += leftValue * rightValue);
            result.SetIfNonZero(row, 0, sum);
        }
        return result;
    }

    /// <summary>This adds two matrices together.</summary>
    /// <param name="left">The left matrix in the sum.</param>
    /// <param name="right">The right matrix in the sum.</param>
    /// <returns>The sum of the two matrices.</returns>
    public static Matrix operator +(Matrix left, Matrix right) =>
        overlay(left, right, (leftValue, rightValue) => leftValue + rightValue);

    /// <summary>This subtracts one matrix from another.</summary>
    /// <param name="left">The left matrix to subtract the right matrix from..</param>
    /// <param name="right">The right matrix to subtract from the left matrix.</param>
    /// <returns>The difference between the matrices.</returns>
    public static Matrix operator -(Matrix left, Matrix right) =>
        overlay(left, right, (leftValue, rightValue) => leftValue - rightValue);

    /// <summary>
    /// This will join two same sized matrices to create a new matrix.
    /// The given joiner will be calle for any entry in the left OR right matrix that is non-zero.
    /// </summary>
    /// <param name="left">The left matrix in the overlay.</param>
    /// <param name="right">The right matrix in the overlay.</param>
    /// <param name="joiner">The function to call when either the left OR right value is non-zero.</param>
    /// <returns>The resulting joined matrix.</returns>
    private static Matrix overlay(Matrix left, Matrix right, Func<double, double, double> joiner) {
        if (left.Rows != right.Rows || left.Columns != right.Columns)
            throw new Exception("The left (" + left.Rows + "x" + left.Columns + ") and " +
                "right (" + right.Rows + "x" + right.Columns + ") matrices need to be the same size.");

        Matrix result = new(left.Rows, left.Columns, left.Epsilon);
        for (int row = 0; row < left.Rows; ++row) {
            zipOr(left.data[row], right.data[row], (column, leftVal, rightVal) =>
                result.SetIfNonZero(row, column, joiner(leftVal, rightVal)));
        }
        return result;
    }

    /// <summary>This negates the matrix.</summary>
    /// <param name="matrix">The matrix to negate.</param>
    /// <returns>The negated matrix.</returns>
    public static Matrix operator -(Matrix matrix) =>
        new(matrix.Rows, matrix.Columns, matrix.Select(e => new Entry(e.Row, e.Column, -e.Value)), matrix.Epsilon);

    /// <summary>This scales te matrix by a specific value.</summary>
    /// <param name="left">The matrix to scale.</param>
    /// <param name="right">The value to scale the matrix by.</param>
    /// <returns>The scaled matrix.</returns>
    public static Matrix operator *(Matrix left, double right) =>
        new(left.Rows, left.Columns, left.Select(e => new Entry(e.Row, e.Column, e.Value*right)), left.Epsilon);

    /// <summary>This scales te matrix by a specific value.</summary>
    /// <param name="left">The value to scale the matrix by.</param>
    /// <param name="right">The matrix to scale.</param>
    /// <returns>The scaled matrix.</returns>
    public static Matrix operator *(double left, Matrix right) =>
        new(right.Rows, right.Columns, right.Select(e => new Entry(e.Row, e.Column, left*e.Value)), right.Epsilon);
}
