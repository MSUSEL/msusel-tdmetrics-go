using Commons.Data.Locations;
using Commons.Data.Yaml;
using Constructs.Exceptions;
using Constructs.Extensions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>A declaration of an interface.</summary>
/// <see cref="../../docs/genFeatureDef.md#object"/>
public class ObjectDecl : IObject, IDeclaration, IInitializable<Project> {

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

    /// <summary>The data contained by this object.</summary>
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
    public IReadOnlyList<ObjectInst> Instances => this.inInstances.AsReadOnly();
    private readonly List<ObjectInst> inInstances = [];

    /// <summary>Optional method this object was defined inside of.</summary>
    public IMethod? Nest { get; private set; }

    /// <summary>The implicit type parameters from the nest method for this object.</summary>
    public IReadOnlyList<TypeParam> ImplicitTypeParams =>
        this.Nest is MethodDecl decl ? decl.TypeParams : [];

    /// <summary>True if this object is generic, false otherwise.</summary>
    public bool Generic => this.TypeParams.Count > 0;

    /// <summary>True if this object is nested, false otherwise.</summary>
    public bool Nested => this.Nest is not null;

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs {
        get {
            foreach (IConstruct c in this.TypeParams) yield return c;
            yield return this.Data;
            yield return this.Interface;
            foreach (IConstruct c in this.Methods) yield return c;
            foreach (IConstruct c in this.Instances) yield return c;
        }
    }

    void IInitializable<Project>.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.Location = obj.TryReadLocation(project.Locations, "loc");
        this.inInterface = obj.ReadIndex("interface", project.InterfaceDescs);
        this.inData = obj.ReadIndex("data", project.StructDescs);
        this.inPackage = obj.ReadIndex("package", project.Packages);
        obj.TryReadIndexList("typeParams", this.inTypeParams, project.TypeParams);
        obj.TryReadIndexList("methods", this.inMethods, project.MethodDecls);
        obj.TryReadIndexList("instances", this.inInstances, project.ObjectInsts);
        this.Nest = obj.TryReadKey<IMethod>(project, "nest");
        this.Data.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        if (j.Long) j.Write("class ");
        if (this.Nest is not null) {
            j.Write(this.Nest.Name);
            if (j.Long)
                j.Write(this.ImplicitTypeParams, "<", ">");
            j.Write(":");
        }
        j.Write(this.Name);
        if (j.Long) {
            j.Write(this.TypeParams, "<", ">");
            j.WriteLine(" {");
            j.Indent.Write(this.Data.Fields, suffix: ";\n", separator: ";\n");
            j.Indent.Write(this.Methods, suffix: ";\n", separator: ";\n");
            j.Write("}");
            foreach (ObjectInst inst in this.Instances) {
                if (inst.IsConcrete(true))
                    j.WriteLine().AsShort.Write("inst ").Write(inst);
            }
        }
    }
}
