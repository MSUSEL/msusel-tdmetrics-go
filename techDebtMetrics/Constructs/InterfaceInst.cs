using Commons.Data.Yaml;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>An interface instance of a generic interface.</summary>
/// <see cref="../../docs/genFeatureDef.md#interface-instance"/>
public class InterfaceInst : IInterface, IInitializable<Project> {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The name of this interface.</summary>
    public string Name => this.Generic.Name;

    /// <summary>The generic interface this is an instance of.</summary>
    public InterfaceDecl Generic => this.inGeneric ??
        throw new UninitializedException(nameof(this.Generic));
    private InterfaceDecl? inGeneric;

    /// <summary>The interface type description for this instance.</summary>
    public InterfaceDesc Interface => this.inInterface ??
        throw new UninitializedException(nameof(this.Interface));
    private InterfaceDesc? inInterface;

    /// <summary>The type arguments put into the type parameters to create this instance.</summary>
    public IReadOnlyList<ITypeDesc> InstanceTypes => this.inInstanceTypes.AsReadOnly();
    private readonly List<ITypeDesc> inInstanceTypes = [];

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs {
        get {
            foreach (IConstruct c in this.InstanceTypes) yield return c;
            yield return this.Interface;
        }
    }

    void IInitializable<Project>.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.inGeneric = obj.ReadIndex("generic", project.InterfaceDecls);
        this.inInterface = obj.ReadIndex("resolved", project.InterfaceDescs);
        obj.ReadKeyList(project, "instanceTypes", this.inInstanceTypes);
        this.Interface.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write(this.Name).
            AsShort.Write(this.InstanceTypes, "<", ">");
        if (j.Long) j.Write(this.Interface);
    }
}
