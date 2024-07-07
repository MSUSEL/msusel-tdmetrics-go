using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Constructs;

public class Signature : ITypeDesc, IInitializable {
    public bool Variadic { get; private set; }

    public IReadOnlyList<Named> Params => this.inParams.AsReadOnly();
    private readonly List<Named> inParams = [];

    public IReadOnlyList<Named> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<Named> inTypeParams = [];

    public ITypeDesc? ReturnType { get; private set; }

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();

        if (obj.Contains("variadic"))
            this.Variadic = obj.ReadBool("variadic");

        obj.ReadIndexTypeList("params", getter, this.inParams);
        obj.ReadIndexTypeList("typeParams", getter, this.inTypeParams);

        if (obj.Contains("return"))
            this.ReturnType = obj.ReadIndexType<ITypeDesc>("return", getter);
    }

    public string ToStub() {
        StringBuilder sb = new();
        if (this.TypeParams.Count > 0) {
            sb.Append('<');
            sb.Append(this.TypeParams.Select(tp => tp.ToStub()).Join());
            sb.Append('>');
        }

        sb.Append('(');
        sb.Append(this.Params.Select(p => p.ToStub()).Join());
        if (this.Variadic) sb.Append(" ...");
        sb.Append(')');

        if (this.ReturnType is not null) {
            sb.Append(" => ");
            sb.Append(this.ReturnType.ToStub());
        }
        return sb.ToString();
    }
}
