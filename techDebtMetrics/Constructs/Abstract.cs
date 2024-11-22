using Constructs.Data;
using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

public class Abstract : ITypeDesc, IInitializable {
    public string Name { get; private set; } = "";

    public Signature Signature => this.inSignature ??
        throw new UninitializedException("signature");
    private Signature? inSignature;

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.inSignature = obj.ReadIndexType<Signature>("signature", getter);
    }

    public override string ToString() => this.Name + ":" + this.inSignature;

    public void ToStub(Journal j) {
        j.Write(this.Name);
        if (!j.Short)
            j.Write(" ").AsShort.Write(this.Signature);
    }
}
