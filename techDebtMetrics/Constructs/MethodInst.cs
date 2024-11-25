using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class MethodInst : IConstruct, IInitializable {

    public MethodDecl Generic => this.inGeneric ??
        throw new UninitializedException(nameof(this.Generic));
    private MethodDecl? inGeneric;

    public IReadOnlyList<ITypeDesc> InstanceTypes => this.inInstanceTypes.AsReadOnly();
    private List<ITypeDesc> inInstanceTypes = [];

    public ObjectInst? Receiver { get; private set; }

    public Signature Resolved => this.inResolved ??
        throw new UninitializedException(nameof(this.Resolved));
    private Signature? inResolved;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inGeneric       = obj.ReadIndex("generic", project.MethodDecls);
        this.inInstanceTypes = obj.ReadKeyList<ITypeDesc>("instanceTypes", project);
        this.Receiver        = obj.TryReadIndex("receiver", project.ObjectInsts);
        this.inResolved      = obj.ReadIndex("resolved", project.Signatures);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Receiver, suffix: ".").Write(this.Generic).
            Write(this.InstanceTypes, "<", ">").Write(this.Resolved);
}
