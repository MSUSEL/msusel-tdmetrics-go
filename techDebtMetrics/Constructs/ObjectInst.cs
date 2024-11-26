using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class ObjectInst : IObject, IInitializable {
    public string Name => this.Generic.Name;

    public ObjectDecl Generic => this.inGeneric ??
        throw new UninitializedException(nameof(this.Generic));
    private ObjectDecl? inGeneric;

    public IReadOnlyList<ITypeDesc> InstanceTypes => this.inInstanceTypes.AsReadOnly();
    private List<ITypeDesc> inInstanceTypes = [];

    public IReadOnlyList<MethodInst> Methods => this.inMethods.AsReadOnly();
    private List<MethodInst> inMethods = [];

    public InterfaceDesc Interface => this.inInterface ??
        throw new UninitializedException(nameof(this.Interface));
    private InterfaceDesc? inInterface;

    public StructDesc Data => this.inResolvedData ??
        throw new UninitializedException(nameof(this.Data));
    private StructDesc? inResolvedData;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inGeneric = obj.ReadIndex("generic", project.ObjectDecls);
        this.inInstanceTypes = obj.ReadKeyList<ITypeDesc>("instanceTypes", project);
        this.inMethods = obj.TryReadKeyList<MethodInst>("methods", project);
        this.inInterface = obj.ReadIndex("resInterface", project.InterfaceDescs);
        this.inResolvedData = obj.ReadIndex("resData", project.StructDescs);
        this.Data.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Name).Write(this.InstanceTypes, "<", ">").Write(this.Data);
}
