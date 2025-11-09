using Commons.Data.Yaml;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>A method instance of a generic method.</summary>
/// <see cref="../../docs/genFeatureDef.md#method-instance"/>
public class MethodInst : IMethod, IConstruct, IInitializable<Project> {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The name of this method.</summary>
    public string Name => this.Generic.Name;

    /// <summary>The generic method this is an instance of.</summary>
    public MethodDecl Generic => this.inGeneric ??
        throw new UninitializedException(nameof(this.Generic));
    private MethodDecl? inGeneric;

    /// <summary>The type arguments put into the type parameters to create this instance.</summary>
    public IReadOnlyList<ITypeDesc> InstanceTypes => this.inInstanceTypes.AsReadOnly();
    private readonly List<ITypeDesc> inInstanceTypes = [];

    /// <summary>
    /// This the optional object instance this method instance is defined as a member of
    /// or null if the method is a function instance on it's own in the package.
    /// </summary>
    public ObjectInst? Receiver { get; private set; }

    /// <summary>The method signature type description for this instance.</summary>
    public Signature Signature => this.inSignature ??
        throw new UninitializedException(nameof(this.Signature));
    private Signature? inSignature;

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs {
        get {
            foreach (IConstruct c in this.InstanceTypes) yield return c;
            yield return this.Signature;
        }
    }

    void IInitializable<Project>.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.inGeneric = obj.ReadIndex("generic", project.MethodDecls);
        this.Receiver = obj.TryReadIndex("receiver", project.ObjectInsts);
        this.inSignature = obj.ReadIndex("resolved", project.Signatures);
        obj.ReadKeyList(project, "instanceTypes", this.inInstanceTypes);
        this.Signature.AddUses(this);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) =>
        j.AsShort.Write(this.Receiver, suffix: ".").Write(this.Name).
            Write(this.InstanceTypes, "<", ">").Write(this.Signature);
}
