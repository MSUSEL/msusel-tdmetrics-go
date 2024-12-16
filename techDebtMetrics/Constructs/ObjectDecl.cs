using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>A declaration of an interface.</summary>
/// <see cref="../../docs/genFeatureDef.md#object"/>
public class ObjectDecl : IObject, IDeclaration, IInitializable {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The name of the object declaration.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The location the object was defined.</summary>
    public Location Location { get; private set; }

    /// <summary>The interface that summarizes the methods of this object.</summary>
    public InterfaceDesc Interface => this.inInterface ??
        throw new UninitializedException(nameof(this.Interface));
    private InterfaceDesc? inInterface;

    /// <summary>The data contated by this object.</summary>
    public StructDesc Data => this.inData ??
        throw new UninitializedException(nameof(this.Data));
    private StructDesc? inData;

    /// <summary>The package the object was declared in.</summary>
    public Package Package => this.inPackage ??
        throw new UninitializedException(nameof(this.Package));
    private Package? inPackage;

    /// <summary>The type parameters for this object if the object is generic.</summary>
    public IReadOnlyList<TypeParam> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<TypeParam> inTypeParams = [];

    /// <summary>The methods that have been declared as members to this object.</summary>
    public IReadOnlyList<MethodDecl> Methods => this.inMethods.AsReadOnly();
    private readonly List<MethodDecl> inMethods = [];

    /// <summary>The instances for this declaration if the object is generic.</summary>
    public IReadOnlyList<MethodInst> Instances => this.inInstances.AsReadOnly();
    private readonly List<MethodInst> inInstances = [];

    /// <summary>True if this object is generic, false otherwise.</summary>
    public bool Generic => this.TypeParams.Count > 0;

    void IInitializable.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.Location = obj.TryReadLocation("loc", project);
        this.inInterface = obj.ReadIndex("interface", project.InterfaceDescs);
        this.inData = obj.ReadIndex("data", project.StructDescs);
        this.inPackage = obj.ReadIndex("package", project.Packages);
        obj.TryReadIndexList("typeParams", this.inTypeParams, project.TypeParams);
        obj.TryReadIndexList("methods", this.inMethods, project.MethodDecls);
        obj.TryReadIndexList("instances", this.inInstances, project.MethodInsts);
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
