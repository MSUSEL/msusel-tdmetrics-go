using Commons.Data.Yaml;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>A signature type description for a method, function, or function literal.</summary>
/// <see cref="../../docs/genFeatureDef.md#signature"/>
public class Signature : ITypeDesc, IInitializable<Project> {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>Indicates that the last parameter of the signature may have zero or more values.</summary>
    public bool Variadic { get; private set; }

    /// <summary>The list of input parameters for this signature.</summary>
    public IReadOnlyList<Argument> Params => this.inParams.AsReadOnly();
    private readonly List<Argument> inParams = [];

    /// <summary>The list of output results for this signature</summary>
    public IReadOnlyList<Argument> Results => this.inResults.AsReadOnly();
    private readonly List<Argument> inResults = [];

    /// <summary>The list of method declarations and instances that have this signature.</summary>
    public IReadOnlyList<IMethod> Uses => this.inUses.AsReadOnly();
    private readonly List<IMethod> inUses = [];
    internal void AddUses(IMethod use) => this.inUses.Add(use);

    /// <summary>Indicates that this method has no parameters or results.</summary>
    public bool IsEmpty => this.Params.Count <= 0 && this.Results.Count <= 0;

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs {
        get {
            foreach (IConstruct c in this.Params) yield return c;
            foreach (IConstruct c in this.Results) yield return c;
        }
    }

    void IInitializable<Project>.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.Variadic = obj.TryReadBool("variadic");
        obj.TryReadIndexList("params", this.inParams, project.Arguments);
        obj.TryReadIndexList("results", this.inResults, project.Arguments);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write("(").Write(this.Params);
        if (this.Variadic) j.Write("...");
        j.Write(")");
        if (this.Results.Count == 1)
            j.AsShort.Write(this.Results[0].Type, prefix: " ");
        else j.Write(this.Results, "(", ")");
    }
}
