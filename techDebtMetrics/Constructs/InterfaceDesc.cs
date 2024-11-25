using Constructs.Data;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class InterfaceDesc : ITypeDesc, IInitializable {

    public IReadOnlyList<Abstract> Abstracts => this.inAbstracts.AsReadOnly();
    private List<Abstract> inAbstracts = [];

    public IReadOnlyList<ITypeDesc> Approx => this.inApprox.AsReadOnly();
    private List<ITypeDesc> inApprox = [];

    public IReadOnlyList<ITypeDesc> Exact => this.inExact.AsReadOnly();
    private List<ITypeDesc> inExact = [];

    public IReadOnlyList<InterfaceDesc> Inherits => this.inInherits.AsReadOnly();
    private List<InterfaceDesc> inInherits = [];

    public bool IsEmpty =>
        this.Abstracts.Count <= 0 && this.Approx.Count <= 0 &&
        this.Exact.Count <= 0 && this.Inherits.Count <= 0;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inAbstracts = obj.TryReadIndexList("abstracts", project.Abstracts);
        this.inApprox    = obj.TryReadKeyList<ITypeDesc>("approx", project);
        this.inExact     = obj.TryReadKeyList<ITypeDesc>("exact", project);
        this.inInherits  = obj.TryReadIndexList("inherits", project.InterfaceDescs);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (this.IsEmpty) {
            j.Write("{}");
            return;
        }

        j.WriteLine("{");
        Journal j2 = j.Indent.AsShort;
        j2.WriteLine(this.Inherits, prefix: "implements: ", suffix: ";");
        if (this.Exact.Count > 0 || this.Approx.Count > 0) {
            j2.Write(this.Exact, separator: "|");
            if (this.Exact.Count > 0 && this.Approx.Count > 0)
                j2.Write("|");
            j2.Write(this.Exact, prefix: "~", separator: "|~");
            j2.WriteLine(";");
        }
        j2.AsLong.WriteLine(this.Abstracts, suffix: ";", separator: ";\n");
        j.Write("}");
    }
}
