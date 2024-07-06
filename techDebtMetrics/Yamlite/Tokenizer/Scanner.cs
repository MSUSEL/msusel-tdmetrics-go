using System;
using System.Collections;
using System.Collections.Generic;
using System.Text;

namespace Yamlite.Tokenizer;

internal class Scanner : IEnumerator<char> {
    private readonly List<char> pending = [];
    private readonly List<char> buffer  = [];
    private readonly IEnumerator<char> source;

    public Scanner(IEnumerator<char> source) {
        this.source = source;
        this.inReset();
    }

    public char Current { get; private set; } = '\0';
    object IEnumerator.Current => this.Current;

    private bool hasStart;
    public int StartOffset { get; private set; }
    public int StartColumn { get; private set; }
    public int StartLine { get; private set; }
    
    private bool hasCurrent;
    public int CurrentOffset { get; private set; }
    public int CurrentColumn { get; private set; }
    public int CurrentLine { get; private set; }

    public int Count => this.buffer.Count;

    private void inReset() {
        this.pending.Clear();
        this.buffer.Clear();
        this.Current = '\0';

        this.hasStart = false;
        this.StartOffset = 0;
        this.StartColumn = 0;
        this.StartLine   = 1;
        
        this.hasCurrent = false;
        this.CurrentOffset = 0;
        this.CurrentColumn = 0;
        this.CurrentLine   = 1;
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

        this.CurrentOffset++;
        this.CurrentColumn++;

        if (this.Current == '\n') {
            this.CurrentColumn = 0;
            this.CurrentLine++;
        }
    }

    public string Take(int count) {
        if (count < 0 || count > this.Count)
            throw new ArgumentOutOfRangeException(nameof(count));
        
        this.hasCurrent    = this.hasStart;
        this.CurrentOffset = this.StartOffset;
        this.CurrentColumn = this.StartColumn;
        this.CurrentLine   = this.StartLine;
        this.Current       = '\0';

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

        this.StartOffset = this.CurrentOffset;
        this.StartColumn = this.CurrentColumn;
        this.StartLine   = this.CurrentLine;
        this.hasStart    = this.hasCurrent;

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
