using Constructs.Exceptions;
using Constructs.Tooling;

namespace Constructs;

public class Method : IConstruct, IInitializable {
    public string Name { get; private set; } = "";

    public ITypeDesc Signature => this.inSignature ??
        throw new UninitializedException("signature");
    private ITypeDesc? inSignature;

    void IInitializable.Initialize(TypeGetter getter, Data.Node node) {
        Data.Object obj = node.AsObject();
        this.Name = obj.ReadString("name");
        this.inSignature = obj.ReadIndexType<ITypeDesc>("signature", getter);

        // TODO: Add "metrics"
    }

    public void ToStub(Journal j) =>
        j.Write(this.Name).AsShort.Write(this.Signature);
}
