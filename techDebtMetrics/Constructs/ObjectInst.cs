using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class ObjectInst : ITypeDesc, IInitializable {

    public ObjectDecl Generic => this.inGeneric ??
        throw new UninitializedException(nameof(this.Generic));
    private ObjectDecl? inGeneric;

    public IReadOnlyList<ITypeDesc> InstanceTypes => this.inInstanceTypes.AsReadOnly();
    private List<ITypeDesc> inInstanceTypes = [];

    public ObjectInst? Receiver { get; private set; }

    public IReadOnlyList<MethodInst> Methods => this.inMethods.AsReadOnly();
    private List<MethodInst> inMethods = [];

    public StructDesc ResolvedData => this.inResolvedData ??
        throw new UninitializedException(nameof(this.ResolvedData));
    private StructDesc? inResolvedData;

    public InterfaceDesc ResolvedInterface => this.inResolvedInterface ??
        throw new UninitializedException(nameof(this.ResolvedInterface));
    private InterfaceDesc? inResolvedInterface;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inGeneric           = obj.ReadIndex("generics", project.ObjectDecls);
        this.inInstanceTypes     = obj.ReadKeyList<ITypeDesc>("instances", project);
        this.Receiver            = obj.TryReadIndex("receiver", project.ObjectInsts);
        this.inMethods           = obj.TryReadKeyList<MethodInst>("methods", project);
        this.inResolvedData      = obj.ReadIndex("resData", project.StructDescs);
        this.inResolvedInterface = obj.ReadIndex("resInterface", project.InterfaceDescs);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Receiver, suffix: ".").Write(this.Generic).
            Write(this.InstanceTypes, "<", ">").Write(this.ResolvedData);
}
