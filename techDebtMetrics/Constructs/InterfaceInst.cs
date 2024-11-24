using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class InterfaceInst : ITypeDesc, IInitializable {

    public InterfaceDecl Generic => this.inGeneric ??
        throw new UninitializedException("generic");
    private InterfaceDecl? inGeneric;

    public IReadOnlyList<ITypeDesc> InstanceTypes => this.inInstanceTypes.AsReadOnly();
    private List<ITypeDesc> inInstanceTypes = [];

    public InterfaceDesc Resolved => this.inResolved ??
        throw new UninitializedException("resolved");
    private InterfaceDesc? inResolved;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inGeneric       = obj.ReadIndex("generics", project.InterfaceDecls);
        this.inInstanceTypes = obj.ReadKeyList<ITypeDesc>("instances", project);
        this.inResolved      = obj.ReadIndex("resolved", project.InterfaceDescs);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Generic).
            AsLong.Write(this.InstanceTypes, "<", ">").
            AsShort.Write(this.Resolved);
}
