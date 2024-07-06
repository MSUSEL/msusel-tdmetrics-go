using System;
using System.Collections;
using System.Collections.Generic;
using System.Text;

namespace Yamlite.Tokenizer;

internal class Scanner : IEnumerator<char> {
    private readonly List<char> pending = [];
    private readonly List<char> buffer  = [];
    private readonly IEnumerator<char> source;

    public Scanner(string source): this(source.GetEnumerator()) { }

    public Scanner(IEnumerator<char> source) {
        this.source = source;
        this.inReset();
    }

    public char Current { get; private set; } = '\0';
    object IEnumerator.Current => this.Current;

    private bool hasCurrent;
    public Position CurrentPos { get; private set; } = new(0, 0, 1);

    private char start;
    private bool hasStart;
    public Position StartPos { get; private set; } = new(0, 0, 1);

    public int Count => this.buffer.Count;

    private void inReset() {
        this.pending.Clear();
        this.buffer.Clear();

        this.Current    = '\0';
        this.hasCurrent = false;
        this.CurrentPos = new(0, 0, 1);

        this.start    = '\0';
        this.hasStart = false;
        this.StartPos = new(0, 0, 1);
    }

    public void Dispose() {
        this.source.Dispose();
        this.inReset();
    }

    public void Reset() {
        this.source.Reset();
        this.inReset();
    }

    private void stepLocation() {
        if (!this.hasCurrent) return;
        this.CurrentPos = this.CurrentPos.Step(this.Current);
    }

    public string Take(int count) {
        if (count < 0 || count > this.Count)
            throw new ArgumentOutOfRangeException(nameof(count));

        this.Current    = this.start;
        this.hasCurrent = this.hasStart;
        this.CurrentPos = this.StartPos;

        StringBuilder sb = new(count);
        for (int i = 0; i < count; ++i) {
            char c = this.buffer[i];

            this.stepLocation();
            sb.Append(c);
            this.Current = c;
            this.hasCurrent = true;
        }

        this.pending.InsertRange(0, this.buffer[count..]);
        this.buffer.Clear();

        this.start    = this.Current;
        this.hasStart = this.hasCurrent;
        this.StartPos = this.CurrentPos;

        return sb.ToString();
    }

    public bool MoveNext() {
        char c;
        if (pending.Count > 0) {
            c = pending[0];
            pending.RemoveAt(0);
        } else if (source.MoveNext()) c = source.Current;
        else return false;

        this.stepLocation();
        this.buffer.Add(c);
        this.Current = c;
        this.hasCurrent = true;
        return true;
    }

    public override string ToString() => new(this.buffer.ToArray());
}
