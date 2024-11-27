using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>An interface instance of a generic interface.</summary>
/// <see cref="../../docs/genFeatureDef.md#interface-instance"/>
public class InterfaceInst : IInterface, IInitializable {

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

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.inGeneric = obj.ReadIndex("generic", project.InterfaceDecls);
        this.inInterface = obj.ReadIndex("resolved", project.InterfaceDescs);
        obj.ReadKeyList("instanceTypes", this.inInstanceTypes, project);
        this.Interface.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write(this.Name).
            AsShort.Write(this.InstanceTypes, "<", ">");
        if (j.Long) j.AsShort.Write(this.Interface);
    }
}
