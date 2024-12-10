using System;
using System.Collections.Generic;

namespace Participation;

public class Vector : Data {

    public readonly int rows;
    public readonly double epsilon;

    private readonly SortedDictionary<int, double> data;

    static public Vector Deserialize(string data, double epsilon = 1.0e-9) =>
        deserialize((rows, columns) => {
            if (columns != 1)
                throw new Exception("Expected the number of columns to be only one for a vector, but got " + columns);
            return new Vector(rows, epsilon);
        }, data);

    public Vector(int rows, double epsilon = 1.0e-9) {
        this.rows = rows;
        this.epsilon = epsilon;
        this.data = [];
    }

    public Vector(double[] data, double epsilon = 1.0e-9) {
        this.rows = data.Length;
        this.epsilon = epsilon;
        this.data = [];
        for (int row = 0; row < this.Rows; ++row)
            this[row] = data[row];
    }

    internal Vector(SortedDictionary<int, double> data, int rows, double epsilon) {
        this.rows = rows;
        this.epsilon = epsilon;
        this.data = data;
    }

    public override int Rows => this.rows;
    public override int Columns => 1;
    public override double Epsilon => this.epsilon;

    public double this[int row] {
        get => this[row, 0];
        set => this[row, 0] = value;
    }

    protected override double GetValue(int row, int column) =>
        this.data.TryGetValue(row, out double value) ? value : 0.0;

    protected override void SetValue(int row, int column, double value) =>
        this.data[row] = value;

    protected override bool RemoveValue(int row, int column) =>
        this.data.Remove(row);

    protected override bool ColumnHasZero(int column) =>
        column != 0 || this.data.Count != this.Rows;

    protected override SortedDictionary<int, double> GetColumnNode(int column) =>
        column == 0 ? this.data : [];

    protected override SortedDictionary<int, double> GetRowNode(int row) {
        SortedDictionary<int, double> result = [];
        if (this.data.TryGetValue(row, out double value))
            result.Add(row, value);
        return result;
    }

    internal SortedDictionary<int, double> GetDictionary() => this.data;

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
}
