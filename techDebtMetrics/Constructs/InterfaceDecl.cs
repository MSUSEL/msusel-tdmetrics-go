using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>A declaration of an interface.</summary>
/// <see cref="../../docs/genFeatureDef.md#interface-declaration"/>
public class InterfaceDecl : IInterface, IDeclaration, IInitializable {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The name of the interface declaration.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The location the interface was defined.</summary>
    public Location Location { get; private set; }

    /// <summary>The interface type declaration.</summary>
    public InterfaceDesc Interface => this.inInterface ??
        throw new UninitializedException(nameof(this.Interface));
    private InterfaceDesc? inInterface;

    /// <summary>The package the interface was declared in.</summary>
    public Package Package => this.inPackage ??
        throw new UninitializedException(nameof(this.Package));
    private Package? inPackage;

    /// <summary>The type parameters for this interface if the interface is generic.</summary>
    public IReadOnlyList<TypeParam> TypeParams => this.inTypeParams.AsReadOnly();
    private readonly List<TypeParam> inTypeParams = [];

    /// <summary>The instances for this declaration if the interface is generic.</summary>
    public IReadOnlyList<InterfaceInst> Instances => this.inInstances.AsReadOnly();
    private readonly List<InterfaceInst> inInstances = [];

    /// <summary>True if this interface is generic, false otherwise.</summary>
    public bool Generic => this.TypeParams.Count > 0;

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs {
        get {
            foreach (IConstruct c in this.TypeParams) yield return c;
            yield return this.Interface;
            foreach (IConstruct c in this.Instances) yield return c;
        }
    }

    void IInitializable.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.Location = obj.TryReadLocation("loc", project);
        this.inInterface = obj.ReadIndex("interface", project.InterfaceDescs);
        this.inPackage = obj.ReadIndex("package", project.Packages);
        obj.TryReadIndexList("typeParams", this.inTypeParams, project.TypeParams);
        obj.TryReadIndexList("instances", this.inInstances, project.InterfaceInsts);
        this.Interface.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (j.Long) j.Write("interface ");
        j.Write(this.Name).Write(this.TypeParams, "<", ">");
        if (j.Long) j.Write(" ").Write(this.Interface);
        foreach (InterfaceInst inst in this.Instances)
            j.WriteLine().AsShort.Write("inst ").Write(inst);
    }
}
