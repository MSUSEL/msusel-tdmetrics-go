using System;
using System.Collections.Generic;
using System.Text;

namespace Constructs.Tooling;

/// <summary>A tool for writing text to represent a complex data structure.</summary>
public class Journal {

    /// <summary>Gets the string for a journal of the given construct.</summary>
    /// <param name="c">The construct to journal.</param>
    /// <returns>The string for the given construct.</returns>
    public static string ToString(IConstruct c) =>
        new Journal().Write(c).ToString();

    private readonly record struct StubMark(IConstruct Construct, bool Short);

    private static readonly HashSet<StubMark> stubbing = [];

    private readonly StringBuilder sb;
    private readonly string indent;

    /// <summary>Indicates the writes to a journal should be short.</summary>
    public readonly bool Short;

    /// <summary>The newline character used in the journal.</summary>
    public const char NewLine = '\n';
    private const string newLineStr = "\n";

    /// <summary>The string that is uses for each indent.</summary>
    public const string IndentText = "  ";

    /// <summary>The default separator to use between items in a list.</summary>
    public const string CommaSeparator = ", ";

    /// <summary>Creates a new journal to write with.</summary>
    public Journal() {
        this.sb = new();
        this.indent = "";
        this.Short = false;
    }

    /// <summary>Creates a journal with the given values.</summary>
    /// <param name="sb">The string builder to write to.</param>
    /// <param name="indent">The indent to use on each line.</param>
    /// <param name="isShort">Indicates if the journal should be short or long.</param>
    private Journal(StringBuilder sb, string indent, bool isShort) {
        this.sb = sb;
        this.indent = indent;
        this.Short = isShort;
    }

    /// <summary>Indicates the writes to a journal should be long.</summary>
    public bool Long => !this.Short;

    /// <summary>Returns a journal that writes to this journal with additional indents.</summary>
    public Journal Indent => new(this.sb, this.indent + IndentText, this.Short);

    /// <summary>Returns a journal that indicates short writes.</summary>
    public Journal AsShort => new(this.sb, this.indent, true);

    /// <summary>Returns a journal that indicates long writes.</summary>
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
        string[] parts = text.ReplaceLineEndings(newLineStr).Split(NewLine);
        int max = parts.Length - 1;
        for (int i = 0; i <= max; ++i)
            this.writeOneLine(parts[i], i == max);
    }

    private void writeStub(IConstruct c) {
        StubMark m = new(c, this.Short);
        try {
            if (!stubbing.Add(m))
                throw new Exception("Journel detected loop in " + c.GetType().FullName + " when " + (this.Short ? "short" : "long") + ".");
            c.ToStub(this);
        } finally {
            stubbing.Remove(m);
        }
    }

    /// <summary>Writes the given text to the journal.</summary>
    /// <param name="text">The text to write to the journal.</param>
    /// <param name="prefix">A prefix to write before the text, if the text is not null nor empty.</param>
    /// <param name="suffix">A suffix to write after the text, if the text is not null nor empty.</param>
    /// <returns>This journal for chaining.</returns>
    public Journal Write(string text, string prefix = "", string suffix = "") {
        if (!string.IsNullOrEmpty(text))
            this.writeText(prefix + text + suffix);
        return this;
    }

    /// <summary>Writes a list of constructs to the journal.</summary>
    /// <typeparam name="T">The type of constructs to write.</typeparam>
    /// <param name="list">The list of constructs to write.</param>
    /// <param name="prefix">The prefix to write before the list, if the list is not empty.</param>
    /// <param name="suffix">The suffix to write after the list, if the list is not empty.</param>
    /// <param name="separator">The text to write between the list items.</param>
    /// <returns>This journal for chaining.</returns>
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

    /// <summary>Writes the given construct to the journal.</summary>
    /// <param name="c">The construct to write.</param>
    /// <param name="prefix">The prefix to write before the construct, if the construct is not null.</param>
    /// <param name="suffix">The suffix to write after the construct, if the construct is not null.</param>
    /// <returns>This journal for chaining.</returns>
    public Journal Write(IConstruct? c, string prefix = "", string suffix = "") {
        if (c is not null) {
            this.writeText(prefix);
            this.writeStub(c);
            this.writeText(suffix);
        }
        return this;
    }

    /// <summary>Writes the given text to the journal, followed by a new line if the text is not null nor empty.</summary>
    /// <param name="text">The text to write to the journal.</param>
    /// <param name="prefix">A prefix to write before the text, if the text is not null nor empty.</param>
    /// <param name="suffix">A suffix to write after the text, if the text is not null nor empty.</param>
    /// <returns>This journal for chaining.</returns>
    public Journal WriteLine(string text, string prefix = "", string suffix = "") =>
        this.Write(text, prefix, suffix + NewLine);

    /// <summary>Writes a list of constructs to the journal, followed by a new line if the list is not empty.</summary>
    /// <typeparam name="T">The type of constructs to write.</typeparam>
    /// <param name="list">The list of constructs to write.</param>
    /// <param name="prefix">The prefix to write before the list, if the list is not empty.</param>
    /// <param name="suffix">The suffix to write after the list, if the list is not empty.</param>
    /// <param name="separator">The text to write between the list items.</param>
    /// <returns>This journal for chaining.</returns>
    public Journal WriteLine<T>(IReadOnlyList<T> list, string prefix = "", string suffix = "", string separator = CommaSeparator)
        where T : IConstruct =>
        this.Write(list, prefix, suffix + NewLine, separator);

    /// <summary>Writes the given construct to the journal, followed by a new line if the construct is not null.</summary>
    /// <param name="c">The construct to write.</param>
    /// <param name="prefix">The prefix to write before the construct, if the construct is not null.</param>
    /// <param name="suffix">The suffix to write after the construct, if the construct is not null.</param>
    /// <returns>This journal for chaining.</returns>
    public Journal WriteLine(IConstruct? c, string prefix = "", string suffix = "") =>
        this.Write(c, prefix, suffix + NewLine);

    /// <summary>Writes a new line to the journal.</summary>
    /// <returns>This journal for chaining.</returns>
    public Journal WriteLine() {
        this.sb.Append(NewLine);
        return this;
    }

    /// <summary>Gets the text that has been written to the journal.</summary>
    /// <returns>The journal's text.</returns>
    public override string ToString() => this.sb.ToString();
}
