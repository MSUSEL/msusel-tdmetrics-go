using System.Collections.Generic;
using System.Text;

namespace Constructs.Tooling;

public class Journal {
    private readonly StringBuilder sb;
    private readonly string indent;

    public readonly bool Short;

    public Journal() {
        this.sb = new();
        this.indent = "";
        this.Short = false;
    }

    private Journal(StringBuilder sb, string indent, bool Short) {
        this.sb = sb;
        this.indent = indent;
        this.Short = Short;
    }

    public Journal Indent => new(this.sb, this.indent+"   ", this.Short);

    public Journal AsShort => new(this.sb, this.indent, true);

    public Journal AsLong => new(this.sb, this.indent, false);

    public Journal Write(string text) {
        if (this.sb.Length > 0 && this.sb[^1] == '\n')
            this.sb.Append(this.indent);
        this.sb.Append(text.ReplaceLineEndings("\n"+this.indent));
        return this;
    }
    
    public Journal Write<T>(IReadOnlyList<T> list, string prefix = "", string suffix = "", string separator = ", ")
        where T: IConstruct {
        int count = list.Count;
        if (count > 0) {
            this.Write(prefix);
            list[0].ToStub(this);
            for (int i = 1; i < count; ++i) {
                this.Write(", ");
                list[i].ToStub(this);
            }
            this.Write(suffix);
        }
        return this;
    }

    public Journal Write<T>(T c) where T: IConstruct {
        c.ToStub(this);
        return this;
    }

    public Journal WriteLine(string text) =>
        this.Write(text).WriteLine();

    public Journal WriteLine() {
        this.sb.Append('\n');
        return this;
    }

    public override string ToString() => this.sb.ToString();
}
