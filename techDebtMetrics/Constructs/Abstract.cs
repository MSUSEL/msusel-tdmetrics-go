using Commons.Data.Yaml;
using Constructs.Exceptions;
using Constructs.Tooling;
using System.Collections.Generic;

namespace Constructs;

/// <summary>A named function signature in an interface.</summary>
/// <see cref="../../docs/genFeatureDef.md#abstract"/>
public class Abstract : IConstruct, IInitializable<Project> {

    /// <summary>Gets the index of this construct in the project list.</summary>
    public int Index { get; private set; } = 0;

    /// <summary>The name of the function signature.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The signature that is being named.</summary>
    public Signature Signature => this.inSignature ??
        throw new UninitializedException(nameof(this.Signature));
    private Signature? inSignature;

    /// <summary>Enumerates all the constructs that are directly part of this construct.</summary>
    public IEnumerable<IConstruct> SubConstructs => [this.Signature];

    void IInitializable<Project>.Initialize(Project project, int index, Node node) {
        this.Index = index;
        Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.inSignature = obj.ReadIndex("signature", project.Signatures);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write(this.Name);
        if (j.Long)
            j.AsShort.Write(this.Signature);
    }
}
