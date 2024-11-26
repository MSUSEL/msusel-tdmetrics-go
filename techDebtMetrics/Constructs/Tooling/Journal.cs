using System.Collections.Generic;
using System.Text;

namespace Constructs.Tooling;

public class Journal {

    public static string ToString<T>(T c) where T : IConstruct =>
        new Journal().Write(c).ToString();

    private readonly StringBuilder sb;
    private readonly string indent;

    public readonly bool Short;

    public const char NewLine = '\n';
    private const string NewLineStr = "\n";
    public const string IndentText = "  ";
    public const string CommaSeparator = ", ";

    public Journal() {
        this.sb = new();
        this.indent = "";
        this.Short = false;
    }

    private Journal(StringBuilder sb, string indent, bool isShort) {
        this.sb = sb;
        this.indent = indent;
        this.Short = isShort;
    }

    public bool Long => !this.Short;

    public Journal Indent => new(this.sb, this.indent + IndentText, this.Short);

    public Journal AsShort => new(this.sb, this.indent, true);

    public Journal AsLong => new(this.sb, this.indent, false);

    private void writeOneLine(string text, bool last) {
        if (text.Length <= 0) {
            if (!last) this.sb.Append(NewLine);
            return;
        }

        if (this.sb.Length > 0 && this.sb[^1] == NewLine)
            this.sb.Append(this.indent);

        this.sb.Append(text);
        if (!last) this.sb.Append('\n');
    }

    private void writeText(string text) {
        string[] parts = text.ReplaceLineEndings(NewLineStr).Split(NewLine);
        int max = parts.Length - 1;
        for (int i = 0; i <= max; ++i)
            this.writeOneLine(parts[i], i == max);
    }

    private void writeStub<T>(T c)
        where T : IConstruct =>
        c.ToStub(this);

    public Journal Write(string text, string prefix = "", string suffix = "") {
        if (!string.IsNullOrEmpty(text))
            this.writeText(prefix + text + suffix);
        return this;
    }

    public Journal Write<T>(IReadOnlyList<T> list, string prefix = "", string suffix = "", string separator = CommaSeparator)
        where T : IConstruct {
        int count = list.Count;
        if (count > 0) {
            this.writeText(prefix); 
            this.writeStub(list[0]);
            for (int i = 1; i < count; ++i) {
                this.writeText(separator);
                this.writeStub(list[i]);
            }
            this.writeText(suffix);
        }
        return this;
    }

    public Journal Write<T>(T? c, string prefix = "", string suffix = "")
        where T : IConstruct {
        if (c is not null) {
            this.writeText(prefix);
            this.writeStub(c);
            this.writeText(suffix);
        }
        return this;
    }

    public Journal WriteLine(string text, string prefix = "", string suffix = "") =>
        this.Write(text, prefix, suffix + NewLine);

    public Journal WriteLine<T>(IReadOnlyList<T> list, string prefix = "", string suffix = "", string separator = CommaSeparator)
        where T : IConstruct =>
        this.Write(list, prefix, suffix + NewLine, separator);

    public Journal WriteLine<T>(T? c, string prefix = "", string suffix = "")
        where T : IConstruct =>
        this.Write(c, prefix, suffix + NewLine);

    public Journal WriteLine() {
        this.sb.Append(NewLine);
        return this;
    }

    public override string ToString() => this.sb.ToString();
}
