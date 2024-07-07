using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class Struct : ITypeDesc, IInitializable {
    public IReadOnlyList<Named> Fields => this.inFields.AsReadOnly();
    private readonly List<Named> inFields = [];

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();
        obj.ReadIndexTypeList("fields", getter, this.inFields);
    }

    public void ToStub(Journal j) {
        j.Write("struct {");
        j.Indent.AsLong.Write(this.Fields, "\n", ";\n", ";\n");
        j.Write("}");
    }
}
