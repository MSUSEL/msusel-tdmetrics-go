using Constructs.Data;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class Metrics : IConstruct, IInitializable {
    /// <summary>The location this metrics was read from.</summary>
    public Location Location { get; private set; }

    /// <summary>The number of lines in the method that are not comments or empty.</summary>
    public int CodeCount { get; private set; }

    /// <summary>The cyclomatic complexity of the method.</summary>
    public int Complexity { get; private set; }

    /// <summary>The indent complexity of the method.</summary>
    public int Indents { get; private set; }

    /// <summary>The number of lines in the method.</summary>
    public int LineCount { get; private set; }

    /// <summary>True indicates the method is a getter pattern.</summary>
    public bool Getter {  get; private set; }

    /// <summary>True indicates the method is a setter pattern.</summary>
    public bool Setter { get; private set; }

    /// <summary>List of methods that were invoked in the method.</summary>
    public IReadOnlyList<IConstruct> Invokess => this.inInvokess.AsReadOnly();
    private List<IConstruct> inInvokess = [];

    /// <summary>List of types that were read from in the method.</summary>
    public IReadOnlyList<IConstruct> Reads => this.inReads.AsReadOnly();
    private List<IConstruct> inReads = [];

    /// <summary>List of types that were written to in the method.</summary>
    public IReadOnlyList<IConstruct> Writes => this.inWrites.AsReadOnly();
    private List<IConstruct> inWrites = [];

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Location   = obj.ReadLocation("loc", project);
        this.CodeCount  = obj.TryReadInt("codeCount");
        this.Complexity = obj.TryReadInt("complexity");
        this.Indents    = obj.TryReadInt("indents");
        this.LineCount  = obj.TryReadInt("lineCounr");
        this.Getter     = obj.TryReadBool("getter");
        this.Setter     = obj.TryReadBool("setter");
        this.inInvokess = obj.TryReadKeyList<IConstruct>("invokes", project);
        this.inReads    = obj.TryReadKeyList<IConstruct>("reads", project);
        this.inWrites   = obj.TryReadKeyList<IConstruct>("writes", project);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.WriteLine("{");
        j.Indent.
            WriteLine("codeCount:  " + this.CodeCount).
            WriteLine("complexity: " + this.Complexity).
            WriteLine("indents:    " + this.Indents).
            WriteLine("lineCount:  " + this.LineCount);
        if (this.Getter) j.Indent.WriteLine("getter");
        if (this.Setter) j.Indent.WriteLine("setter");
        j.WriteLine("}");
    }
}
