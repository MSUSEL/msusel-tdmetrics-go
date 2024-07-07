using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Linq;

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

    public void ToStub(Journal j) {
        j.AsLong.Write(this.TypeParams, "<", ">");

        j.Write("(");
        j.AsLong.Write(this.Params);
        if (this.Variadic) j.Write(" ...");
        j.Write(")");

        if (this.ReturnType is not null)
            j.Write(" => ").Write(this.ReturnType);
    }
}
