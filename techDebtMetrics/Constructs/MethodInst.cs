using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class MethodInst : IConstruct, IInitializable {

    public Method Generic => this.inGeneric ??
        throw new UninitializedException("generic");
    private Method? inGeneric;

    public IReadOnlyList<ITypeDesc> InstanceTypes => this.inInstanceTypes.AsReadOnly();
    private List<ITypeDesc> inInstanceTypes = [];

    public ObjectInst? Receiver { get; private set; }

    public Signature Resolved => this.inResolved ??
        throw new UninitializedException("resolved");
    private Signature? inResolved;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inGeneric       = obj.ReadIndex("generics", project.Methods);
        this.inInstanceTypes = obj.ReadKeyList<ITypeDesc>("instances", project);
        this.Receiver        = obj.TryReadIndex("receiver", project.ObjectInsts);
        this.inResolved      = obj.ReadIndex("resolved", project.Signatures);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Receiver, suffix: ".").
            AsShort.Write(this.Generic).
            AsLong.Write(this.InstanceTypes, "<", ">").
            AsShort.Write(this.Resolved);
}
