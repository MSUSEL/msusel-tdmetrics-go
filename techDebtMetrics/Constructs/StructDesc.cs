using Constructs.Data;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class StructDesc : ITypeDesc, IInitializable {

    public IReadOnlyList<Field> Fields => this.inFields.AsReadOnly();
    private List<Field> inFields = [];

    public bool IsEmpty => this.Fields.Count <= 0;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inFields = obj.TryReadIndexList("fields", project.Fields);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.WriteLine("{");
        j.Indent.Write(this.Fields, separator: "\n");
        j.Write("\n}");
    }
}
