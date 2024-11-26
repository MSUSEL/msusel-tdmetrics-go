using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;
using System.Linq;

namespace Constructs;

public class MethodInst : IMethod, IConstruct, IInitializable {
    public string Name => this.Generic.Name;

    public MethodDecl Generic => this.inGeneric ??
        throw new UninitializedException(nameof(this.Generic));
    private MethodDecl? inGeneric;

    public IReadOnlyList<ITypeDesc> InstanceTypes => this.inInstanceTypes.AsReadOnly();
    private List<ITypeDesc> inInstanceTypes = [];

    public ObjectInst? Receiver { get; private set; }

    public Signature Signature => this.inSignature ??
        throw new UninitializedException(nameof(this.Signature));
    private Signature? inSignature;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inGeneric       = obj.ReadIndex("generic", project.MethodDecls);
        this.inInstanceTypes = obj.ReadKeyList<ITypeDesc>("instanceTypes", project);
        this.Receiver        = obj.TryReadIndex("receiver", project.ObjectInsts);
        this.inSignature     = obj.ReadIndex("resolved", project.Signatures);
        this.Signature.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Receiver, suffix: ".").Write(this.Name).
            Write(this.InstanceTypes, "<", ">").Write(this.Signature);
}
