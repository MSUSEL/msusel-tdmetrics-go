using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class MethodInst : IMethod, IConstruct, IInitializable {
    public string Name => this.Generic.Name;

    public MethodDecl Generic => this.inGeneric ??
        throw new UninitializedException(nameof(this.Generic));
    private MethodDecl? inGeneric;

    public IReadOnlyList<ITypeDesc> InstanceTypes => this.inInstanceTypes.AsReadOnly();
    private readonly List<ITypeDesc> inInstanceTypes = [];

    public ObjectInst? Receiver { get; private set; }

    public Signature Signature => this.inSignature ??
        throw new UninitializedException(nameof(this.Signature));
    private Signature? inSignature;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inGeneric = obj.ReadIndex("generic", project.MethodDecls);
        this.Receiver = obj.TryReadIndex("receiver", project.ObjectInsts);
        this.inSignature = obj.ReadIndex("resolved", project.Signatures);
        obj.ReadKeyList("instanceTypes", this.inInstanceTypes, project);
        this.Signature.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Receiver, suffix: ".").Write(this.Name).
            Write(this.InstanceTypes, "<", ">").Write(this.Signature);
}
