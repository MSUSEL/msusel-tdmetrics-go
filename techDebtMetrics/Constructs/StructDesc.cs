using Constructs.Data;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class StructDesc : ITypeDesc, IInitializable {

    public IReadOnlyList<Field> Fields => this.inFields.AsReadOnly();
    private List<Field> inFields = [];

    public IReadOnlyList<IObject> Uses => this.inUses.AsReadOnly();
    private List<IObject> inUses = [];
    internal void AddUses(IObject use) => this.inUses.Add(use);

    public bool IsEmpty => this.Fields.Count <= 0;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inFields = obj.TryReadIndexList("fields", project.Fields);
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
