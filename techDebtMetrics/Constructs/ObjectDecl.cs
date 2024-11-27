using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

public class ObjectDecl : IObject, IDeclaration, IInitializable {

    public string Name { get; private set; } = "";

    public Location Location { get; private set; }

    public IReadOnlyList<MethodInst> Instances => this.inInstances.AsReadOnly();
    private readonly List<MethodInst> inInstances = [];

    public InterfaceDesc Interface => this.inInterface ??
        throw new UninitializedException(nameof(this.Interface));
    private InterfaceDesc? inInterface;

    public StructDesc Data => this.inData ??
        throw new UninitializedException(nameof(this.Data));
    private StructDesc? inData;

    public Package Package => this.inPackage ??
        throw new UninitializedException(nameof(this.Package));
    private Package? inPackage;

    public IReadOnlyList<TypeParam> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<TypeParam> inTypeParams = [];

    public IReadOnlyList<MethodDecl> Methods => this.inMethods.AsReadOnly();
    private readonly List<MethodDecl> inMethods = [];

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.Location = obj.TryReadLocation("loc", project);
        this.inInterface = obj.ReadIndex("interface", project.InterfaceDescs);
        this.inData = obj.ReadIndex("data", project.StructDescs);
        this.inPackage = obj.ReadIndex("package", project.Packages);
        obj.TryReadIndexList("instances", this.inInstances, project.MethodInsts);
        obj.TryReadIndexList("typeParams", this.inTypeParams, project.TypeParams);
        obj.TryReadIndexList("methods", this.inMethods, project.MethodDecls);
        this.Data.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (j.Long) j.Write("class ");
        j.Write(this.Name);
        if (j.Long) {
            j.WriteLine(" {");
            j.Indent.Write(this.Data.Fields, suffix: ";\n", separator: ";\n");
            j.Indent.Write(this.Methods, suffix: ";\n", separator: ";\n");
            j.Write("}");
        }
    }
}
