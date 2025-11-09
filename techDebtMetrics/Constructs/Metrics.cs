using Commons.Data.Locations;
using Commons.Data.Yaml;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>
/// The metrics gathered from code expressions for either a method
/// or the initializer to a package level variable or constant.
/// </summary>
/// <remarks>
/// Several variables or constants may have the same metrics, determinable via the location,
/// if the initialization expressions were used to initialize multiple at a time,
/// e.g. `var A, B = func() (int, int) { return 13, 45 }()`
/// </remarks>
/// <see cref="../../docs/genFeatureDef.md#metrics"/>
public class Metrics : IConstruct, IInitializable<Project> {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

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
    public bool Getter { get; private set; }

    /// <summary>True indicates the method is a setter pattern.</summary>
    public bool Setter { get; private set; }

    /// <summary>List of methods that were invoked in the method.</summary>
    public IReadOnlyList<IConstruct> Invokess => this.inInvokess.AsReadOnly();
    private readonly List<IConstruct> inInvokess = [];

    /// <summary>List of types that were read from in the method.</summary>
    public IReadOnlyList<IConstruct> Reads => this.inReads.AsReadOnly();
    private readonly List<IConstruct> inReads = [];

    /// <summary>List of types that were written to in the method.</summary>
    public IReadOnlyList<IConstruct> Writes => this.inWrites.AsReadOnly();
    private readonly List<IConstruct> inWrites = [];

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs => [];

    void IInitializable<Project>.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.Location = obj.ReadLocation(project.Locations, "loc");
        this.CodeCount = obj.TryReadInt("codeCount");
        this.Complexity = obj.TryReadInt("complexity");
        this.Indents = obj.TryReadInt("indents");
        this.LineCount = obj.TryReadInt("lineCounr");
        this.Getter = obj.TryReadBool("getter");
        this.Setter = obj.TryReadBool("setter");
        obj.TryReadKeyList(project, "invokes", this.inInvokess);
        obj.TryReadKeyList(project, "reads", this.inReads);
        obj.TryReadKeyList(project, "writes", this.inWrites);
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
