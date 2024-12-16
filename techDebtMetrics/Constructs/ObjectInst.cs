using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>An object instance of a generic object.</summary>
/// <see cref="../../docs/genFeatureDef.md#object-instance"/>
public class ObjectInst : IObject, IInitializable {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The name of the object declaration.</summary>
    public string Name => this.Generic.Name;

    /// <summary>The generic object this is an instance of.</summary>
    public ObjectDecl Generic => this.inGeneric ??
        throw new UninitializedException(nameof(this.Generic));
    private ObjectDecl? inGeneric;

    /// <summary>The interface that summarizes the methods of this object.</summary>
    public InterfaceDesc Interface => this.inInterface ??
        throw new UninitializedException(nameof(this.Interface));
    private InterfaceDesc? inInterface;

    /// <summary>The data contated by this object.</summary>
    public StructDesc Data => this.inResolvedData ??
        throw new UninitializedException(nameof(this.Data));
    private StructDesc? inResolvedData;

    /// <summary>The type arguments put into the type parameters to create this instance.</summary>
    public IReadOnlyList<ITypeDesc> InstanceTypes => this.inInstanceTypes.AsReadOnly();
    private readonly List<ITypeDesc> inInstanceTypes = [];

    /// <summary>The methods that have been declared as members to this object.</summary>
    public IReadOnlyList<MethodInst> Methods => this.inMethods.AsReadOnly();
    private readonly List<MethodInst> inMethods = [];

    void IInitializable.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.inGeneric = obj.ReadIndex("generic", project.ObjectDecls);
        this.inInterface = obj.ReadIndex("resInterface", project.InterfaceDescs);
        this.inResolvedData = obj.ReadIndex("resData", project.StructDescs);
        obj.ReadKeyList("instanceTypes", this.inInstanceTypes, project);
        obj.TryReadKeyList("methods", this.inMethods, project);
        this.Data.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Name).Write(this.InstanceTypes, "<", ">").Write(this.Data);
}
