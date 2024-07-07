using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class Solid : ITypeDesc, IInitializable {
    public ITypeDesc Target => this.inTarget ??
        throw new UninitializedException("target");
    private ITypeDesc? inTarget;

    public IReadOnlyList<ITypeDesc> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<ITypeDesc> inTypeParams = [];

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();
        this.inTarget = obj.ReadIndexType<ITypeDesc>("target", getter);
        obj.ReadIndexTypeList("typeParams", getter, this.inTypeParams);
    }

    public void ToStub(Journal j) =>
        j.Write(this.Target).Write(this.TypeParams, "<", ">");
}
