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

    /// <summary>Optional construct that this interface is pinned to.</summary>
    /// <remarks>Interfaces can be pinned to a construct when they have non-exported methods.</remarks>
    public IConstruct? Pin { get; private set; } = null;

    /// <summary>Indicates this interface is the base type of all other types.</summary>
    public bool IsEmpty =>
        this.Abstracts.Count <= 0 && this.Approx.Count <= 0 &&
        this.Exact.Count <= 0 && this.Inherits.Count <= 0;

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs {
        get {
            foreach (var c in this.Abstracts) yield return c;
            foreach (var c in this.Approx) yield return c;
            foreach (var c in this.Exact) yield return c;
        }
    }

    void IInitializable.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        obj.TryReadIndexList("abstracts", this.inAbstracts, project.Abstracts);
        obj.TryReadKeyList("approx", this.inApprox, project);
        obj.TryReadKeyList("exact", this.inExact, project);
        obj.TryReadIndexList("inherits", this.inInherits, project.InterfaceDescs);
        this.Pin = obj.TryReadKey<IConstruct>("pin", project);
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

        bool hasAbs       = this.Abstracts.Count > 0;
        bool hasExact     = this.Exact.Count > 0;
        bool hasApprox    = this.Approx.Count > 0;
        bool hasInherits  = this.Inherits.Count > 0;
        bool showInherits = j.Long && hasInherits;
        bool showParens   = showInherits || hasAbs;
        Journal j2 = j;

        if (showParens) {
            j.WriteLine("{");
            j2 = j.Indent.AsShort;
        }
        if (showInherits)
            j2.WriteLine(this.Inherits, prefix: "implements: ", suffix: ";");
        if (!hasAbs && !hasExact && !hasApprox) {
            if (showParens) j2.WriteLine("any");
            else j2.Write("any");
        }
        if (hasExact || hasApprox) {
            j2.Write(this.Exact, separator: "|");
            if (hasExact && hasApprox)
                j2.Write("|");
            j2.Write(this.Approx, prefix: "~", separator: "|~");
            if (hasAbs) j2.WriteLine(";");
        }
        j2.AsLong.WriteLine(this.Abstracts, suffix: ";", separator: ";\n");
        if (showParens) j.Write("}");
    }
}
