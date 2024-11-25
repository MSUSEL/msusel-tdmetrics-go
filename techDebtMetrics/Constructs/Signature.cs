using Constructs.Data;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class Signature : ITypeDesc, IInitializable {

    public bool Variadic { get; private set; }

    public IReadOnlyList<Argument> Params => this.inParams.AsReadOnly();
    private List<Argument> inParams = [];

    public IReadOnlyList<Argument> Results => this.inResults.AsReadOnly();
    private List<Argument> inResults = [];

    public bool IsEmpty => this.Variadic && this.Params.Count <= 0 && this.Results.Count <= 0;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Variadic  = obj.TryReadBool("variadic");
        this.inParams  = obj.TryReadIndexList("params", project.Arguments);
        this.inResults = obj.TryReadIndexList("results", project.Arguments);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write("(").Write(this.Params);
        if (this.Variadic) j.Write("...");
        j.Write(")").Write(this.Results, "(", ")");
    }
}
