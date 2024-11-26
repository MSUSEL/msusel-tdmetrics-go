using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class InterfaceInst : IInterface, IInitializable {
    public string Name => this.Generic.Name;

    public InterfaceDecl Generic => this.inGeneric ??
        throw new UninitializedException(nameof(this.Generic));
    private InterfaceDecl? inGeneric;

    public IReadOnlyList<ITypeDesc> InstanceTypes => this.inInstanceTypes.AsReadOnly();
    private List<ITypeDesc> inInstanceTypes = [];

    public InterfaceDesc Interface => this.inInterface ??
        throw new UninitializedException(nameof(this.Interface));
    private InterfaceDesc? inInterface;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inGeneric       = obj.ReadIndex("generic", project.InterfaceDecls);
        this.inInstanceTypes = obj.ReadKeyList<ITypeDesc>("instanceTypes", project);
        this.inInterface     = obj.ReadIndex("resolved", project.InterfaceDescs);
        this.Interface.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write(this.Name).
            AsLong.Write(this.InstanceTypes, "<", ">");
        if (j.Long) j.AsShort.Write(this.Interface);
    }
}
