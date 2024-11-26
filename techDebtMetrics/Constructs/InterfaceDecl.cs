using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class InterfaceDecl : IInterface, IDeclaration, IInitializable {
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
        this.Name = obj.ReadString("name");
        this.Location = obj.TryReadLocation("loc", project);
        this.inInstances = obj.TryReadIndexList("instances", project.InterfaceInsts);
        this.inInterface = obj.ReadIndex("interface", project.InterfaceDescs);
        this.inPackage = obj.ReadIndex("package", project.Packages);
        this.inTypeParams = obj.TryReadIndexList("typeParams", project.TypeParams);
        this.Interface.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (j.Long) j.Write("interface ");
        j.Write(this.Name).Write(this.TypeParams, "<", ">");
        if (j.Long) j.Write(" ").Write(this.Interface);
    }
}
