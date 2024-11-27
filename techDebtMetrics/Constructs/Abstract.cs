using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

/// <summary>A named function signature in an interface.</summary>
/// <see cref="../../docs/genFeatureDef.md#abstract"/>
public class Abstract : IConstruct, IInitializable {

    /// <summary>The name of the function signature.</summary>
    public string Name { get; private set; } = "";

    /// <summary>The signature that is being named.</summary>
    public Signature Signature => this.inSignature ??
        throw new UninitializedException(nameof(this.Signature));
    private Signature? inSignature;

    void IInitializable.Initialize(Project project, Node node) {
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
