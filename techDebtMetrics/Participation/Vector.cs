using System;
using System.Collections.Generic;
using System.Linq;

namespace Participation;

/// <summary>A 1D sparse vector.</summary>
/// <remarks>
/// This vector assumes that the number of non-zero
/// entries is expected to be very small.
/// </remarks>
public class Vector : Data {
    private readonly SortedDictionary<int, double> data;

    /// <summary>Deserialized the given data into a vector.</summary>
    /// <param name="data">The serialized data to populate this vector with.</param>
    /// <param name="epsilon">The epsilon comparator used for determining if a value is zero or not.</param>
    /// <returns>The deserialized vector.</returns>
    static public Vector Deserialize(string data, double epsilon = DefaultEpsilon) =>
        deserialize((rows, columns) => {
            if (columns != 1)
                throw new Exception("Expected the number of columns to be only one for a vector, but got " + columns);
            return new Vector(rows, epsilon);
        }, data);

    /// <summary>Creates a new sparse vector.</summary>
    /// <param name="rows">The number of rows for the vector.</param>
    /// <param name="epsilon">The epsilon comparator used for determining if a value is zero or not.</param>
    public Vector(int rows, double epsilon = DefaultEpsilon) :
        base(rows, 1, epsilon) =>
        this.data = [];

    /// <summary>Creates a new sparse vector.</summary>
    /// <param name="data">The data to populate the vector with.</param>
    /// <param name="epsilon">The epsilon comparator used for determining if a value is zero or not.</param>
    public Vector(double[] data, double epsilon = DefaultEpsilon) :
        this(data.Length, epsilon) {
        for (int row = 0; row < this.Rows; ++row)
            this.SetIfNonZero(row, 0, data[row]);
    }

    /// <summary>Creates a new sparse vector directly given the data to use.</summary>
    /// <param name="data">The data to use in this vector.</param>
    /// <param name="rows">The number of rows for the vector.</param>
    /// <param name="epsilon">The epsilon comparator used for determining if a value is zero or not.</param>
    internal Vector(SortedDictionary<int, double> data, int rows, double epsilon = DefaultEpsilon) :
        base(rows, 1, epsilon) =>
        this.data = data;

    /// <summary>Creates a new sparse vector.</summary>
    /// <param name="rows">The number of rows for the vector.</param>
    /// <param name="entries">The data to populate thr matrix with.</param>
    /// <param name="epsilon">The epsilon comparator used for determining if a value is zero or not.</param>
    public Vector(int rows, IEnumerable<Entry> entries, double epsilon = DefaultEpsilon) :
        this(rows, epsilon) {
        foreach (Entry entry in entries) {
            this.CheckRange(entry.Row, entry.Column);
            this.SetIfNonZero(entry.Row, entry.Column, entry.Value);
        }
    }

    /// <summary>Gets or sets the value at the given row.</summary>
    /// <param name="row">The row to get or set, [0..Rows).</param>
    /// <returns>The value at the given row.</returns>
    public double this[int row] {
        get => this[row, 0];
        set => this[row, 0] = value;
    }

    /// <summary>Gets or sets the value at the given row.</summary>
    /// <param name="row">The row to get or set, [0..Rows).</param>
    /// <returns>The value at the given row.</returns>
    public double this[Index row] {
        get => this[row, 0];
        set => this[row, 0] = value;
    }

    #region Data overrides

    protected override double GetValue(int row, int column) =>
        this.data.TryGetValue(row, out double value) ? value : 0.0;

    protected override void SetValue(int row, int column, double value) =>
        this.data[row] = value;

    protected override bool RemoveValue(int row, int column) =>
        this.data.Remove(row);

    protected override bool ColumnHasZero(int column) =>
        column != 0 || this.data.Count != this.Rows;

    internal override SortedDictionary<int, double> GetColumnNode(int column) =>
        column == 0 ? this.data : [];

    internal override SortedDictionary<int, double> GetRowNode(int row) {
        SortedDictionary<int, double> result = [];
        if (this.data.TryGetValue(row, out double value))
            result.Add(row, value);
        return result;
    }

    public override IEnumerable<Entry> ShortEnumerate() {
        foreach (KeyValuePair<int, double> edge in this.data)
            yield return new(edge.Key, 0, edge.Value);
    }

    public override IEnumerable<Entry> FullEnumerate() {
        int next = 0;
        foreach (KeyValuePair<int, double> edge in this.data) {
            for (int row = next; row < edge.Key; ++row)
                yield return new(row, 0, 0.0);
            yield return new(edge.Key, 0, edge.Value);
            next = edge.Key + 1;
        }
        for (int row = next; row < this.Rows; ++row)
            yield return new(row, 0, 0.0);
    }

    #endregion

    /// <summary>Creates a copy of this vector.</summary>
    /// <returns>The clone of this vector.</returns>
    public Vector Clone() =>
        new(this.Rows, this, this.Epsilon);

    /// <summary>This adds two vectors together.</summary>
    /// <param name="left">The left vector in the sum.</param>
    /// <param name="right">The right vector in the sum.</param>
    /// <returns>The sum of the two vectors.</returns>
    public static Vector operator +(Vector left, Vector right) =>
        overlay(left, right, (leftValue, rightValue) => leftValue + rightValue);

    /// <summary>This subtracts one vector from another.</summary>
    /// <param name="left">The left vector to subtract the right vector from..</param>
    /// <param name="right">The right vector to subtract from the left vector.</param>
    /// <returns>The difference between the vectors.</returns>
    public static Vector operator -(Vector left, Vector right) =>
        overlay(left, right, (leftValue, rightValue) => leftValue - rightValue);

    /// <summary>
    /// This will join two same sized matrices to create a new vector.
    /// The given joiner will be calle for any entry in the left OR right vector that is non-zero.
    /// </summary>
    /// <param name="left">The left vector in the overlay.</param>
    /// <param name="right">The right vector in the overlay.</param>
    /// <param name="joiner">The function to call when either the left OR right value is non-zero.</param>
    /// <returns>The resulting joined vector.</returns>
    private static Vector overlay(Vector left, Vector right, Func<double, double, double> joiner) {
        if (left.Rows != right.Rows)
            throw new Exception("The left (" + left.Rows + ") and " +
                "right (" + right.Rows + ") vectors need to be the same size.");

        Vector result = new(left.Rows, left.Epsilon);
        zipOr(left.data, right.data, (row, leftVal, rightVal) =>
            result.SetIfNonZero(row, 0, joiner(leftVal, rightVal)));
        return result;
    }

    /// <summary>This negates the vector.</summary>
    /// <param name="vector">The vector to negate.</param>
    /// <returns>The negated vector.</returns>
    public static Vector operator -(Vector vector) =>
        new(vector.Rows, vector.Select(e => new Entry(e.Row, e.Column, -e.Value)), vector.Epsilon);

    /// <summary>This scales te vector by a specific value.</summary>
    /// <param name="left">The vector to scale.</param>
    /// <param name="right">The value to scale the vector by.</param>
    /// <returns>The scaled vector.</returns>
    public static Vector operator *(Vector left, double right) =>
        new(left.Rows, left.Select(e => new Entry(e.Row, e.Column, e.Value * right)), left.Epsilon);

    /// <summary>This scales te vector by a specific value.</summary>
    /// <param name="left">The value to scale the vector by.</param>
    /// <param name="right">The vector to scale.</param>
    /// <returns>The scaled vector.</returns>
    public static Vector operator *(double left, Vector right) =>
        new(right.Rows, right.Select(e => new Entry(e.Row, e.Column, left * e.Value)), right.Epsilon);
}
