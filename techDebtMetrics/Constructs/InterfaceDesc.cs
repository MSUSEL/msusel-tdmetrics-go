using Constructs.Data;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>An interface type description.</summary>
/// <see cref="../../docs/genFeatureDef.md#interface-description"/>
public class InterfaceDesc : ITypeDesc, IInitializable {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The list of abstracts, named function signatures, for this interface.</summary>
    public IReadOnlyList<Abstract> Abstracts => this.inAbstracts.AsReadOnly();
    private readonly List<Abstract> inAbstracts = [];

    /// <summary>The list of approximate type constraints.</summary>
    public IReadOnlyList<ITypeDesc> Approx => this.inApprox.AsReadOnly();
    private readonly List<ITypeDesc> inApprox = [];

    /// <summary>The list of exact type constraints.</summary>
    public IReadOnlyList<ITypeDesc> Exact => this.inExact.AsReadOnly();
    private readonly List<ITypeDesc> inExact = [];

    /// <summary>The interfaces that this interface inherits.</summary>
    public IReadOnlyList<InterfaceDesc> Inherits => this.inInherits.AsReadOnly();
    private readonly List<InterfaceDesc> inInherits = [];

    /// <summary>The list of interface declarations and instances that have this interface.</summary>
    public IReadOnlyList<IInterface> Uses => this.inUses.AsReadOnly();
    private readonly List<IInterface> inUses = [];
    internal void AddUses(IInterface use) => this.inUses.Add(use);

    /// <summary>Indicates this interface is the base type of all other types.</summary>
    public bool IsEmpty =>
        this.Abstracts.Count <= 0 && this.Approx.Count <= 0 &&
        this.Exact.Count <= 0 && this.Inherits.Count <= 0;

    void IInitializable.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        obj.TryReadIndexList("abstracts", this.inAbstracts, project.Abstracts);
        obj.TryReadKeyList("approx", this.inApprox, project);
        obj.TryReadKeyList("exact", this.inExact, project);
        obj.TryReadIndexList("inherits", this.inInherits, project.InterfaceDescs);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (j.Short && this.Uses.Count > 0) {
            j.AsShort.Write(this.Uses[0]);
            return;
        }

        if (this.IsEmpty) {
            j.Write("any");
            return;
        }

        bool hasAbs    = this.Abstracts.Count > 0;
        bool hasExact  = this.Exact.Count > 0;
        bool hasApprox = this.Approx.Count > 0;

        if (hasAbs) j.WriteLine("{");
        Journal j2 = j.Indent.AsShort;
        if (j.Long)
            j2.WriteLine(this.Inherits, prefix: "implements: ", suffix: ";");
        if (hasExact || hasApprox) {
            j2.Write(this.Exact, separator: "|");
            if (hasExact && hasApprox)
                j2.Write("|");
            j2.Write(this.Approx, prefix: "~", separator: "|~");
            if (hasAbs) j2.WriteLine(";");
        }
        if (hasAbs) {
            j2.AsLong.WriteLine(this.Abstracts, suffix: ";", separator: ";\n");
            j.Write("}");
        }
    }
}
