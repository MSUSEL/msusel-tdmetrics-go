using Commons.Data.Yaml;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>A structure type description.</summary>
/// <see cref="../../docs/genFeatureDef.md#structure-description"/>
public class StructDesc : ITypeDesc, IInitializable<Project> {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The list of fields for this structure.</summary>
    public IReadOnlyList<Field> Fields => this.inFields.AsReadOnly();
    private readonly List<Field> inFields = [];

    /// <summary>The list of object declarations and instances that have this structure.</summary>
    public IReadOnlyList<IObject> Uses => this.inUses.AsReadOnly();
    private readonly List<IObject> inUses = [];
    internal void AddUses(IObject use) => this.inUses.Add(use);

    /// <summary>Indicates if this structure has no fields.</summary>
    public bool IsEmpty => this.Fields.Count <= 0;

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs => this.Fields;

    void IInitializable<Project>.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        obj.TryReadIndexList("fields", this.inFields, project.Fields);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (j.Short && this.Uses.Count > 0) {
            j.AsShort.Write(this.Uses[0]);
            return;
        }

        if (this.IsEmpty) {
            j.Write("object");
            return;
        }

        j.WriteLine("{");
        j.Indent.WriteLine(this.Fields, separator: "\n");
        j.Write("}");
    }
}
