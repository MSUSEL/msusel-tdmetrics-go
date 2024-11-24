using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

public class Abstract : IConstruct, IInitializable {
    public string Name { get; private set; } = "";

    public Signature Signature => this.inSignature ??
        throw new UninitializedException("signature");
    private Signature? inSignature;

    void IInitializable.Initialize(Project project, Node node) {
        Object obj = node.AsObject();
        this.Name        = obj.ReadString("name");
        this.inSignature = obj.ReadIndex("signature", project.Signatures);
    }

    public override string ToString() => Journal.ToString(this);

    public void ToStub(Journal j) {
        j.Write(this.Name);
        if (j.Long)
            j.Write(" ").AsShort.Write(this.Signature);
    }
}
