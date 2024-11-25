using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class InterfaceDecl : ITypeDesc, IDeclaration, IInitializable {
    public string Name { get; private set; } = "";

    public Location Location { get; private set; }

    public IReadOnlyList<InterfaceInst> Instances => this.inInstances.AsReadOnly();
    private List<InterfaceInst> inInstances = [];

    public InterfaceDesc Interface => this.inInterface ??
        throw new UninitializedException(nameof(this.Interface));
    private InterfaceDesc? inInterface;

    public Package Package => this.inPackage ??
        throw new UninitializedException(nameof(this.Package));
    private Package? inPackage;

    public IReadOnlyList<TypeParam> TypeParams => this.inTypeParams.AsReadOnly();
    private List<TypeParam> inTypeParams = [];

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name         = obj.ReadString("name");
        this.Location     = obj.TryReadLocation("loc", project);
        this.inInstances  = obj.TryReadIndexList("instances", project.InterfaceInsts);
        this.inInterface  = obj.ReadIndex("interface", project.InterfaceDescs);
        this.inPackage    = obj.ReadIndex("package", project.Packages);
        this.inTypeParams = obj.TryReadIndexList("typeParams", project.TypeParams);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write("interface ").Write(this.Name);
        if (j.Long)
            j.Write(this.TypeParams, "<", ">").
                Write(this.Interface);
    }
}
