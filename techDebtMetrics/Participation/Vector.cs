using System;
using System.Collections;
using System.Collections.Generic;

namespace Participation;

public class Vector : IEnumerable<Entry> {
    public readonly int Rows;
    public readonly double Epsilon;

    private readonly SortedDictionary<int, double> data;

    public Vector(int rows, double epsilon = 1.0e-9) {
        this.Rows = rows;
        this.Epsilon = epsilon;
        this.data = [];
    }

    public Vector(double[] data, double epsilon = 1.0e-9) {
        this.Rows = data.Length;
        this.Epsilon = epsilon;
        this.data = [];
        for (int row = 0; row < this.Rows; ++row)
            this[row] = data[row];
    }

    public double this[int row] {
        get {
            this.checkRange(row);
            return this.data.TryGetValue(row, out double value) ? value : 0.0;
        }
        set {
            this.checkRange(row);
            if (double.Abs(value) < this.Epsilon)
                this.data.Remove(row);
            else this.data[row] = value;
        }
    }

    internal SortedDictionary<int, double> getDictionary() => this.data;

    private void checkRange(int row) {
        if (row < 0 || row >= this.Rows)
            throw new IndexOutOfRangeException("Row must be in [0.." + this.Rows + "), the given row was " + row);
    }

    public IEnumerator<Entry> GetEnumerator() => this.FullEnumerate().GetEnumerator();
    IEnumerator IEnumerable.GetEnumerator() => this.FullEnumerate().GetEnumerator();

    public IEnumerable<Entry> ShortEnumerate() {
        foreach (KeyValuePair<int, double> edge in this.data)
            yield return new(edge.Key, 0, edge.Value);
    }

    public IEnumerable<Entry> FullEnumerate() {
        int next = 0;
        foreach (KeyValuePair<int, double> edge in this.data) {
            for (int row = 0; row < edge.Key; ++row)
                yield return new(row, 0, 0.0);
            yield return new(edge.Key, 0, edge.Value);
            next = edge.Key + 1;
        }
        for (int row = next; row < this.Rows; ++row)
            yield return new(row, 0, 0.0);
    }
}
