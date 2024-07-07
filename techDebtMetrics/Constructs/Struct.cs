using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Text;

namespace Constructs;

public class Struct : ITypeDesc, IInitializable {
    public IReadOnlyList<Named> Fields => this.inFields.AsReadOnly();
    private readonly List<Named> inFields = [];

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();
        obj.ReadIndexTypeList("fields", getter, this.inFields);
    }

    public string ToStub() {
        StringBuilder sb = new();
        sb.Append("struct {");
        if (this.inFields.Count > 0) {
            sb.AppendLine();
            foreach (Named field in this.inFields) {
                sb.Append("   ");
                sb.Append(field.ToStub().Indent());
                sb.AppendLine(";");
            }
        }
        sb.Append('}');
        return sb.ToString();
    }
}
