using Constructs.Data;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class Signature : ITypeDesc, IInitializable {

    public bool Variadic { get; private set; }

    public IReadOnlyList<Argument> Params => this.inParams.AsReadOnly();
    private readonly List<Argument> inParams = [];

    public IReadOnlyList<Argument> Results => this.inResults.AsReadOnly();
    private readonly List<Argument> inResults = [];

    public IReadOnlyList<IMethod> Uses => this.inUses.AsReadOnly();
    private readonly List<IMethod> inUses = [];
    internal void AddUses(IMethod use) => this.inUses.Add(use);

    public bool IsEmpty => this.Variadic && this.Params.Count <= 0 && this.Results.Count <= 0;

    void IInitializable.Initialize(Project project, Node node) {
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
